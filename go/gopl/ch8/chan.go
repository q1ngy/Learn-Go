package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("working...")
			time.Sleep(time.Second)
		}
		done <- struct{}{}
	}()
	<-done
}
