package leetcode

import "testing"

func TestMaximumUnits(t *testing.T) {
	tests := []struct {
		boxTypes  [][]int
		truckSize int
		want      int
	}{
		{[][]int{{1, 3}, {2, 2}, {3, 1}}, 4, 8},
		{[][]int{{5, 10}, {2, 5}, {4, 7}, {3, 9}}, 10, 91},
	}
	for _, tc := range tests {
		got := maximumUnits(tc.boxTypes, tc.truckSize)
		if got != tc.want {
			t.Errorf("maximumUnits(%v, %d) = %d, want %d", tc.boxTypes, tc.truckSize, got, tc.want)
		}
	}
}
