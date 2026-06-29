package leetcode

// LeetCode 992 - Subarrays with K Different Integers
//
// Pseudo code:
//   answer = atMostK(nums, k) - atMostK(nums, k-1)
//   atMostK: sliding window, shrink left when distinct count > k
//             add (r - l + 1) to count for each right pointer

func subarraysWithKDistinct(nums []int, k int) int {
	return atMostK(nums, k) - atMostK(nums, k-1)
}

func atMostK(nums []int, k int) int {
	count := make(map[int]int)
	l, res := 0, 0
	for r, v := range nums {
		count[v]++
		for len(count) > k {
			count[nums[l]]--
			if count[nums[l]] == 0 {
				delete(count, nums[l])
			}
			l++
		}
		res += r - l + 1
	}
	return res
}
