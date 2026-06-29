package leetcode

// LeetCode 234 - Palindrome Linked List
//
// Pseudo code:
//   find middle using slow/fast pointers
//   reverse second half of list
//   compare first half and reversed second half
//   return true if all match

func isPalindromeList(head *ListNode) bool {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	// reverse from slow
	var prev *ListNode
	for slow != nil {
		next := slow.Next
		slow.Next = prev
		prev = slow
		slow = next
	}
	l, r := head, prev
	for r != nil {
		if l.Val != r.Val {
			return false
		}
		l = l.Next
		r = r.Next
	}
	return true
}
