package main

import "fmt"

// LeetCode 3355 - Zero Array Transformation I
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
