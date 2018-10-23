package main

import (
	"fmt"
	"runtime"
	"sync"
)

// сюда писать код

var wg = &sync.WaitGroup{}

func testPipe(in, out chan interface{}) {
	defer wg.Done()
	recv := <-in
	if cnv, ok := recv.(int); ok {
		cnv++
		out <- cnv
		runtime.Gosched()
	}
}

func swapPipe(in, out chan interface{}) {
	defer wg.Done()
	fmt.Println("swap")
	recv := <-out
	in <- recv
	runtime.Gosched()
}

func ExecutePipeline(workers ...job) {
	in, out := make(chan interface{}, 1), make(chan interface{})
	for _, worker := range workers {
		wg.Add(1)
		go worker(in, out)
		wg.Add(1)
		go swapPipe(in, out)
	}
	in <- 1
}

func main() {
	workflow := []job{
		job(testPipe),
		job(testPipe),
		job(testPipe),
		job(testPipe),
		job(testPipe),
	}
	ExecutePipeline(workflow...)
	wg.Wait()
}
