package task1

import (
	"fmt"
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	s := []string{"ab", "a"}
	prefix := longestCommonPrefix(s)
	fmt.Println(prefix)
}

func TestPlusOne(t *testing.T) {
	slice := []int{9}
	one := plusOne(slice)
	fmt.Println(one)
}

func TestRemoveDuplicates(t *testing.T) {
	slice := []int{1, 1, 2}
	res := removeDuplicates(slice)
	fmt.Println(res)
}
