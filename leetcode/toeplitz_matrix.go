package leetcode

// LeetCode 766 - Toeplitz Matrix
//
// Pseudo code:
//   for each cell (i,j) not in first row/col:
//     if matrix[i][j] != matrix[i-1][j-1]: return false
//   return true

func isToeplitzMatrix(matrix [][]int) bool {
	m, n := len(matrix), len(matrix[0])
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if matrix[i][j] != matrix[i-1][j-1] {
				return false
			}
		}
	}
	return true
}
