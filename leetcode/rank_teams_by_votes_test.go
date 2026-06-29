package leetcode

import "testing"

func TestRankTeams(t *testing.T) {
	tests := []struct {
		votes []string
		want  string
	}{
		{[]string{"ABC", "ACB", "ABC", "ACB", "ACB"}, "ACB"},
		{[]string{"WXYZ", "XWYZ", "ZWYX"}, "WXZY"},
		{[]string{"BCA", "CAB", "CBA", "ABC", "ACB", "BAC"}, "ABC"},
	}
	for _, tc := range tests {
		got := rankTeams(tc.votes)
		if got != tc.want {
			t.Errorf("rankTeams(%v) = %q, want %q", tc.votes, got, tc.want)
		}
	}
}
