package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "github.com/Mrinal92/Distri-Assignment2/q1/protofiles"
)

var (
	lbAddrFlag   = flag.String("lb_addr", "localhost:9000", "Load Balancer address")
	policyFlag   = flag.String("policy", "pick_first", "Load balancing policy: pick_first | round_robin | least_load")
	requestsFlag = flag.Int("requests", 5, "Number of requests to send")
	concurrency  = flag.Int("concurrency", 1, "Number of concurrent client goroutines")
)

func main() {
	flag.Parse()

	connLB, err := grpc.Dial(*lbAddrFlag, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to LB: %v", err)
	}
	defer connLB.Close()

	lbClient := pb.NewLoadBalancerServiceClient(connLB)

	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			for r := 0; r < *requestsFlag; r++ {
				resp, err := lbClient.GetBestServer(context.Background(), &pb.BestServerRequest{
					Policy: *policyFlag,
				})
				if err != nil || !resp.GetSuccess() {
					log.Printf("[Client %d] Could not get best server: %v\n", clientID, err)
					continue
				}
				serverAddress := resp.GetAddress()

				connServer, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("[Client %d] Could not connect to server %s: %v\n", clientID, serverAddress, err)
					continue
				}
				computeClient := pb.NewComputeServiceClient(connServer)

				taskData := fmt.Sprintf("Task#%d_from_Client%d", r, clientID)
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				res, err := computeClient.ComputeTask(ctx, &pb.ComputeRequest{TaskData: taskData})
				cancel()
				connServer.Close()

				if err != nil {
					log.Printf("[Client %d] ComputeTask error: %v\n", clientID, err)
				} else {
					log.Printf("[Client %d] Response from server: %s\n", clientID, res.GetResult())
				}
				time.Sleep(500 * time.Millisecond)
			}
		}(i)
	}
	wg.Wait()
	log.Println("All client requests completed.")
}
