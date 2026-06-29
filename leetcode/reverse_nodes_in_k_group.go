package leetcode

// LeetCode 25 - Reverse Nodes in k-Group
//
// Pseudo code:
//   check if k nodes remain; if not return head as-is
//   reverse k nodes: standard iterative reversal
//   head.Next = reverseKGroup(remaining, k)
//   return new head (prev after reversal)

func reverseKGroup(head *ListNode, k int) *ListNode {
	curr := head
	count := 0
	for curr != nil && count < k {
		curr = curr.Next
		count++
	}
	if count < k {
		return head
	}
	var prev *ListNode
	curr = head
	for i := 0; i < k; i++ {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}
	head.Next = reverseKGroup(curr, k)
	return prev
}
