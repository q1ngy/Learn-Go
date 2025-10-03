package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID int
}

func (e *Employee) PrintInfo() {
	fmt.Printf("%+v", e)
}

func main() {
	e := Employee{Person: Person{Name: "slim", Age: 18}, EmployeeID: 10001}
	e.PrintInfo()
}
