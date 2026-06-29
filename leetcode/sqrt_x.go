package leetcode

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
