package leetcode

import (
	"reflect"
	"testing"
)

func TestRightSideView(t *testing.T) {
	//     1
	//    / \
	//   2   3
	//    \
	//     5
	root := &TreeNode{Val: 1,
		Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3},
	}
	got := rightSideView(root)
	want := []int{1, 3, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("rightSideView = %v, want %v", got, want)
	}

	if rightSideView(nil) != nil {
		t.Error("expected nil for nil root")
	}
}
