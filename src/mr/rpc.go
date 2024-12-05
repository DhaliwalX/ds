package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.
type DoneArgs struct{}

type DoneReply struct {
	Result bool
}

type MapReduceArgs struct {
	Files []string
}

type TaskType int

const (
	Map TaskType = iota
	Reduce
)

type MapReduceTaskArgs struct{}

type MapReduceTaskReply struct {
	TaskType TaskType
	Filename string
	Contents string

	KvData []string
}

type MapTaskDoneArgs struct {
	InputFilename  string
	OutputFilename string
}

type MapTaskDoneReply struct{}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
