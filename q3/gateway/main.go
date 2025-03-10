// File: gateway/main.go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	//pb "stripe" // generated from stripe.proto
	pb "github.com/Mrinal92/Distri-Assignment2/q2/protofiles"
)

// ----- User and Gateway server definition -----
type User struct {
	Username      string
	Password      string
	Token         string
	AccountNumber string
}

type PaymentGatewayServer struct {
	pb.UnimplementedPaymentGatewayServer
	userStore        map[string]User                // keyed by username
	tokenToUser      map[string]User                // maps token to user
	mu               sync.Mutex
	idempotencyStore map[string]*pb.PaymentResponse // idempotency key -> response
	bankMapping      map[string]string              // mapping of account prefix to bank server address
}

func NewPaymentGatewayServer() *PaymentGatewayServer {
	// For simplicity, no preloaded users.
	// We map accounts starting with "A" to BankA (localhost:50051) and "B" to BankB (localhost:50052).
	bankMapping := map[string]string{
		"A": "localhost:50051",
		"B": "localhost:50052",
	}
	return &PaymentGatewayServer{
		userStore:        make(map[string]User),
		tokenToUser:      make(map[string]User),
		idempotencyStore: make(map[string]*pb.PaymentResponse),
		bankMapping:      bankMapping,
	}
}

// Register a new user.
func (s *PaymentGatewayServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.userStore[req.Username]; exists {
		return &pb.RegisterResponse{
			Success: false,
			Message: "User already exists",
		}, nil
	}
	// Create and store new user.
	user := User{
		Username:      req.Username,
		Password:      req.Password,
		AccountNumber: req.AccountNumber,
	}
	s.userStore[req.Username] = user
	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

// Authenticate user and return a token.
func (s *PaymentGatewayServer) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, exists := s.userStore[req.Username]
	if !exists || user.Password != req.Password {
		return &pb.AuthResponse{
			Success: false,
			Message: "Invalid credentials",
		}, nil
	}
	// Generate a simple random token.
	token := fmt.Sprintf("%d", rand.Int())
	user.Token = token
	s.userStore[req.Username] = user
	s.tokenToUser[token] = user
	return &pb.AuthResponse{
		Success: true,
		Token:   token,
		Message: "Authenticated successfully",
	}, nil
}

