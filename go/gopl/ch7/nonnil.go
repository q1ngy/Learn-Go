package main

import "fmt"

type User struct {
	name string
}

type Run interface {
	Run()
}

func (u *User) Run() {
	fmt.Println(u.name)
}

func main() {
	var u *User = nil
	var r Run = u
	fmt.Println("u == nil ?", u == nil)
	fmt.Println("r == nil ?", r == nil)
	u.Run()
}
