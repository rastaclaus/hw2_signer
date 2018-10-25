package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	for recv := range in {
		if recv, ok := recv.(int); ok {
			data := strconv.Itoa(recv)
			out <- (DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data)))
		}
	}
}

func MultiHash(in, out chan interface{}) {
	for recv := range in {
		if data, ok := recv.(string); ok {
			result := ""
			for th := 0; th < 6; th++ {
				result += DataSignerCrc32(strconv.Itoa(th) + data)
			}
			out <- result
		}
	}
}

func CombineResults(in, out chan interface{}) {
	ss := make([]string, 0)
	for recv := range in {
		if data, ok := recv.(string); ok {
			ss = append(ss, data)
		}
	}
	sort.Strings(ss)
	out <- strings.Join(ss, "_")
}

func do_work(worker job, wg *sync.WaitGroup, in, out chan interface{}) {
	defer wg.Done()
	defer close(out)
	worker(in, out)
}

func ExecutePipeline(workers ...job) {
	var wg = &sync.WaitGroup{}
	channels := make([]chan interface{}, len(workers)+1)
	for i, _ := range channels {
		channels[i] = make(chan interface{})
	}
	for i, worker := range workers {
		in, out := channels[i], channels[i+1]
		wg.Add(1)
		go do_work(worker, wg, in, out)
	}
	wg.Wait()
}

func main() {
	workflow := make([]job, 5)
	workflow[0] = job(func(in, out chan interface{}) {
		for i := 0; i < 2; i++ {
			out <- i
		}
	})
	workflow[1] = job(SingleHash)
	workflow[2] = job(MultiHash)
	workflow[3] = job(CombineResults)
	workflow[4] = job(func(in, out chan interface{}) {
		for res := range in {
			fmt.Println(res)
		}
	})
	ExecutePipeline(workflow...)
}
