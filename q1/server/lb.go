package main

import (
	"context"
	"log"
	"math"
	"net"
	"sync"

	"google.golang.org/grpc"
	pb "github.com/Mrinal92/Distri-Assignment2/q1/protofiles"
)

// serverInfo holds registration and load details for a backend server.
type serverInfo struct {
	address string
	load    float64
}

// loadBalancerServer implements pb.LoadBalancerServiceServer.
type loadBalancerServer struct {
	pb.UnimplementedLoadBalancerServiceServer

	mu           sync.Mutex
	servers      map[string]serverInfo // map of server_id -> serverInfo
	roundRobinIx int
}

func newLoadBalancerServer() *loadBalancerServer {
	return &loadBalancerServer{
		servers: make(map[string]serverInfo),
	}
}

// RegisterServer is called by a backend server to register itself.
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

// ReportLoad is called periodically by a backend server to update its load.
func (lb *loadBalancerServer) ReportLoad(ctx context.Context, req *pb.LoadReportRequest) (*pb.LoadReportResponse, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if _, exists := lb.servers[req.GetServerId()]; !exists {
		return &pb.LoadReportResponse{
			Success: false,
			Message: "Server not registered",
		}, nil
	}

	si := lb.servers[req.GetServerId()]
	si.load = float64(req.GetCurrentLoad())
	lb.servers[req.GetServerId()] = si

	log.Printf("[LB] Updated load for server %s: %.2f\n", req.GetServerId(), si.load)
	return &pb.LoadReportResponse{
		Success: true,
		Message: "Load updated",
	}, nil
}

// GetBestServer selects a backend server based on the requested policy.
func (lb *loadBalancerServer) GetBestServer(ctx context.Context, req *pb.BestServerRequest) (*pb.BestServerResponse, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.servers) == 0 {
		return &pb.BestServerResponse{
			Success: false,
			Address: "",
			Message: "No servers available",
		}, nil
	}

	policy := req.GetPolicy()
	var address string

	switch policy {
	case "pick_first":
		// Return the first available server.
		for _, si := range lb.servers {
			address = si.address
			break
		}
	case "round_robin":
		// Use a round-robin strategy.
		keys := make([]string, 0, len(lb.servers))
		for k := range lb.servers {
			keys = append(keys, k)
		}
		serverID := keys[lb.roundRobinIx%len(keys)]
		lb.roundRobinIx++
		address = lb.servers[serverID].address
	case "least_load":
		// Find the server with the minimum load.
		minLoad := math.MaxFloat64
		var minServerID string
		for id, si := range lb.servers {
			if si.load < minLoad {
				minLoad = si.load
				minServerID = id
			}
		}
		address = lb.servers[minServerID].address
	default:
		// Default to pick_first if the policy is unknown.
		for _, si := range lb.servers {
			address = si.address
			break
		}
	}

	log.Printf("[LB] Policy=%s -> Chosen server: %s\n", policy, address)
	return &pb.BestServerResponse{
		Success: true,
		Address: address,
		Message: "Success",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	lbServer := newLoadBalancerServer()

	pb.RegisterLoadBalancerServiceServer(grpcServer, lbServer)

	log.Printf("Load Balancer listening on :9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
