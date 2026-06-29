package leetcode

import "testing"

func TestIsZeroArray(t *testing.T) {
	tests := []struct {
		nums    []int
		queries [][]int
		want    bool
	}{
		{[]int{1, 0, 1}, [][]int{{0, 2}}, true},
		{[]int{4, 3, 2, 1}, [][]int{{1, 3}, {0, 2}}, false},
		{[]int{0, 0, 0}, [][]int{}, true},
	}
	for _, tc := range tests {
		got := isZeroArray(tc.nums, tc.queries)
		if got != tc.want {
			t.Errorf("isZeroArray(%v, %v) = %v, want %v", tc.nums, tc.queries, got, tc.want)
		}
	}
}
