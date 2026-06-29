package leetcode

import "testing"

func TestIsToeplitzMatrix(t *testing.T) {
	tests := []struct {
		matrix [][]int
		want   bool
	}{
		{[][]int{{1, 2, 3, 4}, {5, 1, 2, 3}, {9, 5, 1, 2}}, true},
		{[][]int{{1, 2}, {2, 2}}, false},
		{[][]int{{1}}, true},
	}
	for _, tc := range tests {
		got := isToeplitzMatrix(tc.matrix)
		if got != tc.want {
			t.Errorf("isToeplitzMatrix(%v) = %v, want %v", tc.matrix, got, tc.want)
		}
	}
}
