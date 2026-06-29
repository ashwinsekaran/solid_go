package main

import "fmt"

// LeetCode 199 - Binary Tree Right Side View
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
