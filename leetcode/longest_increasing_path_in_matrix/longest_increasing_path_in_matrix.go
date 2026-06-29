package main

import "fmt"

// LeetCode 329 - Longest Increasing Path in a Matrix
//
// Pseudo code:
//   DFS with memoization from every cell
//   from cell (i,j): try all 4 neighbors with strictly greater value
//   memo[i][j] = 1 + max(dfs(neighbor))
//   return max over all cells

func longestIncreasingPath(matrix [][]int) int {
	m, n := len(matrix), len(matrix[0])
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n)
	}
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		if memo[i][j] != 0 {
			return memo[i][j]
		}
		best := 1
		for _, d := range dirs {
			ni, nj := i+d[0], j+d[1]
			if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] > matrix[i][j] {
				if v := 1 + dfs(ni, nj); v > best {
					best = v
				}
			}
		}
		memo[i][j] = best
		return best
	}
	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if v := dfs(i, j); v > res {
				res = v
			}
		}
	}
	return res
}

func main() {
	fmt.Println(longestIncreasingPath([][]int{{9, 9, 4}, {6, 6, 8}, {2, 1, 1}}))   // 4
	fmt.Println(longestIncreasingPath([][]int{{3, 4, 5}, {3, 2, 6}, {2, 2, 1}}))   // 4
	fmt.Println(longestIncreasingPath([][]int{{1}}))                                 // 1
}
