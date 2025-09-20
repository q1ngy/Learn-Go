package main

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	name string
	age  int
}

func main() {
	//basic()
	//cancel()

	// deadline
	context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	// timeout
	context.WithTimeout(context.Background(), 2*time.Second)
}

func deadline() {

}

func cancel() {
	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		defer close(done)
		data, err := getRemoteData(ctx)
		fmt.Println(data, err)
	}()
	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()
	<-done
}

func basic() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "a", &User{
		name: "slim",
		age:  18,
	})
	getUser(ctx)
}

func getUser(ctx context.Context) {
	fmt.Println(ctx.Value("a").(*User).name)
}

func getRemoteData(ctx context.Context) (data string, err error) {
	select {
	case <-time.After(4 * time.Second):
		fmt.Println("remote data fetched")
		return "data", nil
	case <-ctx.Done():
		fmt.Println("goroutine canceled")
		return "", ctx.Err()
	}
}
