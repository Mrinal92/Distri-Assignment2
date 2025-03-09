package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Mrinal92/Distri-Assignment2/q1/protofiles"
)

// setupLogging configures logging to both a file and stdout.
func setupLogging() {
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file client.log: %v", err)
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

var (
	lbAddrFlag   = flag.String("lb_addr", "localhost:9000", "Load Balancer address")
	policyFlag   = flag.String("policy", "pick_first", "Load balancing policy: pick_first | round_robin | least_load")
	requestsFlag = flag.Int("requests", 5, "Number of requests to send")
	concurrency  = flag.Int("concurrency", 1, "Number of concurrent client goroutines")
	taskTypeFlag = flag.String("task_type", "regular", "Type of task: regular or heavy")
)

func main() {
	flag.Parse()
	setupLogging()

	log.Printf("[Client] Connecting to Load Balancer at %s...\n", *lbAddrFlag)
	connLB, err := grpc.Dial(*lbAddrFlag, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[Client] Could not connect to Load Balancer: %v", err)
	}
	defer connLB.Close()
	lbClient := pb.NewLoadBalancerServiceClient(connLB)
	log.Println("[Client] Successfully connected to Load Balancer!")

	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			for r := 0; r < *requestsFlag; r++ {
				log.Printf("[Client %d] Requesting best server from LB (Policy: %s)\n", clientID, *policyFlag)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				resp, err := lbClient.GetBestServer(ctx, &pb.BestServerRequest{
					Policy: *policyFlag,
				})
				cancel()
				if err != nil || !resp.GetSuccess() {
					log.Printf("[Client %d] Error getting best server: %v\n", clientID, err)
					continue
				}
				serverAddress := resp.GetAddress()
				log.Printf("[Client %d] Assigned to server: %s\n", clientID, serverAddress)

				connServer, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("[Client %d] Could not connect to server %s: %v\n", clientID, serverAddress, err)
					continue
				}
				computeClient := pb.NewComputeServiceClient(connServer)

				// Prepare task data with the specified task type.
				taskData := fmt.Sprintf("%s:Task#%d_from_Client%d", *taskTypeFlag, r, clientID)
				log.Printf("[Client %d] Sending task: %s to server %s\n", clientID, taskData, serverAddress)
				reqCtx, reqCancel := context.WithTimeout(context.Background(), 3*time.Second)
				res, err := computeClient.ComputeTask(reqCtx, &pb.ComputeRequest{TaskData: taskData})
				reqCancel()
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
	log.Println("[Client] All client requests completed.")
}
