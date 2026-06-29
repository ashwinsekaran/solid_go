package main

import "fmt"

// LeetCode 1760 - Minimum Limit of Balls in a Bag
//
// Pseudo code:
//   binary search on the answer (max balls allowed in any bag)
//   for a given max size: ops needed = sum of ceil(n/max)-1 for each bag
//   find smallest max where ops <= maxOperations

func minimumSize(nums []int, maxOperations int) int {
	lo, hi := 1, 1_000_000_000
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

func main() {
	fmt.Println(minimumSize([]int{9}, 2))              // 3
	fmt.Println(minimumSize([]int{2, 4, 8, 2}, 4))    // 2
	fmt.Println(minimumSize([]int{7, 17}, 2))          // 7
}
