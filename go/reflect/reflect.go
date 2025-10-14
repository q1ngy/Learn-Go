package main

import (
	"fmt"
	"reflect"
)

func main() {
	a := 1
	rt := reflect.TypeOf(a)
	fmt.Println(rt)
	fmt.Println(rt.Kind())
	rv := reflect.ValueOf(a)
	fmt.Println(rv)
}
