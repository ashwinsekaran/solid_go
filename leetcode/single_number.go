package leetcode

// LeetCode 136 - Single Number
//
// Pseudo code:
//   result = 0
//   XOR all numbers: duplicates cancel out (a ^ a = 0), leaving the single number
//   return result

func singleNumber(nums []int) int {
	result := 0
	for _, n := range nums {
		result ^= n
	}
	return result
}
