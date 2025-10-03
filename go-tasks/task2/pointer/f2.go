package main

import (
	"fmt"
	"sync"
	"time"
)

type Task func()
type TaskResult struct {
	d time.Duration
}

func runTasks(tasks []Task) []TaskResult {
	var wg sync.WaitGroup
	result := make([]TaskResult, len(tasks))
	for i, t := range tasks {
		wg.Add(1)
		go func(index int, task Task) {
			defer wg.Done()
			start := time.Now()
			task()
			end := time.Now()
			d := end.Sub(start)
			ts := TaskResult{d: d}
			result[index] = ts
		}(i, t)
	}
	wg.Wait()
	return result
}

func main() {
	tasks := []Task{
		func() {
			time.Sleep(1 * time.Second)
		},
		func() {
			time.Sleep(2 * time.Second)
		},
		func() {
			time.Sleep(3 * time.Second)
		},
	}
	results := runTasks(tasks)
	for _, t := range results {
		fmt.Println(t.d)
	}
}
