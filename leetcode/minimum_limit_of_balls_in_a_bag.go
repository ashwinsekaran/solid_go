package leetcode

// LeetCode 1760 - Minimum Limit of Balls in a Bag
//
// Pseudo code:
//   binary search on answer (penalty = max balls in any bag)
//   for a given max size, count operations needed:
//     for each bag of n balls: ops += ceil(n/maxSize) - 1
//   find smallest maxSize where ops <= maxOperations

func minimumSize(nums []int, maxOperations int) int {
	lo, hi := 1, 1000000000
	for lo < hi {
		mid := lo + (hi-lo)/2
		ops := 0
		for _, n := range nums {
			ops += (n - 1) / mid
		}
		if ops <= maxOperations {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}
