package main

import "fmt"

// LeetCode 69 - Sqrt(x)
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
