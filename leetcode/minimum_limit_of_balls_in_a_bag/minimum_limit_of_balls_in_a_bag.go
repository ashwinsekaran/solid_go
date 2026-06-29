package main

import "fmt"

// LeetCode 1760 - Minimum Limit of Balls in a Bag
//
// Problem:
//   You have bags of balls. In one operation you can split a bag into two bags.
//   After at most maxOperations splits, minimise the maximum number of balls
//   in any single bag (the "penalty").
//
// Example 1:
//   Input:  nums = [9], maxOperations = 2
//   Output: 3
//   Explanation:
//   Split 9 → [6,3] (1 op), then [3,3,3] (2 ops). Max = 3.
//
// Example 2:
//   Input:  nums = [2,4,8,2], maxOperations = 4
//   Output: 2
//   Explanation: Split 8 twice and 4 once. Every bag has ≤ 2 balls.
//
// Pseudo code:
//   binary search on the answer (max balls allowed in any bag)
//   for a given max size: ops needed = sum of ceil(n/max)-1 for each bag
//   find smallest max where ops <= maxOperations

func minimumSize(nums []int, maxOperations int) int {
	lo, hi := 1, 1_000_000_000
	for lo < hi {
		mid := lo + (hi-lo)/2
		ops := 0
		for _, n := range nums {
			ops += (n - 1) / mid
		}
		if ops <= maxOperations {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func main() {
	fmt.Println(minimumSize([]int{9}, 2))              // 3
	fmt.Println(minimumSize([]int{2, 4, 8, 2}, 4))    // 2
	fmt.Println(minimumSize([]int{7, 17}, 2))          // 7
}
