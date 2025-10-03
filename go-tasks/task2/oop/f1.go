package main

import "fmt"

type Shape interface {
	Area() int
	Perimeter() int
}

type Rectangle struct {
	Width  int
	Height int
}

func (r Rectangle) Area() int {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() int {
	return (r.Width + r.Height) * 2
}

type Circle struct {
}

func (c Circle) Area() (a int) {
	return a
}

func (c Circle) Perimeter() (p int) {
	return p
}

func main() {
	var s Shape
	s = Rectangle{Width: 5, Height: 3}
	area := s.Area()
	perimeter := s.Perimeter()
	fmt.Printf("Area: %d, Perimeter: %d", area, perimeter)

	s = Circle{}
	s.Area()
	s.Perimeter()
}
