package task1

// https://leetcode.cn/problems/valid-parentheses/
func isValid(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	var stack []byte
	for i := range s {
		if pairs[s[i]] > 0 {
			// 右括号
			// 栈顶不是与右括号匹配的左括号
			if len(stack) == 0 || stack[len(stack)-1] != pairs[s[i]] {
				return false
			}
			// 出栈
			stack = stack[:len(stack)-1]
		} else {
			// 左括号
			// 入栈
			stack = append(stack, s[i])
		}
	}
	return len(stack) == 0
}

func stackForGo() {
	s := "hello"
	var stack []rune
	for _, v := range s {
		stack = append(stack, v)
	}
	for len(stack) > 0 {
		_ = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}
