package task1

// https://leetcode.cn/problems/longest-common-prefix/description/
func longestCommonPrefix(strs []string) string {
	res := strs[0]
	for i := 1; i < len(strs); i++ {
		str := strs[i]
		for j, c := range res {
			if j > len(str)-1 || c != rune(str[j]) {
				res = str[:j]
				break
			} else {
				continue
			}
		}
	}
	return res
}
