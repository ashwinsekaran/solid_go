package main

import "fmt"

// LeetCode 103 - Binary Tree Zigzag Level Order Traversal
//
// Pseudo code:
//   BFS level by level
//   even levels: fill left-to-right
//   odd levels: fill right-to-left (write from end of slice)
//   return collected levels

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func zigzagLevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	result := [][]int{}
	queue := []*TreeNode{root}
	leftToRight := true
	for len(queue) > 0 {
		size := len(queue)
		level := make([]int, size)
		for i := 0; i < size; i++ {
			node := queue[i]
			if leftToRight {
				level[i] = node.Val
			} else {
				level[size-1-i] = node.Val
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		queue = queue[size:]
		result = append(result, level)
		leftToRight = !leftToRight
	}
	return result
}

func main() {
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
	fmt.Println(zigzagLevelOrder(root)) // [[3] [20 9] [15 7]]
	fmt.Println(zigzagLevelOrder(nil))  // []
}
