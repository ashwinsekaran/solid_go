package leetcode

// LeetCode 426 - Convert Binary Search Tree to Sorted Doubly Linked List
//
// Pseudo code:
//   in-order traversal (left, root, right)
//   maintain prev pointer
//   link prev.right = curr, curr.left = prev
//   track head (first node visited)
//   after traversal, link head and tail to make circular
//   return head

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
