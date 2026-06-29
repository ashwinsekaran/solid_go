package leetcode

import "testing"

func TestKthSmallest(t *testing.T) {
	tests := []struct {
		matrix [][]int
		k      int
		want   int
	}{
		{[][]int{{1, 5, 9}, {10, 11, 13}, {12, 13, 15}}, 8, 13},
		{[][]int{{-5}}, 1, -5},
		{[][]int{{1, 2}, {1, 3}}, 2, 1},
	}
	for _, tc := range tests {
		got := kthSmallest(tc.matrix, tc.k)
		if got != tc.want {
			t.Errorf("kthSmallest(k=%d) = %d, want %d", tc.k, got, tc.want)
		}
	}
}
