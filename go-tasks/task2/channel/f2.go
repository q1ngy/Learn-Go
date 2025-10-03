package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	bc := make(chan int, 100)
	wg.Add(100)

	go Producer(bc)
	go Consumer(bc)

	wg.Wait()
}

func Producer(bc chan int) {
	for i := 0; i < 100; i++ {
		bc <- i
	}
}

func Consumer(bc chan int) {
	//for m := range bc {
	//	fmt.Println(m)
	//	wg.Done()
	//}
	for {
		select {
		case m, ok := <-bc:
			if !ok {
				break
			}
			fmt.Println(m)
			wg.Done()
		}
	}
}
