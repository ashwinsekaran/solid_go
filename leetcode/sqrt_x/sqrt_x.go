package main

import "fmt"

// LeetCode 69 - Sqrt(x)
//
// Problem:
//   Given a non-negative integer x, return the integer square root of x
//   (i.e., floor(sqrt(x))). Do not use built-in exponent or sqrt functions.
//
// Example 1:
//   Input:  x = 4    Output: 2
//
// Example 2:
//   Input:  x = 8    Output: 2
//   Explanation: sqrt(8) ≈ 2.828..., floor gives 2.
//
// Example 3:
//   Input:  x = 9    Output: 3
//
// Pseudo code:
//   binary search in [0, x]
//   find largest m where m*m <= x

func mySqrt(x int) int {
	lo, hi := 0, x
	for lo < hi {
		mid := lo + (hi-lo+1)/2
		if mid*mid <= x {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func main() {
	fmt.Println(mySqrt(0))  // 0
	fmt.Println(mySqrt(1))  // 1
	fmt.Println(mySqrt(4))  // 2
	fmt.Println(mySqrt(8))  // 2
	fmt.Println(mySqrt(9))  // 3
	fmt.Println(mySqrt(16)) // 4
}
