package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	pb "github.com/Mrinal92/Distri-Assignment2/q1/protofiles"
)

type computeServer struct {
	pb.UnimplementedComputeServiceServer
	serverID    string
	lbAddress   string
	currentLoad int64 // Tracks active tasks.
}

// ComputeTask simulates a computation.
func (s *computeServer) ComputeTask(ctx context.Context, req *pb.ComputeRequest) (*pb.ComputeResponse, error) {
	atomic.AddInt64(&s.currentLoad, 1)
	defer atomic.AddInt64(&s.currentLoad, -1)

	log.Printf("[Server %s] Received task: %s\n", s.serverID, req.GetTaskData())
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return &pb.ComputeResponse{
		Result: fmt.Sprintf("Result from server %s for input '%s'", s.serverID, req.GetTaskData()),
	}, nil
}

func (s *computeServer) registerWithLB() {
	conn, err := grpc.Dial(s.lbAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to LB: %v", err)
	}
	defer conn.Close()

	lbClient := pb.NewLoadBalancerServiceClient(conn)
	_, err = lbClient.RegisterServer(context.Background(), &pb.RegisterRequest{
		ServerId: s.serverID,
		Address:  fmt.Sprintf("localhost:%s", *portFlag),
	})
	if err != nil {
		log.Fatalf("RegisterServer failed: %v", err)
	}
	log.Printf("[Server %s] Registered with LB at %s\n", s.serverID, s.lbAddress)
}

func (s *computeServer) reportLoadToLB() {
	conn, err := grpc.Dial(s.lbAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to LB: %v", err)
	}
	lbClient := pb.NewLoadBalancerServiceClient(conn)

	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		loadVal := float32(atomic.LoadInt64(&s.currentLoad))
		_, err := lbClient.ReportLoad(context.Background(), &pb.LoadReportRequest{
			ServerId:    s.serverID,
			CurrentLoad: loadVal,
		})
		if err != nil {
			log.Printf("[Server %s] ReportLoad error: %v", s.serverID, err)
		}
	}
}

var (
	portFlag   = flag.String("port", "9100", "Port to run the server on")
	lbAddrFlag = flag.String("lb_addr", "localhost:9000", "Address of the load balancer")
	idFlag     = flag.String("id", "server1", "Unique server ID")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+*portFlag)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *portFlag, err)
	}

	s := &computeServer{
		serverID:  *idFlag,
		lbAddress: *lbAddrFlag,
	}

	go s.registerWithLB()
	go s.reportLoadToLB()

	grpcServer := grpc.NewServer()
	pb.RegisterComputeServiceServer(grpcServer, s)

	log.Printf("Backend Server %s listening on port %s", s.serverID, *portFlag)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
