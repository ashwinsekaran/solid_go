package main

import "fmt"

// LeetCode 41 - First Missing Positive
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
