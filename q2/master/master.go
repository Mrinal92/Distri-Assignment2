// master.go
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "github.com/Mrinal92/Distri-Assignment2/q2/protofiles"
)

// TaskState represents the state of a task.
type TaskState int

const (
	Idle TaskState = iota
	InProgress
	Completed
)

// TaskInfo holds information for a map or reduce task.
type TaskInfo struct {
	TaskType  pb.TaskType
	TaskID    int32
	InputFile string   // valid for map tasks
	ReduceID  int32    // valid for reduce tasks
	NReduce   int32    // total number of reduce tasks (for map tasks)
	State     TaskState
	Worker    string
	StartTime time.Time
}

// Master implements the MasterService gRPC server.
type Master struct {
	pb.UnimplementedMasterServiceServer

	mu          sync.Mutex
	mapTasks    []TaskInfo
	reduceTasks []TaskInfo
	nMap        int
	nReduce     int
	phase       string        // "map" or "reduce"
	jobType     string        // "wordcount" or "invertedindex"
	startTime   time.Time     // master start time
	regDelay    time.Duration // registration delay
	workerLoad  map[string]int
}

// GetRandomIdleTask selects a random idle task from a slice.
func (m *Master) GetRandomIdleTask(taskList []TaskInfo) (int, bool) {
	var idleIndices []int
	for i, task := range taskList {
		if task.State == Idle {
			idleIndices = append(idleIndices, i)
		}
	}
	if len(idleIndices) == 0 {
		return 0, false
	}
	idx := idleIndices[rand.Intn(len(idleIndices))]
	return idx, true
}

// minWorkerLoad returns the minimum load among all registered workers.
func (m *Master) minWorkerLoad() int {
	min := math.MaxInt64
	for _, load := range m.workerLoad {
		if load < min {
			min = load
		}
	}
	return min
}

// RegisterWorker registers a worker and returns its ID along with master start time and regDelay.
func (m *Master) RegisterWorker(ctx context.Context, req *pb.RegisterWorkerRequest) (*pb.RegisterWorkerResponse, error) {
	workerID := req.WorkerAddress
	m.mu.Lock()
	if m.workerLoad == nil {
		m.workerLoad = make(map[string]int)
	}
	m.workerLoad[workerID] = 0
	m.mu.Unlock()
	log.Printf("Registered worker: %s", workerID)
	return &pb.RegisterWorkerResponse{
		WorkerId:        workerID,
		MasterStartTime: m.startTime.UnixNano(),
		RegDelay:        int64(m.regDelay.Seconds()),
	}, nil
}

// GetTask assigns a randomly selected idle task to the requesting worker.
// It waits until the registration delay has elapsed, then only assigns a task if the worker's load equals the minimum load.
func (m *Master) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	workerID := req.WorkerId

	if time.Since(m.startTime) < m.regDelay {
		return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_WAIT}}, nil
	}

	if m.workerLoad[workerID] > m.minWorkerLoad() {
		return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_WAIT}}, nil
	}

	if m.phase == "map" {
		if idx, ok := m.GetRandomIdleTask(m.mapTasks); ok {
			m.mapTasks[idx].State = InProgress
			m.mapTasks[idx].Worker = workerID
			m.mapTasks[idx].StartTime = time.Now()
			task := m.mapTasks[idx]
			log.Printf("Assigning MAP task %d (file: %s) to worker %s", task.TaskID, task.InputFile, workerID)
			return &pb.GetTaskResponse{Task: &pb.Task{
				TaskType:  pb.TaskType_MAP,
				TaskId:    task.TaskID,
				InputFile: task.InputFile,
				NReduce:   task.NReduce,
			}}, nil
		}
		allDone := true
		for _, task := range m.mapTasks {
			if task.State != Completed {
				allDone = false
				break
			}
		}
		if allDone {
			m.phase = "reduce"
			m.reduceTasks = make([]TaskInfo, m.nReduce)
			for i := 0; i < m.nReduce; i++ {
				m.reduceTasks[i] = TaskInfo{
					TaskType: pb.TaskType_REDUCE,
					TaskID:   int32(i),
					ReduceID: int32(i),
					State:    Idle,
				}
			}
		} else {
			return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_WAIT}}, nil
		}
	}

	if m.phase == "reduce" {
		if idx, ok := m.GetRandomIdleTask(m.reduceTasks); ok {
			m.reduceTasks[idx].State = InProgress
			m.reduceTasks[idx].Worker = workerID
			m.reduceTasks[idx].StartTime = time.Now()
			task := m.reduceTasks[idx]
			log.Printf("Assigning REDUCE task %d to worker %s", task.TaskID, workerID)
			return &pb.GetTaskResponse{Task: &pb.Task{
				TaskType: pb.TaskType_REDUCE,
				TaskId:   task.TaskID,
				ReduceId: task.ReduceID,
				NMap:     int32(m.nMap),
			}}, nil
		}
		allDone := true
		for _, task := range m.reduceTasks {
			if task.State != Completed {
				allDone = false
				break
			}
		}
		if allDone {
			m.printTaskDistribution()
			mergeOutputs(m.nReduce)
			return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_EXIT}}, nil
		}
		return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_WAIT}}, nil
	}

	return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_WAIT}}, nil
}

