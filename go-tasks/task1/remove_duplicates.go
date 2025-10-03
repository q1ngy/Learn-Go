package task1

// https://leetcode.cn/problems/remove-duplicates-from-sorted-array/description/
func removeDuplicates(nums []int) int {
	i := 0
	j := 1
	for j < len(nums) {
		if nums[i] == nums[j] {
			j++
		} else {
			i++
			nums[i] = nums[j]
		}
	}
	return i + 1
}
