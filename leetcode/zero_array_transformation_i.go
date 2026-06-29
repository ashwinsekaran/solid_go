package leetcode

// LeetCode 3355 - Zero Array Transformation I
//
// Pseudo code:
//   use difference array to apply range decrements
//   diff[l] -= 1, diff[r+1] += 1 for each query
//   compute prefix sum of diff
//   for each index i: if nums[i] + prefix[i] > 0, return false
//   return true

func isZeroArray(nums []int, queries [][]int) bool {
	n := len(nums)
	diff := make([]int, n+1)
	for _, q := range queries {
		diff[q[0]]--
		diff[q[1]+1]++
	}
	prefix := 0
	for i := 0; i < n; i++ {
		prefix += diff[i]
		if nums[i]+prefix > 0 {
			return false
		}
	}
	return true
}
