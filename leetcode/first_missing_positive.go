package leetcode

// LeetCode 41 - First Missing Positive
//
// Pseudo code:
//   place each number n in nums[n-1] if 1 <= n <= len(nums)
//   scan array: first index i where nums[i] != i+1 -> return i+1
//   if all placed correctly: return n+1

func firstMissingPositive(nums []int) int {
	n := len(nums)
	for i := 0; i < n; i++ {
		for nums[i] > 0 && nums[i] <= n && nums[nums[i]-1] != nums[i] {
			nums[i], nums[nums[i]-1] = nums[nums[i]-1], nums[i]
		}
	}
	for i := 0; i < n; i++ {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return n + 1
}
