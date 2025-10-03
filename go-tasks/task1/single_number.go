package task1

// https://leetcode.cn/problems/single-number/
func singleNumber(nums []int) int {
	var m map[int]int = make(map[int]int)
	for _, num := range nums {
		if _, ok := m[num]; ok {
			m[num] = 2
		} else {
			m[num] = 1
		}
	}
	for k, v := range m {
		if v == 1 {
			return k
		}
	}
	return 0
}
