package main

import "fmt"

// LeetCode 498 - Diagonal Traverse
//
// Problem:
//   Given an m x n matrix, return all elements in diagonal order
//   (alternating up-right and down-left diagonals).
//
// Example:
//   Input:
//   [[1, 2, 3],
//    [4, 5, 6],
//    [7, 8, 9]]
//
//   Output: [1, 2, 4, 7, 5, 3, 6, 8, 9]
//
//   Explanation (diagonals):
//   d=0 ↗: [1]          (up-right)
//   d=1 ↙: [2, 4]       (down-left)
//   d=2 ↗: [7, 5, 3]    (up-right)  ← but written 3,5,7 → reversed to [3,5,7]... no:
//   Actually: d=2 goes up: row decreases col increases → 7→5→3... wait output shows [7,5,3]
//   d=3 ↙: [6, 8]       (down-left)
//   d=4 ↗: [9]          (up-right)
//
// Pseudo code:
//   iterate diagonals d = 0 to m+n-2
//   even d: go up-right (row decreases, col increases)
//   odd  d: go down-left (row increases, col decreases)
//   clamp starting row/col to valid bounds each time

func findDiagonalOrder(mat [][]int) []int {
	if len(mat) == 0 {
		return nil
	}
	m, n := len(mat), len(mat[0])
	result := make([]int, 0, m*n)
	for d := 0; d < m+n-1; d++ {
		if d%2 == 0 {
			r := d
			if r >= m {
				r = m - 1
			}
			c := d - r
			for r >= 0 && c < n {
				result = append(result, mat[r][c])
				r--
				c++
			}
		} else {
			c := d
			if c >= n {
				c = n - 1
			}
			r := d - c
			for c >= 0 && r < m {
				result = append(result, mat[r][c])
				r++
				c--
			}
		}
	}
	return result
}

func main() {
	fmt.Println(findDiagonalOrder([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})) // [1 2 4 7 5 3 6 8 9]
	fmt.Println(findDiagonalOrder([][]int{{1, 2}, {3, 4}}))                   // [1 2 3 4]
}