// ReportTask marks a task as completed by the worker and increments its load.
func (m *Master) ReportTask(ctx context.Context, req *pb.ReportTaskRequest) (*pb.ReportTaskResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	workerID := req.WorkerId
	taskType := req.TaskType
	taskID := req.TaskId

	if taskType == pb.TaskType_MAP {
		for i, task := range m.mapTasks {
			if task.TaskID == taskID && task.Worker == workerID {
				m.mapTasks[i].State = Completed
				m.workerLoad[workerID]++
				log.Printf("MAP task %d completed by worker %s", taskID, workerID)
				break
			}
		}
	} else if taskType == pb.TaskType_REDUCE {
		for i, task := range m.reduceTasks {
			if task.TaskID == taskID && task.Worker == workerID {
				m.reduceTasks[i].State = Completed
				m.workerLoad[workerID]++
				log.Printf("REDUCE task %d completed by worker %s", taskID, workerID)
				break
			}
		}
	}
	return &pb.ReportTaskResponse{Ack: true}, nil
}

// mergeOutputs reads each mr-out file from the worker folder and merges their contents into a single file final.txt.
func mergeOutputs(nReduce int) {
	workerDir := filepath.Join("..", "worker")
	finalPath := filepath.Join(workerDir, "final.txt")
	finalFile, err := os.Create(finalPath)
	if err != nil {
		log.Fatalf("Failed to create final output file: %v", err)
	}
	defer finalFile.Close()

	writer := bufio.NewWriter(finalFile)
	for i := 0; i < nReduce; i++ {
		mrOutName := filepath.Join(workerDir, fmt.Sprintf("mr-out-%d", i))
		inFile, err := os.Open(mrOutName)
		if err != nil {
			log.Printf("Warning: could not open %s: %v", mrOutName, err)
			continue
		}
		scanner := bufio.NewScanner(inFile)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Fprintln(writer, line)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading %s: %v", mrOutName, err)
		}
		inFile.Close()
	}
	writer.Flush()
	log.Printf("Merged final output is written to %s", finalPath)
}

// printTaskDistribution prints a summary of task assignments.
func (m *Master) printTaskDistribution() {
	distribution := make(map[string]int)
	for _, t := range m.mapTasks {
		if t.Worker != "" {
			distribution[t.Worker]++
		}
	}
	for _, t := range m.reduceTasks {
		if t.Worker != "" {
			distribution[t.Worker]++
		}
	}
	log.Printf("Task Distribution Summary:")
	for worker, count := range distribution {
		log.Printf("Worker %s: %d tasks", worker, count)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	numReduce := flag.Int("r", 3, "number of reduce tasks")
	jobType := flag.String("job", "wordcount", "job type: wordcount or invertedindex")
	regDelay := flag.Int("regdelay", 10, "registration delay in seconds")
	flag.Parse()
	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		log.Fatalf("Usage: master -r <numReduce> -job <jobType> -regdelay <seconds> <inputfile1> <inputfile2> ...")
	}

	master := &Master{
		nMap:       len(inputFiles),
		nReduce:    *numReduce,
		phase:      "map",
		jobType:    *jobType,
		startTime:  time.Now(),
		regDelay:   time.Duration(*regDelay) * time.Second,
		workerLoad: make(map[string]int),
	}
	// Each input file is a separate map task.
	master.mapTasks = make([]TaskInfo, len(inputFiles))
	for i, file := range inputFiles {
		master.mapTasks[i] = TaskInfo{
			TaskType:  pb.TaskType_MAP,
			TaskID:    int32(i),
			InputFile: file,
			NReduce:   int32(*numReduce),
			State:     Idle,
		}
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMasterServiceServer(s, master)
	log.Printf("Master started on port 50051 with job type %s", master.jobType)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
