package task1

import (
	"slices"
)

// https://leetcode.cn/problems/merge-intervals/
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}
	slices.SortFunc(intervals, func(a, b []int) int {
		return a[0] - b[0]
	})
	var res [][]int
	for i := 0; i < len(intervals); i++ {
		l, r := intervals[i][0], intervals[i][1]
		if len(res) == 0 || l > res[len(res)-1][1] {
			res = append(res, []int{l, r})
		} else {
			if res[len(res)-1][1] < r {
				res[len(res)-1][1] = r
			}
		}
	}
	return res
}
