package main

import "fmt"

// LeetCode 136 - Single Number
//
// Problem:
//   Given a non-empty array of integers where every element appears twice
//   except for one, find that single one. Must run in O(n) time and O(1) space.
//
// Example 1:
//   Input:  nums = [2, 2, 1]
//   Output: 1
//
// Example 2:
//   Input:  nums = [4, 1, 2, 1, 2]
//   Output: 4
//
//   Explanation:
//   XOR is self-inverse: a ^ a = 0 and a ^ 0 = a.
//   So XOR-ing all values cancels every duplicate, leaving only the unique one.
//   4 ^ 1 ^ 2 ^ 1 ^ 2  =  4 ^ (1^1) ^ (2^2)  =  4 ^ 0 ^ 0  =  4
//
// Pseudo code:
//   XOR all numbers together
//   duplicate pairs cancel out (a ^ a = 0), leaving only the unique number

func singleNumber(nums []int) int {
	result := 0
	for _, n := range nums {
		result ^= n
	}
	return result
}

func main() {
	fmt.Println(singleNumber([]int{2, 2, 1}))         // 1
	fmt.Println(singleNumber([]int{4, 1, 2, 1, 2}))   // 4
	fmt.Println(singleNumber([]int{1}))                // 1
}
