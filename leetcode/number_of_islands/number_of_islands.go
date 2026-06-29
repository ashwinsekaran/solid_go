package main

import "fmt"

// LeetCode 200 - Number of Islands
//
// Problem:
//   Given an m x n 2D grid of '1's (land) and '0's (water), return the number
//   of islands. An island is surrounded by water and formed by connecting
//   adjacent land cells horizontally or vertically.
//
// Example 1:
//   Input:
//   [["1","1","1","1","0"],
//    ["1","1","0","1","0"],
//    ["1","1","0","0","0"],
//    ["0","0","0","0","0"]]
//   Output: 1  (all land is connected)
//
// Example 2:
//   Input:
//   [["1","1","0","0","0"],
//    ["1","1","0","0","0"],
//    ["0","0","1","0","0"],
//    ["0","0","0","1","1"]]
//   Output: 3  (three separate islands)
//
// Pseudo code:
//   for each cell that is '1': increment count, DFS to mark all connected '1's as '0'
//   return count

func numIslands(grid [][]byte) int {
	m, n := len(grid), len(grid[0])
	count := 0
	var dfs func(i, j int)
	dfs = func(i, j int) {
		if i < 0 || i >= m || j < 0 || j >= n || grid[i][j] != '1' {
			return
		}
		grid[i][j] = '0'
		dfs(i+1, j)
		dfs(i-1, j)
		dfs(i, j+1)
		dfs(i, j-1)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				count++
				dfs(i, j)
			}
		}
	}
	return count
}

func main() {
	g1 := [][]byte{
		{'1', '1', '1', '1', '0'},
		{'1', '1', '0', '1', '0'},
		{'1', '1', '0', '0', '0'},
		{'0', '0', '0', '0', '0'},
	}
	fmt.Println(numIslands(g1)) // 1

	g2 := [][]byte{
		{'1', '1', '0', '0', '0'},
		{'1', '1', '0', '0', '0'},
		{'0', '0', '1', '0', '0'},
		{'0', '0', '0', '1', '1'},
	}
	fmt.Println(numIslands(g2)) // 3
}
