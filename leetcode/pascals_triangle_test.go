package leetcode

import (
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		numRows int
		want    [][]int
	}{
		{1, [][]int{{1}}},
		{3, [][]int{{1}, {1, 1}, {1, 2, 1}}},
		{5, [][]int{{1}, {1, 1}, {1, 2, 1}, {1, 3, 3, 1}, {1, 4, 6, 4, 1}}},
	}
	for _, tc := range tests {
		got := generate(tc.numRows)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("generate(%d) = %v, want %v", tc.numRows, got, tc.want)
		}
	}
}
