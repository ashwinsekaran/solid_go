package leetcode

import (
	"reflect"
	"testing"
)

func TestZigzagLevelOrder(t *testing.T) {
	//     3
	//    / \
	//   9  20
	//     /  \
	//    15   7
	root := &TreeNode{Val: 3,
		Left: &TreeNode{Val: 9},
		Right: &TreeNode{Val: 20,
			Left:  &TreeNode{Val: 15},
			Right: &TreeNode{Val: 7},
		},
	}
	want := [][]int{{3}, {20, 9}, {15, 7}}
	got := zigzagLevelOrder(root)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("zigzagLevelOrder = %v, want %v", got, want)
	}

	if zigzagLevelOrder(nil) != nil {
		t.Error("expected nil for nil root")
	}
}
