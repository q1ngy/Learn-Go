package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)
	go even()
	go odd()
	wg.Wait()
}

func even() {
	defer wg.Done()
	for i := 2; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Printf("Even: %d\n", i)
		}
	}
}

func odd() {
	defer wg.Done()
	for i := 1; i < 10; i++ {
		if i%2 != 0 {
			fmt.Printf("Odd: %d\n", i)
		}
	}
}
