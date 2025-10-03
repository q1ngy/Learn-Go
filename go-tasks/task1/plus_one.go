package task1

// https://leetcode.cn/problems/plus-one/description/
func plusOne(digits []int) []int {
	var res []int
	carry := false
	for i := len(digits) - 1; i >= 0; i-- {
		num := digits[i]

		if i == len(digits)-1 {
			if num == 9 {
				num = 0
				carry = true
			} else {
				num = num + 1
			}
		} else {
			if num == 9 && carry {
				num = 0
				carry = true
			} else if carry {
				num = num + 1
				carry = false
			} else {
				carry = false
			}
		}

		res = append(res, num)
	}
	if carry {
		res = append(res, 1)
	}
	for i := 0; i < len(res)/2; i++ {
		res[i], res[len(res)-i-1] = res[len(res)-i-1], res[i]
	}
	return res
}
