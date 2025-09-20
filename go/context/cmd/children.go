package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()

	parent, cancel := context.WithCancel(ctx)
	child1, _ := context.WithCancel(parent)
	child2, _ := context.WithCancel(parent)
	grandChild, _ := context.WithCancel(child1)

	go worker(parent, "parent")
	go worker(child1, "child1")
	go worker(child2, "child2")
	go worker(grandChild, "grandChild")

	time.Sleep(2 * time.Second)
	fmt.Println(">>> calling cancelParent()")
	cancel()
	time.Sleep(1 * time.Second)
}

func worker(ctx context.Context, name string) {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println(name, "work finished")
	case <-ctx.Done():
		fmt.Println(name, "canceled", ctx.Err())
	}
}
