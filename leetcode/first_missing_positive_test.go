package leetcode

import "testing"

func TestFirstMissingPositive(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{1, 2, 0}, 3},
		{[]int{3, 4, -1, 1}, 2},
		{[]int{7, 8, 9, 11, 12}, 1},
		{[]int{1}, 2},
		{[]int{1, 2, 3}, 4},
	}
	for _, tc := range tests {
		got := firstMissingPositive(tc.nums)
		if got != tc.want {
			t.Errorf("firstMissingPositive(%v) = %d, want %d", tc.nums, got, tc.want)
		}
	}
}
