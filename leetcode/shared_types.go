package leetcode

// Shared data structures used across multiple solutions

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type ListNode struct {
	Val  int
	Next *ListNode
}

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Prev  *Node
	Next  *Node
	Child *Node
}
