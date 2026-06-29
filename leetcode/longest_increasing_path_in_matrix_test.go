package leetcode

import "testing"

func TestLongestIncreasingPath(t *testing.T) {
	tests := []struct {
		matrix [][]int
		want   int
	}{
		{[][]int{{9, 9, 4}, {6, 6, 8}, {2, 1, 1}}, 4},
		{[][]int{{3, 4, 5}, {3, 2, 6}, {2, 2, 1}}, 4},
		{[][]int{{1}}, 1},
	}
	for _, tc := range tests {
		got := longestIncreasingPath(tc.matrix)
		if got != tc.want {
			t.Errorf("longestIncreasingPath(%v) = %d, want %d", tc.matrix, got, tc.want)
		}
	}
}
