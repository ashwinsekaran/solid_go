package main

import "fmt"

// LeetCode 136 - Single Number
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
