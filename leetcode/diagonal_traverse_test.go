package leetcode

import (
	"reflect"
	"testing"
)

func TestFindDiagonalOrder(t *testing.T) {
	tests := []struct {
		mat  [][]int
		want []int
	}{
		{
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[]int{1, 2, 4, 7, 5, 3, 6, 8, 9},
		},
		{
			[][]int{{1, 2}, {3, 4}},
			[]int{1, 2, 3, 4},
		},
		{
			[][]int{{1}},
			[]int{1},
		},
	}
	for _, tc := range tests {
		got := findDiagonalOrder(tc.mat)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("findDiagonalOrder(%v) = %v, want %v", tc.mat, got, tc.want)
		}
	}
}
