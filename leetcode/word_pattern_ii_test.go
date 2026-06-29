package leetcode

import "testing"

func TestWordPatternMatch(t *testing.T) {
	tests := []struct {
		pattern, s string
		want       bool
	}{
		{"abab", "redblueredblue", true},
		{"aaaa", "asdasdasdasd", true},
		{"aabb", "xyzxyz", false},
		{"ab", "aa", false},
	}
	for _, tc := range tests {
		got := wordPatternMatch(tc.pattern, tc.s)
		if got != tc.want {
			t.Errorf("wordPatternMatch(%q, %q) = %v, want %v", tc.pattern, tc.s, got, tc.want)
		}
	}
}
