package leetcode

// LeetCode 103 - Binary Tree Zigzag Level Order Traversal
//
// Pseudo code:
//   result = []
//   queue = [root]
//   leftToRight = true
//   while queue not empty:
//     level = []
//     for each node in current queue level:
//       add node.val to level (append or prepend based on leftToRight)
//       enqueue children
//     result.append(level)
//     leftToRight = !leftToRight
//   return result

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
