package leetcode

// LeetCode 64 - Minimum Path Sum
//
// Pseudo code:
//   m, n = grid dimensions
//   dp = copy of grid
//   for i in rows:
//     for j in cols:
//       if i==0 && j==0: skip
//       else if i==0: dp[i][j] += dp[i][j-1]
//       else if j==0: dp[i][j] += dp[i-1][j]
//       else: dp[i][j] += min(dp[i-1][j], dp[i][j-1])
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
				grid[i][j] += min2(grid[i-1][j], grid[i][j-1])
			}
		}
	}
	return grid[m-1][n-1]
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}
