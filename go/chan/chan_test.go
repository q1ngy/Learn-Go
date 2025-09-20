package chan_test

import (
	"fmt"
	"testing"
)

func TestChan(t *testing.T) {
	ints := make(chan []int, 10)
	fmt.Println(ints)
}

type User struct {
	name string
}
