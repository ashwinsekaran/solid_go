package leetcode

import "testing"

func TestIsIsomorphic(t *testing.T) {
	tests := []struct {
		s, t string
		want bool
	}{
		{"egg", "add", true},
		{"foo", "bar", false},
		{"paper", "title", true},
		{"badc", "baba", false},
	}
	for _, tc := range tests {
		got := isIsomorphic(tc.s, tc.t)
		if got != tc.want {
			t.Errorf("isIsomorphic(%q, %q) = %v, want %v", tc.s, tc.t, got, tc.want)
		}
	}
}
