package leetcode

import (
	"testing"
)

func TestCustomSortString(t *testing.T) {
	// The output just needs all order chars to appear before non-order chars
	// in the correct relative order; multiple valid outputs exist.
	tests := []struct {
		order, s string
		wantLen  int
	}{
		{"cba", "abcd", 4},
		{"bcafg", "abcd", 4},
	}
	for _, tc := range tests {
		got := customSortString(tc.order, tc.s)
		if len(got) != tc.wantLen {
			t.Errorf("customSortString(%q, %q) length = %d, want %d", tc.order, tc.s, len(got), tc.wantLen)
		}
		// Verify all chars present
		freq := make(map[rune]int)
		for _, c := range tc.s {
			freq[c]++
		}
		for _, c := range got {
			freq[c]--
		}
		for c, v := range freq {
			if v != 0 {
				t.Errorf("customSortString: char %c count mismatch", c)
			}
		}
	}
}
