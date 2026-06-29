package main

import "fmt"

// LeetCode 33 - Search in Rotated Sorted Array
//
// Problem:
//   An ascending array was rotated at some pivot. Given the rotated array and a
//   target, return the index of target or -1 if not found. Must be O(log n).
//
// Example 1:
//   Input:  nums = [4,5,6,7,0,1,2], target = 0
//   Output: 4
//   Explanation: The array was rotated at index 4. Target 0 is at index 4.
//
// Example 2:
//   Input:  nums = [4,5,6,7,0,1,2], target = 3
//   Output: -1
//
//   Key insight: at any mid point, one of the two halves is always sorted.
//   Check which half is sorted, then decide which half the target falls in.
//
// Pseudo code:
//   binary search; one half is always sorted
//   if left half sorted and target in [nums[lo], nums[mid]): go left
//   else go right
//   mirror logic when right half is sorted

func search(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if nums[mid] == target {
			return mid
		}
		if nums[lo] <= nums[mid] { // left half sorted
			if nums[lo] <= target && target < nums[mid] {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		} else { // right half sorted
			if nums[mid] < target && target <= nums[hi] {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
	}
	return -1
}

func main() {
	fmt.Println(search([]int{4, 5, 6, 7, 0, 1, 2}, 0))  // 4
	fmt.Println(search([]int{4, 5, 6, 7, 0, 1, 2}, 3))  // -1
	fmt.Println(search([]int{1}, 0))                      // -1
	fmt.Println(search([]int{1, 3}, 3))                   // 1
}
