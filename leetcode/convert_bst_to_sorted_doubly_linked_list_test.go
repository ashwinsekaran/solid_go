package leetcode

import (
	"testing"
)

func TestTreeToDoublyList(t *testing.T) {
	//   4
	//  / \
	// 2   5
	// /\
	//1  3
	root := &Node{Val: 4,
		Left: &Node{Val: 2,
			Left:  &Node{Val: 1},
			Right: &Node{Val: 3},
		},
		Right: &Node{Val: 5},
	}
	head := treeToDoublyList(root)

	// Collect values walking forward until we loop back to head
	vals := []int{}
	curr := head
	for {
		vals = append(vals, curr.Val)
		curr = curr.Next
		if curr == head {
			break
		}
	}
	want := []int{1, 2, 3, 4, 5}
	for i, v := range vals {
		if v != want[i] {
			t.Errorf("doubly list[%d] = %d, want %d", i, v, want[i])
		}
	}
	// Verify circular: tail.Next == head and head.Prev == tail
	tail := head.Prev
	if tail.Next != head {
		t.Error("list is not circular: tail.Next != head")
	}
}
