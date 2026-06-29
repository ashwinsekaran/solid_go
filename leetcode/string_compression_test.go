package leetcode

import "testing"

func TestCompress(t *testing.T) {
	tests := []struct {
		chars []byte
		want  int
		out   string
	}{
		{[]byte{'a', 'a', 'b', 'b', 'c', 'c', 'c'}, 6, "a2b2c3"},
		{[]byte{'a'}, 1, "a"},
		{[]byte{'a', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b'}, 4, "ab12"},
	}
	for _, tc := range tests {
		n := compress(tc.chars)
		if n != tc.want {
			t.Errorf("compress length = %d, want %d", n, tc.want)
		}
		if string(tc.chars[:n]) != tc.out {
			t.Errorf("compress result = %q, want %q", string(tc.chars[:n]), tc.out)
		}
	}
}
