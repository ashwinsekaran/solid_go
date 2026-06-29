package leetcode

// LeetCode 498 - Diagonal Traverse
//
// Pseudo code:
//   result = []
//   for d from 0 to m+n-2:
//     if d is even: traverse diagonal upward (row decreases, col increases)
//     else: traverse diagonal downward (row increases, col decreases)
//     clamp row and col to valid bounds
//   return result

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
