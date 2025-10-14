package main

import "fmt"

func main() {
	arr := [...]string{"a", "b", "c", "d", "e"}
	var a []func()

	for _, s := range arr {
		a = append(a, func() {
			fmt.Println(s)
		})
	}

	for _, f := range a {
		f()
	}
}
