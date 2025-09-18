package main

func main() {
	invoke()
}

func invoke() {
	//i := returnDefer()
	//i := returnDefer()
	//deferClosureLoop()
	//deferClosureLoop2()
	deferClosureLoop3()
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
