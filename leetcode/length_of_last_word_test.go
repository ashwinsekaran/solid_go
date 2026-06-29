package leetcode

import "testing"

func TestLengthOfLastWord(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"Hello World", 5},
		{"   fly me   to   the moon  ", 4},
		{"luffy is still joyboy", 6},
		{"a", 1},
	}
	for _, tc := range tests {
		got := lengthOfLastWord(tc.s)
		if got != tc.want {
			t.Errorf("lengthOfLastWord(%q) = %d, want %d", tc.s, got, tc.want)
		}
	}
}
