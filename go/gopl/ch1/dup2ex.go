package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines2(os.Stdin, counts)
	} else {
		for _, file := range files {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			b := false
			input := bufio.NewScanner(f)
			for input.Scan() {
				line := input.Text()
				if _, ok := counts[line]; ok {
					b = true
					break
				}
				counts[line]++
			}
			if b {
				fmt.Println("filename: " + f.Name())
			}
			f.Close()
		}
	}
}

func countLines2(in *os.File, counts map[string]int) {
	input := bufio.NewScanner(in)
	for input.Scan() {
		counts[input.Text()]++
	}
}
