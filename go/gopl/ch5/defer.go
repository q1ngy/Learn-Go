package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 100; i++ {
		defer fmt.Println(i)
		fmt.Println("start")
		func() {
			defer fmt.Println(i)
		}()
	}
}
