package leetcode

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		x    int
		want bool
	}{
		{121, true},
		{-121, false},
		{10, false},
		{0, true},
		{1221, true},
	}
	for _, tc := range tests {
		got := isPalindrome(tc.x)
		if got != tc.want {
			t.Errorf("isPalindrome(%d) = %v, want %v", tc.x, got, tc.want)
		}
	}
}
