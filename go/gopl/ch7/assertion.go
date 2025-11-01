package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var w io.Writer
	var rw io.ReadWriter
	var f *os.File
	w = rw
	st := os.Stdout
	w = os.Stdout
	rw = os.Stdout

	w = w.(*os.File)
	fmt.Println(w, st, f)
}
