package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var num int32 = 0
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				atomic.AddInt32(&num, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(num)
}
