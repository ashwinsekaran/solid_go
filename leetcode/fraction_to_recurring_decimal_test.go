package leetcode

import "testing"

func TestFractionToDecimal(t *testing.T) {
	tests := []struct {
		num, den int
		want     string
	}{
		{1, 2, "0.5"},
		{2, 1, "2"},
		{4, 333, "0.(012)"},
		{1, 6, "0.1(6)"},
		{0, 3, "0"},
		{-1, 2, "-0.5"},
		{-50, 8, "-6.25"},
	}
	for _, tc := range tests {
		got := fractionToDecimal(tc.num, tc.den)
		if got != tc.want {
			t.Errorf("fractionToDecimal(%d, %d) = %q, want %q", tc.num, tc.den, got, tc.want)
		}
	}
}
