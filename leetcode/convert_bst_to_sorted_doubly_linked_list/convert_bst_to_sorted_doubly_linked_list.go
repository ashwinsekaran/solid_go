package main

import "fmt"

// LeetCode 426 - Convert Binary Search Tree to Sorted Doubly Linked List
//
// Pseudo code:
//   in-order traversal (left, root, right)
//   link each visited node to prev; track head (first visited)
//   after traversal: circularize head <-> tail

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Prev  *Node
	Next  *Node
}

func treeToDoublyList(root *Node) *Node {
	if root == nil {
		return nil
	}
	var head, prev *Node
	var inorder func(node *Node)
	inorder = func(node *Node) {
		if node == nil {
			return
		}
		inorder(node.Left)
		if prev == nil {
			head = node
		} else {
			prev.Next = node
			node.Prev = prev
		}
		prev = node
		inorder(node.Right)
	}
	inorder(root)
	head.Prev = prev
	prev.Next = head
	return head
}

func main() {
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

	// Print circular list forward (5 elements)
	fmt.Print("Forward: ")
	curr := head
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", curr.Val)
		curr = curr.Next
	}
	fmt.Println()

	// Print backward from tail
	fmt.Print("Backward: ")
	tail := head.Prev
	curr = tail
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", curr.Val)
		curr = curr.Prev
	}
	fmt.Println()
}
