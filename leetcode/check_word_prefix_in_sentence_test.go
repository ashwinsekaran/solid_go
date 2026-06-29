package leetcode

import "testing"

func TestIsPrefixOfWord(t *testing.T) {
	tests := []struct {
		sentence, searchWord string
		want                 int
	}{
		{"i love eating burger", "burg", 4},
		{"this problem is an easy problem", "pro", 2},
		{"i am tired", "you", -1},
		{"hello from the other side", "they", -1},
	}
	for _, tc := range tests {
		got := isPrefixOfWord(tc.sentence, tc.searchWord)
		if got != tc.want {
			t.Errorf("isPrefixOfWord(%q, %q) = %d, want %d", tc.sentence, tc.searchWord, got, tc.want)
		}
	}
}
