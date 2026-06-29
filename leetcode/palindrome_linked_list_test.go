package leetcode

import "testing"

func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func TestIsPalindromeList(t *testing.T) {
	tests := []struct {
		vals []int
		want bool
	}{
		{[]int{1, 2, 2, 1}, true},
		{[]int{1, 2}, false},
		{[]int{1}, true},
		{[]int{1, 0, 0, 1}, true},
	}
	for _, tc := range tests {
		got := isPalindromeList(makeList(tc.vals))
		if got != tc.want {
			t.Errorf("isPalindromeList(%v) = %v, want %v", tc.vals, got, tc.want)
		}
	}
}
