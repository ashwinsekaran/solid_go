package main

import "fmt"

// LeetCode 992 - Subarrays with K Different Integers
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
