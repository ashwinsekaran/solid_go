package main

import "fmt"

// LeetCode 41 - First Missing Positive
//
// Problem:
//   Given an unsorted integer array nums, return the smallest missing positive integer.
//   Must run in O(n) time and O(1) extra space.
//
// Example 1:
//   Input:  [1, 2, 0]      Output: 3
//
// Example 2:
//   Input:  [3, 4, -1, 1]  Output: 2
//   Explanation: 1 is present, 2 is missing.
//
// Example 3:
//   Input:  [7, 8, 9, 11]  Output: 1
//
//   Key idea: use the array itself as a hash table.
//   Place each number n in position n-1 (swap into correct slot).
//   Then scan: first index i where nums[i] ≠ i+1 → answer is i+1.
//
// Pseudo code:
//   use the array itself as a hash: place number n at index n-1 (swap into place)
//   scan: first index i where nums[i] != i+1 → return i+1
//   if all placed: return n+1

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

func main() {
	fmt.Println(firstMissingPositive([]int{1, 2, 0}))       // 3
	fmt.Println(firstMissingPositive([]int{3, 4, -1, 1}))   // 2
	fmt.Println(firstMissingPositive([]int{7, 8, 9, 11, 12})) // 1
}
