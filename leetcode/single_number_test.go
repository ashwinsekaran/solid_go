package leetcode

import "testing"

func TestSingleNumber(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{2, 2, 1}, 1},
		{[]int{4, 1, 2, 1, 2}, 4},
		{[]int{1}, 1},
	}
	for _, tc := range tests {
		got := singleNumber(tc.nums)
		if got != tc.want {
			t.Errorf("singleNumber(%v) = %d, want %d", tc.nums, got, tc.want)
		}
	}
}
