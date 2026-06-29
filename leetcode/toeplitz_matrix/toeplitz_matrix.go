package main

import "fmt"

// LeetCode 766 - Toeplitz Matrix
//
// Pseudo code:
//   for every cell (i,j) not in the first row or first col:
//     it must equal matrix[i-1][j-1] (its top-left diagonal neighbour)
//   return true if all match

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

func main() {
	fmt.Println(isToeplitzMatrix([][]int{{1, 2, 3, 4}, {5, 1, 2, 3}, {9, 5, 1, 2}})) // true
	fmt.Println(isToeplitzMatrix([][]int{{1, 2}, {2, 2}}))                             // false
}
