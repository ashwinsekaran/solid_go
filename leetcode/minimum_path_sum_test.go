package leetcode

import "testing"

func TestMinPathSum(t *testing.T) {
	tests := []struct {
		grid [][]int
		want int
	}{
		{[][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}, 7},
		{[][]int{{1, 2, 3}, {4, 5, 6}}, 12},
		{[][]int{{1}}, 1},
	}
	for _, tc := range tests {
		got := minPathSum(tc.grid)
		if got != tc.want {
			t.Errorf("minPathSum(%v) = %d, want %d", tc.grid, got, tc.want)
		}
	}
}
