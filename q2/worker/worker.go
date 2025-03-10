// worker.go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	pb "github.com/Mrinal92/Distri-Assignment2/q2/protofiles"
)

// KeyValue holds a key/value pair.
type KeyValue struct {
	Key   string
	Value string
}

// Global job type variable.
var jobType string

// wordRegex matches words that consist of letters and digits and allow an internal apostrophe.
// It captures contractions (e.g. "couldn’t") as one token.
var wordRegex = regexp.MustCompile(`[A-Za-z0-9]+(?:['’][A-Za-z0-9]+)*`)

// ----------------- Map Functions -----------------

// wordCountMap tokenizes the file using wordRegex and emits (word, "1") pairs.
func wordCountMap(filename string, contents string) []KeyValue {
	var kva []KeyValue
	words := wordRegex.FindAllString(contents, -1)
	for _, w := range words {
		w = strings.ToLower(w)
		kva = append(kva, KeyValue{Key: w, Value: "1"})
	}
	return kva
}

// invertedIndexMap tokenizes the file using wordRegex and emits (word, filename) pairs.
func invertedIndexMap(filename string, contents string) []KeyValue {
	var kva []KeyValue
	words := wordRegex.FindAllString(contents, -1)
	for _, w := range words {
		w = strings.ToLower(w)
		kva = append(kva, KeyValue{Key: w, Value: filename})
	}
	return kva
}

// ----------------- Reduce Functions -----------------

// wordCountReduce sums counts for a word.
func wordCountReduce(key string, values []string) string {
	sum := 0
	for _, v := range values {
		n, _ := strconv.Atoi(v)
		sum += n
	}
	return strconv.Itoa(sum)
}

// invertedIndexReduce aggregates filenames for a word and deduplicates them.
func invertedIndexReduce(key string, values []string) string {
	fileSet := make(map[string]bool)
	for _, filename := range values {
		fileSet[filename] = true
	}
	var files []string
	for f := range fileSet {
		files = append(files, f)
	}
	sort.Strings(files)
	return strings.Join(files, ",")
}

// ----------------- Utility Functions -----------------

// ihash returns a non-negative hash for a given key.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// ----------------- RPC Helper Functions -----------------

// registerWorker registers this worker with the master using the provided worker ID.
func registerWorker(client pb.MasterServiceClient, workerID string) (*pb.RegisterWorkerResponse, error) {
	req := &pb.RegisterWorkerRequest{WorkerAddress: workerID}
	return client.RegisterWorker(context.Background(), req)
}

// callGetTask requests a task from the master.
func callGetTask(client pb.MasterServiceClient, workerID string) *pb.Task {
	req := &pb.GetTaskRequest{WorkerId: workerID}
	resp, err := client.GetTask(context.Background(), req)
	if err != nil {
		log.Fatalf("GetTask error: %v", err)
	}
	return resp.Task
}

// callReportTask notifies the master that a task is complete.
func callReportTask(client pb.MasterServiceClient, workerID string, taskType pb.TaskType, taskID int32) {
	req := &pb.ReportTaskRequest{
		WorkerId: workerID,
		TaskType: taskType,
		TaskId:   taskID,
	}
	_, err := client.ReportTask(context.Background(), req)
	if err != nil {
		log.Fatalf("ReportTask error: %v", err)
	}
}

// ----------------- Task Execution Functions -----------------

