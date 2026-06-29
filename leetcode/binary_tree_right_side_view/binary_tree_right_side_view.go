package main

import "fmt"

// LeetCode 199 - Binary Tree Right Side View
//
// Problem:
//   Given the root of a binary tree, imagine yourself standing on the right
//   side of it. Return the values of the nodes you can see, ordered top to bottom.
//
// Example:
//   Input:
//     1
//    / \
//   2   3
//    \
//     5
//
//   Output: [1, 3, 5]
//
//   Explanation:
//   Level 0: rightmost visible = 1
//   Level 1: rightmost visible = 3  (3 hides 2)
//   Level 2: rightmost visible = 5  (only node at this level)
//
// Pseudo code:
//   BFS level-order; for each level record only the last node's value
//   return collected values

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func rightSideView(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	result := []int{}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[i]
			if i == size-1 {
				result = append(result, node.Val)
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		queue = queue[size:]
	}
	return result
}

func main() {
	//     1
	//    / \
	//   2   3
	//    \
	//     5
	root := &TreeNode{Val: 1,
		Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3},
	}
	fmt.Println(rightSideView(root)) // [1 3 5]
	fmt.Println(rightSideView(nil))  // []
}
