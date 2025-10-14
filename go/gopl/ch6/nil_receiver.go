package main

import "fmt"

type M map[string]string

func (m M) get(key string) string {
	if m == nil {
		return "nil"
	}
	return m[key]
}

func main() {
	m := M{"a": "aaa"}
	fmt.Println(m.get("a"))
	m = nil
	fmt.Println(m.get("a"))
}
