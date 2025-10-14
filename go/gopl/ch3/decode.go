package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "Hello, 世界"
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%d\t%c\n", i, r)
		i += size
	}

	for i, r := range s {
		fmt.Printf("%d\t%q\t%d\n", i, r, r)
	}

	// 统计字符数
	n := 0
	for range s {
		n++
	}
}
