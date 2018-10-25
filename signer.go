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
	close(out)
}

func MultiHash(in, out chan interface{}) {
	result := ""
	for recv := range in {
		if data, ok := recv.(string); ok {
			for th := 0; th < 6; th++ {
				result += DataSignerCrc32(strconv.Itoa(th) + data)
			}
		}
	}
	out <- result
	close(out)
	result = ""
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
	close(out)
}

func ExecutePipeline(workers ...job) {
	var wg = &sync.WaitGroup{}
	in, out := make(chan interface{}, 100), make(chan interface{}, 100)
	for _, work := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(in, out)
			close(out)
		}()
		in, out = out, make(chan interface{}, 100)
	}
	wg.Wait()
}

func main() {
	workflow := make([]job, 4)
	workflow[0] = job(func(in, out chan interface{}) {
		for i := 0; i < 2; i++ {
			fmt.Println(i)
			out <- i
		}
		close(out)
	})
	workflow[1] = job(SingleHash)
	workflow[2] = job(MultiHash)
	workflow[3] = job(CombineResults)
	ExecutePipeline(workflow...)
}
