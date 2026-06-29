package leetcode

// LeetCode 329 - Longest Increasing Path in a Matrix
//
// Pseudo code:
//   memo = cache for each cell's longest increasing path
//   for each cell (i,j): dfs to find longest path starting there
//   dfs(i,j): if cached return; try 4 neighbors with strictly greater value
//             result = 1 + max(dfs(neighbor))
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
				v := 1 + dfs(ni, nj)
				if v > best {
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
			v := dfs(i, j)
			if v > res {
				res = v
			}
		}
	}
	return res
}
