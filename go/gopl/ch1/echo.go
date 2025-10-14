package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(os.Args[0])
	//fmt.Println(os.Args[1])

	var s, sep string
	for i := 0; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)

	s, sep = "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)

	fmt.Println(strings.Join(os.Args[1:], " "))

	// exercise 1
	s, sep = "", ""
	for _, arg := range os.Args[:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)

	// exercise 2
	s = ""
	for i, arg := range os.Args[1:] {
		fmt.Printf("index: %d, arg: %s\n", i, arg)
	}
	fmt.Println(s)

}
