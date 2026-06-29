package leetcode

// LeetCode 200 - Number of Islands
//
// Pseudo code:
//   count = 0
//   for each cell (i,j):
//     if grid[i][j] == '1':
//       count++
//       BFS/DFS to mark all connected '1's as visited ('0')
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
