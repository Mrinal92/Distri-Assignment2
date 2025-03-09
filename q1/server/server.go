package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/Mrinal92/Distri-Assignment2/q1/protofiles"
)

// setupLogging configures logging to both a file and stdout.
func setupLogging(serverID string) {
	filename := fmt.Sprintf("server_%s.log", serverID)
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", filename, err)
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type computeServer struct {
	pb.UnimplementedComputeServiceServer
	serverID    string
	lbAddress   string
	currentLoad int64 // Active tasks counter.
}

// ComputeTask simulates a computation.
// If task data starts with "heavy:" then it performs a CPU-intensive computation;
// otherwise, it simulates a regular task with a random sleep.
func (s *computeServer) ComputeTask(ctx context.Context, req *pb.ComputeRequest) (*pb.ComputeResponse, error) {
	atomic.AddInt64(&s.currentLoad, 1)
	defer atomic.AddInt64(&s.currentLoad, -1)

	log.Printf("[Server %s] Received task: %s\n", s.serverID, req.GetTaskData())
	if strings.HasPrefix(req.GetTaskData(), "heavy:") {
		// Simulate CPU-heavy task by computing a large sum.
		sum := 0
		for i := 0; i < 100000000; i++ {
			sum += i % 10
		}
		log.Printf("[Server %s] Heavy task computed sum: %d\n", s.serverID, sum)
	} else {
		// Regular task: sleep for a random duration (0 to 1 second).
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}

	result := fmt.Sprintf("Result from server %s for input '%s'", s.serverID, req.GetTaskData())
	log.Printf("[Server %s] Completed task: %s\n", s.serverID, req.GetTaskData())
	return &pb.ComputeResponse{Result: result}, nil
}

// registerWithLB registers the server with the Load Balancer.
func (s *computeServer) registerWithLB() {
	for {
		conn, err := grpc.Dial(s.lbAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("[Server %s] Could not connect to LB, retrying in 5 seconds...", s.serverID)
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()

		lbClient := pb.NewLoadBalancerServiceClient(conn)
		address := fmt.Sprintf("localhost:%s", *portFlag)
		_, err = lbClient.RegisterServer(context.Background(), &pb.RegisterRequest{
			ServerId: s.serverID,
			Address:  address,
		})
		if err != nil {
			log.Printf("[Server %s] RegisterServer failed: %v, retrying in 5 seconds...", s.serverID, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[Server %s] Successfully registered with LB at %s", s.serverID, s.lbAddress)
		break
	}
}

// reportLoadToLB periodically reports the current load to the Load Balancer.
func (s *computeServer) reportLoadToLB() {
	for {
		conn, err := grpc.Dial(s.lbAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("[Server %s] Could not connect to LB, retrying in 5 seconds...", s.serverID)
			time.Sleep(5 * time.Second)
			continue
		}
		lbClient := pb.NewLoadBalancerServiceClient(conn)
		loadVal := float32(atomic.LoadInt64(&s.currentLoad))
		_, err = lbClient.ReportLoad(context.Background(), &pb.LoadReportRequest{
			ServerId:    s.serverID,
			CurrentLoad: loadVal,
		})
		if err != nil {
			log.Printf("[Server %s] ReportLoad error: %v, retrying in 5 seconds...", s.serverID, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[Server %s] Reported load: %.2f", s.serverID, loadVal)
		conn.Close()
		time.Sleep(2 * time.Second)
	}
}

var (
	portFlag   = flag.String("port", "9100", "Port to run the server on")
	lbAddrFlag = flag.String("lb_addr", "localhost:9000", "Address of the load balancer")
	idFlag     = flag.String("id", "server1", "Unique server ID")
)

func main() {
	flag.Parse()
	setupLogging(*idFlag)

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
	log.Printf("[Server %s] Backend Server listening on port %s", s.serverID, *portFlag)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[Server %s] Failed to serve: %v", s.serverID, err)
	}
}
