package task1

import (
	"strconv"
)

// https://leetcode-cn.com/problems/palindrome-number/
func isPalindrome(x int) bool {
	s := strconv.Itoa(x)
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-i-1] {
			return false
		}
	}
	return true
}
