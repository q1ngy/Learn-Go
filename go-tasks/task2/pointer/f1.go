package main

import "fmt"

func main() {
	// f1
	num := 0
	f1(&num)
	fmt.Println(num)

	// f2
	s := []int{1, 2, 3}
	f2(&s)
	fmt.Println(s)
}

func f1(p *int) {
	*p = 10
}

func f2(s *[]int) {
	for i := 0; i < len(*s); i++ {
		(*s)[i] = (*s)[i] * 2
	}
}
