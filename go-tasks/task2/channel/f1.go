package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	c := make(chan int)
	wg.Add(10)
	go func() {
		for i := 1; i <= 10; i++ {
			c <- i
		}
		close(c)
	}()
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Println(<-c)
			wg.Done()
		}
	}()
	wg.Wait()
}
