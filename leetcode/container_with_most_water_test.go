package leetcode

import "testing"

func TestMaxArea(t *testing.T) {
	tests := []struct {
		height []int
		want   int
	}{
		{[]int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{[]int{1, 1}, 1},
		{[]int{4, 3, 2, 1, 4}, 16},
		{[]int{1, 2, 1}, 2},
	}
	for _, tc := range tests {
		got := maxArea(tc.height)
		if got != tc.want {
			t.Errorf("maxArea(%v) = %d, want %d", tc.height, got, tc.want)
		}
	}
}
