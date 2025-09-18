package main

import "fmt"

type Number interface {
	~int | uint
}

func Max[T Number](vals []T) T {
	if vals == nil || len(vals) == 0 {
		panic("illegal vals")
	}
	e := vals[0]
	for _, elem := range vals {
		if elem > e {
			e = elem
		}
	}
	return e
}

func Find[T any](vals []T, filter func(t T) bool) T {
	for _, val := range vals {
		if filter(val) {
			return val
		}
	}
	var t T
	return t
}

func Insert[T Number](arr []T, idx int, val T) []T {
	if idx < 0 || idx > len(arr) {
		panic("illegal vals")
	}
	arr = append(arr, val)
	for i := len(arr) - 1; i > idx; i-- {
		if i > idx {
			arr[i], arr[i-1] = arr[i-1], arr[i]
		}
	}
	return arr
}

func Insert2[T Number](arr []T, idx int, val T) []T {
	if idx < 0 || idx > len(arr) {
		panic("illegal vals")
	}
	arr = append(arr, val)
	for i := len(arr) - 1; i > idx; i-- {
		if i > idx {
			arr[i] = arr[i-1]
		}
	}
	arr[idx] = val
	return arr
}

func main() {
	arr := []int{1, 2, 7, 5, 3}
	i := Max(arr)
	fmt.Println(i)

	arr2 := Insert(arr, 3, 99)
	fmt.Println(arr2)

	arr3 := Insert2(arr, 2, 99)
	fmt.Println(arr3)
}
