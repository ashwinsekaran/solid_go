package main

import "fmt"

// LeetCode 992 - Subarrays with K Different Integers
//
// Problem:
//   Given an integer array nums and an integer k, return the number of
//   subarrays with exactly k different integers.
//
// Example 1:
//   Input:  nums = [1, 2, 1, 2, 3], k = 2
//   Output: 7
//   Explanation: Subarrays with exactly 2 distinct values:
//   [1,2], [2,1], [1,2], [2,3], [1,2,1], [2,1,2], [1,2,1,2]
//
// Example 2:
//   Input:  nums = [1, 2, 1, 3, 4], k = 3
//   Output: 3
//
//   Key trick: exactly(k) = atMost(k) - atMost(k-1)
//   atMost(k) counts subarrays with ≤ k distinct values via a sliding window.
//
// Pseudo code:
//   exactly(k) = atMost(k) - atMost(k-1)
//   atMost(k): sliding window; shrink left when distinct count exceeds k
//              add (r - l + 1) subarrays ending at r

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

func main() {
	fmt.Println(subarraysWithKDistinct([]int{1, 2, 1, 2, 3}, 2)) // 7
	fmt.Println(subarraysWithKDistinct([]int{1, 2, 1, 3, 4}, 3)) // 3
	fmt.Println(subarraysWithKDistinct([]int{1, 1, 1, 1}, 1))    // 10
}
