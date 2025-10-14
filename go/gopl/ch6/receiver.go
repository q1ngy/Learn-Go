package main

import (
	"fmt"
	"time"
)

type H interface {
	h()
}
type Student struct {
}

func (s *Student) h() {
	fmt.Println("hello")
}

func (s Student) f() {
	fmt.Println("method expression")
}

func main() {
	var h H
	h = &Student{}
	h.h()
	f := h.h
	time.AfterFunc(time.Second*3, f)
	time.Sleep(5 * time.Second)

	h2 := Student.f
	h2(Student{})
}
