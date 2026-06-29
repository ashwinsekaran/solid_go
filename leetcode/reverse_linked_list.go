package leetcode

// LeetCode 206 - Reverse Linked List
//
// Pseudo code:
//   prev = nil, curr = head
//   while curr != nil:
//     next = curr.Next
//     curr.Next = prev
//     prev = curr
//     curr = next
//   return prev

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	curr := head
	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}
	return prev
}
