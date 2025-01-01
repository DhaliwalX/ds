package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type LastResult struct {
	seq   int64
	value string
}

type KVServer struct {
	mu sync.Mutex

	// Your definitions here.

	data map[string]string

	lastResultPerClient map[int64]*LastResult
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	// Your code here.

	if data, ok := kv.data[args.Key]; ok {
		reply.Value = data
	}
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	lastResult, ok := kv.lastResultPerClient[args.Id]
	if ok && (lastResult.seq >= args.Seq) {
		return
	}

	kv.data[args.Key] = args.Value
	reply.Value = args.Value
	kv.lastResultPerClient[args.Id] = &LastResult{seq: args.Seq, value: reply.Value}
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.

	kv.mu.Lock()
	defer kv.mu.Unlock()

	lastResult, ok := kv.lastResultPerClient[args.Id]
	if ok && (lastResult.seq == args.Seq) {
		reply.Value = lastResult.value
		return
	}

	if _, ok := kv.data[args.Key]; !ok {
		kv.data[args.Key] = args.Value
		kv.lastResultPerClient[args.Id] = &LastResult{seq: args.Seq, value: reply.Value}
	} else {
		reply.Value = kv.data[args.Key]
		kv.data[args.Key] += args.Value
		kv.lastResultPerClient[args.Id] = &LastResult{seq: args.Seq, value: reply.Value}
	}
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	// You may need initialization code here.
	kv.data = make(map[string]string)
	kv.lastResultPerClient = make(map[int64]*LastResult)

	return kv
}
