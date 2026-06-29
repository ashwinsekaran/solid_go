package leetcode

import "testing"

func TestJobScheduling(t *testing.T) {
	tests := []struct {
		start, end, profit []int
		want               int
	}{
		{[]int{1, 2, 3, 3}, []int{3, 4, 5, 6}, []int{50, 10, 40, 70}, 120},
		{[]int{1, 2, 3, 4, 6}, []int{3, 5, 10, 6, 9}, []int{20, 20, 100, 70, 60}, 150},
		{[]int{1, 1, 1}, []int{2, 3, 4}, []int{5, 6, 4}, 6},
	}
	for _, tc := range tests {
		got := jobScheduling(tc.start, tc.end, tc.profit)
		if got != tc.want {
			t.Errorf("jobScheduling = %d, want %d", got, tc.want)
		}
	}
}
