package main

import "fmt"

// LeetCode 378 - Kth Smallest Element in a Sorted Matrix
//
// Pseudo code:
//   binary search on the value range [matrix[0][0], matrix[n-1][n-1]]
//   for a mid value, count elements <= mid using a staircase walk from bottom-left
//   find smallest value where count >= k

func kthSmallest(matrix [][]int, k int) int {
	n := len(matrix)
	lo, hi := matrix[0][0], matrix[n-1][n-1]
	for lo < hi {
		mid := lo + (hi-lo)/2
		count, r, c := 0, n-1, 0
		for r >= 0 && c < n {
			if matrix[r][c] <= mid {
				count += r + 1
				c++
			} else {
				r--
			}
		}
		if count >= k {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func main() {
	fmt.Println(kthSmallest([][]int{{1, 5, 9}, {10, 11, 13}, {12, 13, 15}}, 8)) // 13
	fmt.Println(kthSmallest([][]int{{-5}}, 1))                                   // -5
	fmt.Println(kthSmallest([][]int{{1, 2}, {1, 3}}, 2))                         // 1
}
