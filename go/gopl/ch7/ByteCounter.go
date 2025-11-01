package main

import (
	"fmt"
)

type ByteCounter int

func (c *ByteCounter) String() string {
	return "bytecount"
}

func (c *ByteCounter) Write(p []byte) (n int, err error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

func main() {
	var b ByteCounter
	b.Write([]byte("hello"))
	fmt.Println(b)

	fmt.Fprintf(&b, "hello %s", "world")
	fmt.Println(b)
}
