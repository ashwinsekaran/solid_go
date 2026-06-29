package leetcode

import (
	"reflect"
	"testing"
)

func listToSlice(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

func TestReverseKGroup(t *testing.T) {
	tests := []struct {
		vals []int
		k    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{2, 1, 4, 3, 5}},
		{[]int{1, 2, 3, 4, 5}, 3, []int{3, 2, 1, 4, 5}},
		{[]int{1, 2, 3, 4}, 4, []int{4, 3, 2, 1}},
	}
	for _, tc := range tests {
		got := listToSlice(reverseKGroup(makeList(tc.vals), tc.k))
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("reverseKGroup(%v, k=%d) = %v, want %v", tc.vals, tc.k, got, tc.want)
		}
	}
}
