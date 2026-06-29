package leetcode

import (
	"reflect"
	"testing"
)

func TestReverseList(t *testing.T) {
	tests := []struct {
		vals []int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{1}, []int{1}},
	}
	for _, tc := range tests {
		got := listToSlice(reverseList(makeList(tc.vals)))
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("reverseList(%v) = %v, want %v", tc.vals, got, tc.want)
		}
	}
}