// doMapTask reads the input file, applies the appropriate map function,
// partitions the output into n_reduce buckets, and writes intermediate files
// (named "mr-<mapTaskID>-<reduceID>") in the current directory.
func doMapTask(task *pb.Task) {
	filename := task.InputFile
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}

	var kva []KeyValue
	if jobType == "wordcount" {
		kva = wordCountMap(filename, string(data))
	} else if jobType == "invertedindex" {
		kva = invertedIndexMap(filename, string(data))
	} else {
		log.Fatalf("Unknown job type: %s", jobType)
	}

	nReduce := int(task.NReduce)
	buckets := make([][]KeyValue, nReduce)
	for _, kv := range kva {
		index := ihash(kv.Key) % nReduce
		buckets[index] = append(buckets[index], kv)
	}
	for i := 0; i < nReduce; i++ {
		oname := fmt.Sprintf("mr-%d-%d", task.TaskId, i)
		file, err := os.Create(oname)
		if err != nil {
			log.Fatalf("cannot create %v", oname)
		}
		enc := json.NewEncoder(file)
		for _, kv := range buckets[i] {
			if err := enc.Encode(&kv); err != nil {
				log.Fatalf("cannot encode kv pair: %v", err)
			}
		}
		file.Close()
	}
	log.Printf("MAP task %d completed", task.TaskId)
}

// doReduceTask reads intermediate files from the current directory,
// applies the appropriate reduce function, and writes the final output
// to "mr-out-<reduceID>".
func doReduceTask(task *pb.Task) {
	nMap := int(task.NMap)
	reduceID := int(task.ReduceId)
	var kva []KeyValue

	for i := 0; i < nMap; i++ {
		filename := fmt.Sprintf("mr-%d-%d", i, reduceID)
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			continue
		}
		dec := json.NewDecoder(strings.NewReader(string(data)))
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			kva = append(kva, kv)
		}
	}

	kvMap := make(map[string][]string)
	for _, kv := range kva {
		kvMap[kv.Key] = append(kvMap[kv.Key], kv.Value)
	}

	oname := fmt.Sprintf("mr-out-%d", reduceID)
	file, err := os.Create(oname)
	if err != nil {
		log.Fatalf("cannot create output file %v", oname)
	}

	for key, values := range kvMap {
		var output string
		if jobType == "wordcount" {
			output = wordCountReduce(key, values)
		} else if jobType == "invertedindex" {
			output = invertedIndexReduce(key, values)
		}
		fmt.Fprintf(file, "%v %v\n", key, output)
	}
	file.Close()
	log.Printf("REDUCE task %d completed", task.TaskId)
}

func main() {
	masterAddr := flag.String("master", "localhost:50051", "master address")
	jobTypeFlag := flag.String("job", "wordcount", "job type: wordcount or invertedindex")
	workerIDFlag := flag.String("worker", "worker-1", "unique worker ID")
	flag.Parse()
	jobType = *jobTypeFlag

	conn, err := grpc.Dial(*masterAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMasterServiceClient(conn)

	resp, err := registerWorker(client, *workerIDFlag)
	if err != nil {
		log.Fatalf("Worker registration failed: %v", err)
	}
	workerID := resp.WorkerId
	log.Printf("Worker registered with ID: %s (job type: %s)", workerID, jobType)

	masterStart := time.Unix(0, resp.MasterStartTime)
	delay := time.Duration(resp.RegDelay) * time.Second
	startPollTime := masterStart.Add(delay)
	now := time.Now()
	if startPollTime.After(now) {
		waitDuration := startPollTime.Sub(now)
		log.Printf("Waiting for %v before starting task polling...", waitDuration)
		time.Sleep(waitDuration)
	}

	for {
		task := callGetTask(client, workerID)
		switch task.TaskType {
		case pb.TaskType_WAIT:
			time.Sleep(time.Second)
		case pb.TaskType_EXIT:
			log.Println("No more tasks; exiting.")
			return
		case pb.TaskType_MAP:
			log.Printf("Received MAP task: %d", task.TaskId)
			doMapTask(task)
			callReportTask(client, workerID, pb.TaskType_MAP, task.TaskId)
		case pb.TaskType_REDUCE:
			log.Printf("Received REDUCE task: %d", task.TaskId)
			doReduceTask(task)
			callReportTask(client, workerID, pb.TaskType_REDUCE, task.TaskId)
		}
	}
}
