package main

import "fmt"

/*
实现删除切片特定下标元素的方法。

	要求一：能够实现删除操作就可以。
	要求二：考虑使用比较高性能的实现。
	要求三：改造为泛型方法
	要求四：支持缩容，并旦设计缩容机制。
*/
// ❌
func d[T any](arr []T, idx int) []T {
	if idx < 0 || idx >= len(arr) {
		panic("illegal idx")
	}
	newArr := make([]T, 0, len(arr)-1)
	newArr = append(newArr, arr[:idx]...)
	newArr = append(newArr, arr[idx+1:]...)
	return newArr
}

func d2[T any](arr []T, idx int) (res []T, val T) {
	if idx < 0 || idx >= len(arr) {
		panic("illegal idx")
	}
	val = arr[idx]
	for i := idx; i < len(arr)-1; i++ {
		arr[i] = arr[i+1]
	}
	arr = arr[:len(arr)-1]
	// Shrink
	return arr, val
}

func main() {
	arr := []int{0, 1, 2, 3, 4, 5}
	res, val := d2(arr, 3)
	fmt.Println(res, val)
}
