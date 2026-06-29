package leetcode

import "testing"

func TestSubarraysWithKDistinct(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want int
	}{
		{[]int{1, 2, 1, 2, 3}, 2, 7},
		{[]int{1, 2, 1, 3, 4}, 3, 3},
		{[]int{1, 1, 1, 1}, 1, 10},
	}
	for _, tc := range tests {
		got := subarraysWithKDistinct(tc.nums, tc.k)
		if got != tc.want {
			t.Errorf("subarraysWithKDistinct(%v, %d) = %d, want %d", tc.nums, tc.k, got, tc.want)
		}
	}
}
