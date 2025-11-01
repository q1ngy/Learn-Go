package main

type T struct {
	X int
}

func (t *T) SetX(x int) {
	t.X = x
}

type MyInt int

// 指针方法
func (m *MyInt) Increment() {
	*m = *m + 1
}

func main() {
	var a T = T{X: 1} // addressable
	a.SetX(1)
	//T{}.SetX(1) not addressable

	var b MyInt = 10
	b.Increment() // addressable
	//MyInt(1).Increment()
}