// ProcessPayment implements 2PC with idempotency.
func (s *PaymentGatewayServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	// Check idempotency first.
	s.mu.Lock()
	if resp, exists := s.idempotencyStore[req.IdempotencyKey]; exists {
		s.mu.Unlock()
		log.Printf("Idempotent payment detected, returning previous response")
		return resp, nil
	}
	s.mu.Unlock()

	// Retrieve authenticated user from metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing metadata")
	}
	tokens := md["authorization"]
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Missing auth token")
	}
	token := tokens[0]
	s.mu.Lock()
	user, exists := s.tokenToUser[token]
	s.mu.Unlock()
	if !exists {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}
	// Ensure the sender account matches the authenticated user.
	if req.SenderAccount != user.AccountNumber {
		return &pb.PaymentResponse{
			Success: false,
			Message: "Sender account does not match authenticated user",
		}, nil
	}

	// Determine bank servers based on the first character of account numbers.
	senderPrefix := string(req.SenderAccount[0])
	receiverPrefix := string(req.ReceiverAccount[0])
	senderBankAddr, ok1 := s.bankMapping[senderPrefix]
	receiverBankAddr, ok2 := s.bankMapping[receiverPrefix]
	if !ok1 || !ok2 {
		return &pb.PaymentResponse{
			Success: false,
			Message: "Bank server not found for account",
		}, nil
	}

	// Set timeout for the 2PC process.
	timeout := 5 * time.Second
	ctx2, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create TLS credentials for bank connections.
	creds, err := loadTLSCredentials("certs/ca.crt", "certs/client.crt", "certs/client.key")
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to load TLS credentials for bank connection")
	}

	// Connect to sender bank.
	senderConn, err := grpc.Dial(senderBankAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return &pb.PaymentResponse{
			Success: false,
			Message: "Failed to connect to sender bank",
		}, nil
	}
	defer senderConn.Close()
	senderClient := pb.NewBankClient(senderConn)

	// Connect to receiver bank.
	receiverConn, err := grpc.Dial(receiverBankAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return &pb.PaymentResponse{
			Success: false,
			Message: "Failed to connect to receiver bank",
		}, nil
	}
	defer receiverConn.Close()
	receiverClient := pb.NewBankClient(receiverConn)

	// ----- 2PC Phase 1: Prepare -----
	prepareSenderReq := &pb.PrepareRequest{
		TransactionId: req.TransactionId,
		AccountNumber: req.SenderAccount,
		Amount:        req.Amount,
		IsDebit:       true,
	}
	senderPrepResp, err := senderClient.PreparePayment(ctx2, prepareSenderReq)
	if err != nil || !senderPrepResp.CanCommit {
		if err != nil {
			log.Printf("Sender bank prepare error: %v", err)
		} else {
			log.Printf("Sender bank cannot commit: %s", senderPrepResp.Message)
		}
		// Rollback sender if needed.
		senderClient.RollbackPayment(ctx2, &pb.RollbackRequest{
			TransactionId: req.TransactionId,
			AccountNumber: req.SenderAccount,
			IsDebit:       true,
		})
		return &pb.PaymentResponse{
			Success: false,
			Message: "Transaction aborted during prepare phase at sender bank",
		}, nil
	}

	prepareReceiverReq := &pb.PrepareRequest{
		TransactionId: req.TransactionId,
		AccountNumber: req.ReceiverAccount,
		Amount:        req.Amount,
		IsDebit:       false,
	}
	receiverPrepResp, err := receiverClient.PreparePayment(ctx2, prepareReceiverReq)
	if err != nil || !receiverPrepResp.CanCommit {
		if err != nil {
			log.Printf("Receiver bank prepare error: %v", err)
		} else {
			log.Printf("Receiver bank cannot commit: %s", receiverPrepResp.Message)
		}
		// Rollback sender if needed.
		senderClient.RollbackPayment(ctx2, &pb.RollbackRequest{
			TransactionId: req.TransactionId,
			AccountNumber: req.SenderAccount,
			IsDebit:       true,
		})
		return &pb.PaymentResponse{
			Success: false,
			Message: "Transaction aborted during prepare phase at receiver bank",
		}, nil
	}

	// ----- 2PC Phase 2: Commit -----
	commitSenderReq := &pb.CommitRequest{
		TransactionId: req.TransactionId,
		AccountNumber: req.SenderAccount,
		Amount:        req.Amount,
		IsDebit:       true,
	}
	commitSenderResp, err := senderClient.CommitPayment(ctx2, commitSenderReq)
	if err != nil || !commitSenderResp.Success {
		receiverClient.RollbackPayment(ctx2, &pb.RollbackRequest{
			TransactionId: req.TransactionId,
			AccountNumber: req.ReceiverAccount,
			IsDebit:       false,
		})
		return &pb.PaymentResponse{
			Success: false,
			Message: "Transaction commit failed at sender bank",
		}, nil
	}

	commitReceiverReq := &pb.CommitRequest{
		TransactionId: req.TransactionId,
		AccountNumber: req.ReceiverAccount,
		Amount:        req.Amount,
		IsDebit:       false,
	}
	commitReceiverResp, err := receiverClient.CommitPayment(ctx2, commitReceiverReq)
	if err != nil || !commitReceiverResp.Success {
		// In a real system, additional compensating actions might be necessary.
		return &pb.PaymentResponse{
			Success: false,
			Message: "Transaction commit failed at receiver bank",
		}, nil
	}

	// Transaction successful; record in idempotency store.
	finalResp := &pb.PaymentResponse{
		Success: true,
		Message: "Transaction completed successfully",
	}
	s.mu.Lock()
	s.idempotencyStore[req.IdempotencyKey] = finalResp
	s.mu.Unlock()
	return finalResp, nil
}

// loadTLSCredentials loads TLS credentials from certificate files.
func loadTLSCredentials(caFile, certFile, keyFile string) (credentials.TransportCredentials, error) {
	// Load the CA certificate.
	pemServerCA, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, errors.New("failed to add server CA's certificate")
	}
	// Load the server's certificate and key.
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}
	return credentials.NewTLS(config), nil
}

// ----- gRPC Interceptors -----
// Logging interceptor logs every incoming request and its outcome.
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	log.Printf("Received request for method: %s, request: %+v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	if err != nil {
		log.Printf("Method %s failed: %v (Duration: %s)", info.FullMethod, err, duration)
	} else {
		log.Printf("Method %s succeeded (Duration: %s)", info.FullMethod, duration)
	}
	return resp, err
}

// Auth interceptor checks for the authorization token in metadata.
func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Allow unauthenticated access for Register and Authenticate.
	if info.FullMethod == "/stripe.PaymentGateway/Register" || info.FullMethod == "/stripe.PaymentGateway/Authenticate" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing metadata")
	}
	tokens := md["authorization"]
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Missing auth token")
	}
	// In a production system, token verification would be more robust.
	return handler(ctx, req)
}

// Helper to chain multiple interceptors.
func grpc_middleware_chain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		chainedHandler := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			currentInterceptor := interceptors[i]
			next := chainedHandler
			chainedHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return currentInterceptor(currentCtx, currentReq, info, next)
			}
		}
		return chainedHandler(ctx, req)
	}
}

func main() {
	port := flag.Int("port", 5000, "Payment Gateway server port")
	flag.Parse()

	tlsCreds, err := loadTLSCredentials("certs/ca.crt", "certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(tlsCreds),
		grpc.UnaryInterceptor(
			grpc_middleware_chain(loggingInterceptor, authInterceptor),
		),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterPaymentGatewayServer(server, NewPaymentGatewayServer())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Payment Gateway server listening on port %d", *port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
