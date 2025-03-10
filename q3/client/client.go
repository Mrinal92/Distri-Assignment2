// File: client/client.go
package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	//pb "stripe" // generated from stripe.proto
	pb "github.com/Mrinal92/Distri-Assignment2/q2/protofiles"
)

type PaymentClient struct {
	gatewayAddr  string
	conn         *grpc.ClientConn
	client       pb.PaymentGatewayClient
	tlsCreds     credentials.TransportCredentials
	token        string
	offlineQueue []*pb.PaymentRequest
	queueMu      sync.Mutex
}

func NewPaymentClient(gatewayAddr string) (*PaymentClient, error) {
	tlsCreds, err := loadTLSCredentials("certs/ca.crt", "certs/client.crt", "certs/client.key")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(gatewayAddr, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		return nil, err
	}
	client := pb.NewPaymentGatewayClient(conn)
	return &PaymentClient{
		gatewayAddr:  gatewayAddr,
		conn:         conn,
		client:       client,
		tlsCreds:     tlsCreds,
		offlineQueue: make([]*pb.PaymentRequest, 0),
	}, nil
}

func loadTLSCredentials(caFile, certFile, keyFile string) (credentials.TransportCredentials, error) {
	pemServerCA, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, errors.New("failed to add CA's certificate")
	}
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

func (pc *PaymentClient) Close() {
	pc.conn.Close()
}

func (pc *PaymentClient) Register(username, password, accountNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.RegisterRequest{
		Username:      username,
		Password:      password,
		AccountNumber: accountNumber,
	}
	resp, err := pc.client.Register(ctx, req)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("registration failed: %s", resp.Message)
	}
	log.Printf("Registration successful: %s", resp.Message)
	return nil
}

func (pc *PaymentClient) Authenticate(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.AuthRequest{
		Username: username,
		Password: password,
	}
	resp, err := pc.client.Authenticate(ctx, req)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("authentication failed: %s", resp.Message)
	}
	pc.token = resp.Token
	log.Printf("Authenticated successfully. Token: %s", resp.Token)
	return nil
}

func (pc *PaymentClient) ProcessPayment(pr *pb.PaymentRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Attach auth token in metadata.
	md := metadata.Pairs("authorization", pc.token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	resp, err := pc.client.ProcessPayment(ctx, pr)
	if err != nil {
		log.Printf("Payment processing error: %v. Queuing payment for retry.", err)
		pc.queuePayment(pr)
		return
	}
	if !resp.Success {
		log.Printf("Payment failed: %s. Queuing payment for retry.", resp.Message)
		pc.queuePayment(pr)
	} else {
		log.Printf("Payment succeeded: %s", resp.Message)
	}
}

func (pc *PaymentClient) queuePayment(pr *pb.PaymentRequest) {
	pc.queueMu.Lock()
	defer pc.queueMu.Unlock()
	pc.offlineQueue = append(pc.offlineQueue, pr)
}

// StartRetryLoop periodically resends queued payments.
func (pc *PaymentClient) StartRetryLoop() {
	go func() {
		for {
			time.Sleep(10 * time.Second) // configurable retry interval
			pc.queueMu.Lock()
			if len(pc.offlineQueue) == 0 {
				pc.queueMu.Unlock()
				continue
			}
			log.Printf("Retrying %d queued payments...", len(pc.offlineQueue))
			remaining := []*pb.PaymentRequest{}
			for _, pr := range pc.offlineQueue {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				md := metadata.Pairs("authorization", pc.token)
				ctx = metadata.NewOutgoingContext(ctx, md)
				resp, err := pc.client.ProcessPayment(ctx, pr)
				cancel()
				if err != nil || !resp.Success {
					log.Printf("Retry failed for transaction %s: %v", pr.TransactionId, err)
					remaining = append(remaining, pr)
				} else {
					log.Printf("Retry succeeded for transaction %s", pr.TransactionId)
				}
			}
			pc.offlineQueue = remaining
			pc.queueMu.Unlock()
		}
	}()
}

func main() {
	gatewayAddr := flag.String("gateway", "localhost:5000", "Payment Gateway address")
	username := flag.String("username", "user1", "Username")
	password := flag.String("password", "pass1", "Password")
	accountNumber := flag.String("account", "A123", "Account number")
	flag.Parse()

	pc, err := NewPaymentClient(*gatewayAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer pc.Close()

	// Register (ignore error if already registered).
	if err := pc.Register(*username, *password, *accountNumber); err != nil {
		log.Printf("Registration error (might be already registered): %v", err)
	}
	// Authenticate.
	if err := pc.Authenticate(*username, *password); err != nil {
		log.Fatalf("Authentication error: %v", err)
	}
	// Start the offline retry loop.
	pc.StartRetryLoop()

	// Interactive input for payments.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter payment details in format: transaction_id receiver_account amount idempotency_key")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}
		var txnID, receiverAcc, idemKey string
		var amount float64
		n, err := fmt.Sscanf(input, "%s %s %f %s", &txnID, &receiverAcc, &amount, &idemKey)
		if n < 4 || err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}
		paymentReq := &pb.PaymentRequest{
			TransactionId:   txnID,
			SenderAccount:   *accountNumber,
			ReceiverAccount: receiverAcc,
			Amount:          amount,
			IdempotencyKey:  idemKey,
		}
		pc.ProcessPayment(paymentReq)
	}
}
