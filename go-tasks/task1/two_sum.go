package task1

// https://leetcode.cn/problems/two-sum/
func twoSum(nums []int, target int) []int {
	j := 1
	for j < len(nums) {
		i := 0
		for i < j {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
			i++
		}
		j++
	}
	return nil
}
