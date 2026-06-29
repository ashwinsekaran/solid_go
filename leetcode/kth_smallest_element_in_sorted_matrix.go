package leetcode

// LeetCode 378 - Kth Smallest Element in a Sorted Matrix
//
// Pseudo code:
//   binary search on value range [matrix[0][0], matrix[n-1][n-1]]
//   for a mid value, count elements <= mid:
//     walk from bottom-left: if matrix[r][c] <= mid, count += r+1, move right
//     else move up
//   find smallest value where count >= k

func kthSmallest(matrix [][]int, k int) int {
	n := len(matrix)
	lo, hi := matrix[0][0], matrix[n-1][n-1]
	for lo < hi {
		mid := lo + (hi-lo)/2
		count := 0
		r, c := n-1, 0
		for r >= 0 && c < n {
			if matrix[r][c] <= mid {
				count += r + 1
				c++
			} else {
				r--
			}
		}
		if count >= k {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}
