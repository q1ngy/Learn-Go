package main

import "fmt"

func main() {
	invoke()
}

func invoke() {
	//i := returnDefer()
	//i := returnDefer()

	f := adder()
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())

	// 下面的几个例子已经修复了，1.20版本之前会返回10
	//deferClosureLoop()
	//deferClosureLoop2()
	//deferClosureLoop3()
}

func returnDefer() int {
	i := 0
	defer func() {
		i = 1
	}()
	return i
}
func returnDefer2() (i int) {
	i = 0
	defer func() {
		i = 1
	}()
	return i
}

func deferClosureLoop() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

func deferClosureLoop2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}
func deferClosureLoop3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}

func adder() func() int {
	sum := 0
	return func() int {
		sum += 1
		return sum
	}
}
