package main

import (
	"fmt"
	"sync"
)

type SharedObject struct {
	increment int
	mu        sync.Mutex
	wg        *sync.WaitGroup
}

func runThread(obj *SharedObject) {
	obj.mu.Lock()
	defer obj.mu.Unlock()

	obj.increment++
	obj.wg.Done()
}

func main() {
	fmt.Println("Hello world")

	obj := SharedObject{
		wg: &sync.WaitGroup{},
	}
	for i := 0; i < 10; i++ {
		obj.wg.Add(1)
		go runThread(&obj)
	}

	obj.wg.Wait()

	fmt.Println(obj.increment)

	TestChannels()
}
