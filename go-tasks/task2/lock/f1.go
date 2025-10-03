package main

import (
	"fmt"
	"sync"
)

func main() {
	num := 0
	var lock sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				lock.Lock()
				num++
				lock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(num)
}
