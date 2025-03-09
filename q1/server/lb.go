package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math"
	"net"
	"os"
	"sort"
	"sync"

	"google.golang.org/grpc"
	pb "github.com/Mrinal92/Distri-Assignment2/q1/protofiles"
)

type serverInfo struct {
	address string
	load    float64
}

type loadBalancerServer struct {
	pb.UnimplementedLoadBalancerServiceServer
	mu           sync.Mutex
	servers      map[string]serverInfo
	roundRobinIx int
	policy       string
}

func newLoadBalancerServer(policy string) *loadBalancerServer {
	return &loadBalancerServer{
		servers: make(map[string]serverInfo),
		policy:  policy,
	}
}

// setupLogging configures logging to both a file and stdout.
func setupLogging() {
	logFile, err := os.OpenFile("lb.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file lb.log: %v", err)
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (lb *loadBalancerServer) RegisterServer(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.servers[req.GetServerId()] = serverInfo{
		address: req.GetAddress(),
		load:    0.0,
	}
	log.Printf("[LB] Registered server %s at address %s\n", req.GetServerId(), req.GetAddress())
	return &pb.RegisterResponse{
		Success: true,
		Message: "Server registered successfully",
	}, nil
}

func (lb *loadBalancerServer) ReportLoad(ctx context.Context, req *pb.LoadReportRequest) (*pb.LoadReportResponse, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if _, exists := lb.servers[req.GetServerId()]; !exists {
		log.Printf("[LB] Error: Server %s attempted to report load but is not registered.", req.GetServerId())
		return &pb.LoadReportResponse{
			Success: false,
			Message: "Server not registered",
		}, nil
	}

	si := lb.servers[req.GetServerId()]
	si.load = float64(req.GetCurrentLoad())
	lb.servers[req.GetServerId()] = si
	log.Printf("[LB] Updated load for server %s: %.2f", req.GetServerId(), si.load)
	return &pb.LoadReportResponse{
		Success: true,
		Message: "Load updated",
	}, nil
}

func (lb *loadBalancerServer) GetBestServer(ctx context.Context, req *pb.BestServerRequest) (*pb.BestServerResponse, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.servers) == 0 {
		log.Println("[LB] No servers available to process request.")
		return &pb.BestServerResponse{
			Success: false,
			Address: "",
			Message: "No servers available",
		}, nil
	}

	var address string
	switch lb.policy {
	case "pick_first":
		address = lb.pickFirst()
	case "round_robin":
		address = lb.roundRobin()
	case "least_load":
		address = lb.leastLoad()
	default:
		log.Printf("[LB] Unknown policy '%s', defaulting to pick_first", lb.policy)
		address = lb.pickFirst()
	}

	log.Printf("[LB] Policy=%s -> Assigned server: %s\n", lb.policy, address)
	return &pb.BestServerResponse{
		Success: true,
		Address: address,
		Message: "Success",
	}, nil
}

func (lb *loadBalancerServer) pickFirst() string {
	for _, si := range lb.servers {
		return si.address
	}
	return ""
}

func (lb *loadBalancerServer) roundRobin() string {
	keys := make([]string, 0, len(lb.servers))
	for k := range lb.servers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	serverID := keys[lb.roundRobinIx%len(keys)]
	lb.roundRobinIx = (lb.roundRobinIx + 1) % len(keys)
	log.Printf("[LB] Round Robin -> Selected server: %s", lb.servers[serverID].address)
	return lb.servers[serverID].address
}

func (lb *loadBalancerServer) leastLoad() string {
	minLoad := math.MaxFloat64
	var minServerID string
	for id, si := range lb.servers {
		if si.load < minLoad {
			minLoad = si.load
			minServerID = id
		}
	}
	log.Printf("[LB] Least Load -> Selected server: %s with load %.2f", lb.servers[minServerID].address, minLoad)
	return lb.servers[minServerID].address
}

var (
	portFlag   = flag.String("port", "9000", "Port to run the Load Balancer on")
	policyFlag = flag.String("policy", "round_robin", "Load balancing policy: pick_first, round_robin, least_load")
)

func main() {
	flag.Parse()
	setupLogging()

	lis, err := net.Listen("tcp", ":"+*portFlag)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *portFlag, err)
	}

	lbServer := newLoadBalancerServer(*policyFlag)
	grpcServer := grpc.NewServer()
	pb.RegisterLoadBalancerServiceServer(grpcServer, lbServer)

	log.Printf("[LB] Load Balancer listening on :%s with policy: %s", *portFlag, *policyFlag)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
