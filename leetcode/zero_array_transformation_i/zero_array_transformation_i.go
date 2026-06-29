package main

import "fmt"

// LeetCode 3355 - Zero Array Transformation I
//
// Problem:
//   Given an integer array nums and a list of queries [l, r], each query
//   decrements every element in nums[l..r] by 1. Return true if nums can
//   become an all-zero array after processing all queries.
//
// Example 1:
//   Input:  nums = [1, 0, 1], queries = [[0,2]]
//   Output: true
//   Explanation: After query [0,2]: nums = [0, -1, 0] → all ≤ 0, so yes.
//
// Example 2:
//   Input:  nums = [4, 3, 2, 1], queries = [[1,3],[0,2]]
//   Output: false
//   Explanation:
//   After both queries, index 0 receives only 1 decrement but nums[0]=4 → cannot reach 0.
//
// Pseudo code:
//   difference array: for each query [l,r] do diff[l]--, diff[r+1]++
//   compute prefix sum of diff
//   at each index i: if nums[i] + prefix[i] > 0, cannot zero it → false
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

func main() {
	fmt.Println(isZeroArray([]int{1, 0, 1}, [][]int{{0, 2}}))              // true
	fmt.Println(isZeroArray([]int{4, 3, 2, 1}, [][]int{{1, 3}, {0, 2}}))  // false
	fmt.Println(isZeroArray([]int{0, 0, 0}, [][]int{}))                    // true
}
