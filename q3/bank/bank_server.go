// File: bank/bank_server.go
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
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	//pb "stripe" // generated from stripe.proto
	pb "github.com/Mrinal92/Distri-Assignment2/q2/protofiles"
)

type Account struct {
	AccountNumber string
	Balance       float64
}

type PendingTransaction struct {
	Amount float64
}

type BankServer struct {
	pb.UnimplementedBankServer
	bankName    string
	accounts    map[string]*Account
	pendingTxns map[string]PendingTransaction // key: transaction id
	mu          sync.Mutex
}

func NewBankServer(bankName string) *BankServer {
	// For simulation, create dummy accounts.
	accounts := make(map[string]*Account)
	// For example, if bankName is "BankA", we create accounts like "A123", "A456".
	prefix := string(bankName[len(bankName)-1])
	accounts[prefix+"123"] = &Account{AccountNumber: prefix + "123", Balance: 1000.0}
	accounts[prefix+"456"] = &Account{AccountNumber: prefix + "456", Balance: 2000.0}
	return &BankServer{
		bankName:    bankName,
		accounts:    accounts,
		pendingTxns: make(map[string]PendingTransaction),
	}
}

// PreparePayment locks funds (for debit) or marks a pending credit.
func (s *BankServer) PreparePayment(ctx context.Context, req *pb.PrepareRequest) (*pb.PrepareResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account, exists := s.accounts[req.AccountNumber]
	if !exists {
		return &pb.PrepareResponse{
			CanCommit: false,
			Message:   "Account not found",
		}, nil
	}
	if req.IsDebit {
		if account.Balance < req.Amount {
			return &pb.PrepareResponse{
				CanCommit: false,
				Message:   "Insufficient funds",
			}, nil
		}
		// Lock funds by deducting temporarily.
		account.Balance -= req.Amount
		s.pendingTxns[req.TransactionId] = PendingTransaction{Amount: req.Amount}
		log.Printf("[%s] Prepared debit for account %s, amount: %f", s.bankName, req.AccountNumber, req.Amount)
	} else {
		// For credit, record a pending transaction.
		s.pendingTxns[req.TransactionId] = PendingTransaction{Amount: req.Amount}
		log.Printf("[%s] Prepared credit for account %s, amount: %f", s.bankName, req.AccountNumber, req.Amount)
	}
	return &pb.PrepareResponse{
		CanCommit: true,
		Message:   "Prepared successfully",
	}, nil
}

// CommitPayment finalizes the transaction.
func (s *BankServer) CommitPayment(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.pendingTxns[req.TransactionId]
	if !exists {
		return &pb.CommitResponse{
			Success: false,
			Message: "No pending transaction found",
		}, nil
	}
	account, exists := s.accounts[req.AccountNumber]
	if !exists {
		return &pb.CommitResponse{
			Success: false,
			Message: "Account not found",
		}, nil
	}
	if req.IsDebit {
		// Funds were already deducted in prepare.
		log.Printf("[%s] Committed debit for account %s, amount: %f", s.bankName, req.AccountNumber, req.Amount)
	} else {
		// Add funds for credit.
		account.Balance += req.Amount
		log.Printf("[%s] Committed credit for account %s, amount: %f", s.bankName, req.AccountNumber, req.Amount)
	}
	delete(s.pendingTxns, req.TransactionId)
	return &pb.CommitResponse{
		Success: true,
		Message: "Commit successful",
	}, nil
}

// RollbackPayment reverts the changes made during prepare.
func (s *BankServer) RollbackPayment(ctx context.Context, req *pb.RollbackRequest) (*pb.RollbackResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pending, exists := s.pendingTxns[req.TransactionId]
	if !exists {
		return &pb.RollbackResponse{
			Success: false,
			Message: "No pending transaction found",
		}, nil
	}
	account, exists := s.accounts[req.AccountNumber]
	if !exists {
		return &pb.RollbackResponse{
			Success: false,
			Message: "Account not found",
		}, nil
	}
	if req.IsDebit {
		// Revert the deducted amount.
		account.Balance += pending.Amount
		log.Printf("[%s] Rolled back debit for account %s, amount: %f", s.bankName, req.AccountNumber, pending.Amount)
	} else {
		log.Printf("[%s] Rolled back credit for account %s, amount: %f", s.bankName, req.AccountNumber, pending.Amount)
	}
	delete(s.pendingTxns, req.TransactionId)
	return &pb.RollbackResponse{
		Success: true,
		Message: "Rollback successful",
	}, nil
}

// loadTLSCredentials loads TLS credentials for the bank server.
func loadTLSCredentials(caFile, certFile, keyFile string) (credentials.TransportCredentials, error) {
	pemServerCA, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, errors.New("failed to add CA's certificate")
	}
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	return credentials.NewTLS(config), nil
}

func main() {
	bankName := flag.String("bank", "BankA", "Name of the bank")
	port := flag.Int("port", 50051, "Bank server port")
	flag.Parse()

	tlsCreds, err := loadTLSCredentials("certs/ca.crt", "certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(tlsCreds)}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterBankServer(grpcServer, NewBankServer(*bankName))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("%s server listening on port %d", *bankName, *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
