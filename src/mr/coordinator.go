package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Coordinator struct {
	// Your definitions here.
	mu             sync.Mutex
	currentMapFile int

	files   []string
	nReduce int

	jobs chan MapReduceArgs
	done chan bool

	doneMapTasks    map[string]bool
	doneReduceTasks map[string]bool
}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) GetNextTask(arg *MapReduceTaskArgs, reply *MapReduceTaskReply) error {
	c.mu.Lock()
	c.currentMapFile += 1
	currentMapFile := c.currentMapFile
	c.mu.Unlock()

	if c.currentMapFile >= len(c.files) {
		return nil
		// c.replyReduceTask(arg, reply)
	}
	fmt.Printf("Sending file: %s\n", c.files[currentMapFile])
	reply.Filename = c.files[currentMapFile]
	reply.TaskType = Map
	return nil
}

func (c *Coordinator) MarkMapTaskDone(args *MapTaskDoneArgs, reply *MapTaskDoneReply) error {
	c.mu.Lock()
	c.doneMapTasks[args.InputFilename] = true
	c.mu.Unlock()

	

	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	select {
	case <-c.done:
		ret = true
	default:
		ret = false
	}

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.files = files
	c.currentMapFile = -1
	c.nReduce = nReduce
	c.server()
	return &c
}
