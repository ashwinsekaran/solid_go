package leetcode

import "testing"

func TestMyAtoi(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"42", 42},
		{"   -42", -42},
		{"4193 with words", 4193},
		{"words and 987", 0},
		{"-91283472332", -2147483648},
		{"2147483648", 2147483647},
		{"", 0},
		{"+1", 1},
	}
	for _, tc := range tests {
		got := myAtoi(tc.s)
		if got != tc.want {
			t.Errorf("myAtoi(%q) = %d, want %d", tc.s, got, tc.want)
		}
	}
}
