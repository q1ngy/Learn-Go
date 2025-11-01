package main

import "fmt"

type Animal interface {
	Run()
	Eat()
}
type Cat struct {
}

func (c *Cat) Run() {

}
func (c *Cat) Eat() {
}

func main() {
	//var a Animal = Cat{}
	var b Animal = &Cat{}
	fmt.Println(b)
}
