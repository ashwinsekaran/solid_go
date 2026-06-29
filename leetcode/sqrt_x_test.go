package leetcode

import "testing"

func TestMySqrt(t *testing.T) {
	tests := []struct {
		x, want int
	}{
		{0, 0},
		{1, 1},
		{4, 2},
		{8, 2},
		{9, 3},
		{2147395600, 46340},
	}
	for _, tc := range tests {
		got := mySqrt(tc.x)
		if got != tc.want {
			t.Errorf("mySqrt(%d) = %d, want %d", tc.x, got, tc.want)
		}
	}
}
