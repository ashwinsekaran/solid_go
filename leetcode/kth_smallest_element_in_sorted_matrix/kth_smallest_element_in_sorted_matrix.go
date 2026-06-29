package main

import "fmt"

// LeetCode 378 - Kth Smallest Element in a Sorted Matrix
//
// Problem:
//   Given an n×n matrix where each row and column is sorted in ascending order,
//   return the k-th smallest element.
//
// Example:
//   Input:
//     matrix = [[1, 5, 9],
//               [10,11,13],
//               [12,13,15]]
//     k = 8
//   Output: 13
//
//   Explanation:
//   Sorted elements: [1,5,9,10,11,12,13,13,15]
//   The 8th smallest is 13.
//   We binary search on the value range and count how many elements ≤ mid
//   using a staircase walk from the bottom-left corner.
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
