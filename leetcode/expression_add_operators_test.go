package leetcode

import (
	"sort"
	"reflect"
	"testing"
)

func TestAddOperators(t *testing.T) {
	tests := []struct {
		num    string
		target int
		want   []string
	}{
		{"123", 6, []string{"1*2*3", "1+2+3"}},
		{"232", 8, []string{"2*3+2", "2+3*2"}},
		{"3456237490", 9191, []string{}},
		{"105", 5, []string{"1*0+5", "10-5"}},
	}
	for _, tc := range tests {
		got := addOperators(tc.num, tc.target)
		sort.Strings(got)
		sort.Strings(tc.want)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("addOperators(%q, %d) = %v, want %v", tc.num, tc.target, got, tc.want)
		}
	}
}
