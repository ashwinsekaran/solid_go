package leetcode

import "testing"

func TestMinimumSize(t *testing.T) {
	tests := []struct {
		nums          []int
		maxOperations int
		want          int
	}{
		{[]int{9}, 2, 3},
		{[]int{2, 4, 8, 2}, 4, 2},
		{[]int{7, 17}, 2, 7},
	}
	for _, tc := range tests {
		got := minimumSize(tc.nums, tc.maxOperations)
		if got != tc.want {
			t.Errorf("minimumSize(%v, %d) = %d, want %d", tc.nums, tc.maxOperations, got, tc.want)
		}
	}
}
