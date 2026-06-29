package main

import "fmt"

// LeetCode 64 - Minimum Path Sum
//
// Problem:
//   Given an m x n grid of non-negative integers, find a path from the
//   top-left to the bottom-right that minimises the sum of numbers along the
//   path. You can only move right or down.
//
// Example:
//   Input:  grid = [[1,3,1],[1,5,1],[4,2,1]]
//   Output: 7
//
//   Explanation:
//   Path: 1 → 3 → 1 → 1 → 1  =  7
//   (right, right, down, down)
//
// Pseudo code:
//   for each cell (i,j):
//     if top-left corner: skip
//     if top row: add from left
//     if left col: add from above
//     else: add min(above, left)
//   return dp[m-1][n-1]

func minPathSum(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if i == 0 && j == 0 {
				continue
			} else if i == 0 {
				grid[i][j] += grid[i][j-1]
			} else if j == 0 {
				grid[i][j] += grid[i-1][j]
			} else {
				a, b := grid[i-1][j], grid[i][j-1]
				if b < a {
					a = b
				}
				grid[i][j] += a
			}
		}
	}
	return grid[m-1][n-1]
}

func main() {
	fmt.Println(minPathSum([][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}})) // 7
	fmt.Println(minPathSum([][]int{{1, 2, 3}, {4, 5, 6}}))             // 12
}
