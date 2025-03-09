// master.go
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	//pb "mapreduce/proto" // adjust this import path to where your generated proto files are located
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
	phase       string // "map" or "reduce"
	jobType     string // "wordcount" or "invertedindex"
}

// RegisterWorker registers a worker and returns its ID.
func (m *Master) RegisterWorker(ctx context.Context, req *pb.RegisterWorkerRequest) (*pb.RegisterWorkerResponse, error) {
	workerID := req.WorkerAddress
	log.Printf("Registered worker: %s", workerID)
	return &pb.RegisterWorkerResponse{WorkerId: workerID}, nil
}

// GetTask assigns an idle task to the requesting worker (or asks it to wait).
func (m *Master) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	workerID := req.WorkerId

	if m.phase == "map" {
		// Look for an idle map task.
		for i, task := range m.mapTasks {
			if task.State == Idle {
				m.mapTasks[i].State = InProgress
				m.mapTasks[i].Worker = workerID
				m.mapTasks[i].StartTime = time.Now()
				log.Printf("Assigning MAP task %d (file: %s) to worker %s", task.TaskID, task.InputFile, workerID)
				return &pb.GetTaskResponse{Task: &pb.Task{
					TaskType:  pb.TaskType_MAP,
					TaskId:    task.TaskID,
					InputFile: task.InputFile,
					NReduce:   task.NReduce,
				}}, nil
			}
		}
		// Check if all map tasks are completed.
		allDone := true
		for _, task := range m.mapTasks {
			if task.State != Completed {
				allDone = false
				break
			}
		}
		if allDone {
			// Transition to reduce phase.
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
			// Ask the worker to wait.
			return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_WAIT}}, nil
		}
	}

	if m.phase == "reduce" {
		// Look for an idle reduce task.
		for i, task := range m.reduceTasks {
			if task.State == Idle {
				m.reduceTasks[i].State = InProgress
				m.reduceTasks[i].Worker = workerID
				m.reduceTasks[i].StartTime = time.Now()
				log.Printf("Assigning REDUCE task %d to worker %s", task.TaskID, workerID)
				return &pb.GetTaskResponse{Task: &pb.Task{
					TaskType: pb.TaskType_REDUCE,
					TaskId:   task.TaskID,
					ReduceId: task.ReduceID,
					NMap:     int32(m.nMap),
				}}, nil
			}
		}
		// Check if all reduce tasks are completed.
		allDone := true
		for _, task := range m.reduceTasks {
			if task.State != Completed {
				allDone = false
				break
			}
		}
		if allDone {
			return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_EXIT}}, nil
		}
		return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_WAIT}}, nil
	}

	return &pb.GetTaskResponse{Task: &pb.Task{TaskType: pb.TaskType_WAIT}}, nil
}

// ReportTask marks a task as completed by the worker.
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
				log.Printf("MAP task %d completed by worker %s", taskID, workerID)
				break
			}
		}
	} else if taskType == pb.TaskType_REDUCE {
		for i, task := range m.reduceTasks {
			if task.TaskID == taskID && task.Worker == workerID {
				m.reduceTasks[i].State = Completed
				log.Printf("REDUCE task %d completed by worker %s", taskID, workerID)
				break
			}
		}
	}
	return &pb.ReportTaskResponse{Ack: true}, nil
}

func main() {
	// Command-line flags.
	numReduce := flag.Int("r", 3, "number of reduce tasks")
	jobType := flag.String("job", "wordcount", "job type: wordcount or invertedindex")
	flag.Parse()
	inputFiles := flag.Args()
	if len(inputFiles) == 0 {
		log.Fatalf("Usage: master -r <numReduce> -job <jobType> <inputfile1> <inputfile2> ...")
	}

	master := &Master{
		nMap:    len(inputFiles),
		nReduce: *numReduce,
		phase:   "map",
		jobType: *jobType,
	}
	// Each input file is assigned a separate map task.
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
