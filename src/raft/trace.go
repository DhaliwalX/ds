package raft

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Slice struct {
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
}

type Trace struct {
	mu     sync.Mutex
	slices []Slice
}

func (trace *Trace) Trace(since time.Time, format string, args ...interface{}) {
	slice := Slice{Start: since,
		End:         time.Now(),
		Description: fmt.Sprintf(format, args...)}
	trace.mu.Lock()
	defer trace.mu.Unlock()
	trace.slices = append(trace.slices, slice)
}

func (trace *Trace) Dump() error {
	file, err := os.OpenFile("trace.out", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	log.Printf("Generated %d slices", len(trace.slices))
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(trace.slices)
}

var trace Trace

func T(since time.Time, format string, args ...interface{}) {
	trace.Trace(since, format, args...)
}

func D() error {
	return trace.Dump()
}
