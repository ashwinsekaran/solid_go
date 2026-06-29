package leetcode

import "testing"

func TestNumBusesToDestination(t *testing.T) {
	tests := []struct {
		routes        [][]int
		source, target int
		want          int
	}{
		{[][]int{{1, 2, 7}, {3, 6, 7}}, 1, 6, 2},
		{[][]int{{7, 12}, {4, 5, 15}, {6}, {15, 19}, {9, 12, 13}}, 15, 12, -1},
		{[][]int{{1, 2, 7}}, 1, 1, 0},
	}
	for _, tc := range tests {
		got := numBusesToDestination(tc.routes, tc.source, tc.target)
		if got != tc.want {
			t.Errorf("numBusesToDestination routes=%v src=%d tgt=%d = %d, want %d",
				tc.routes, tc.source, tc.target, got, tc.want)
		}
	}
}
