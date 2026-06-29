package main

import "fmt"

// LeetCode 234 - Palindrome Linked List
//
// Pseudo code:
//   find middle using slow/fast pointers
//   reverse second half in-place
//   compare first and reversed second half node by node

type ListNode struct {
	Val  int
	Next *ListNode
}

func isPalindromeList(head *ListNode) bool {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
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

func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func main() {
	fmt.Println(isPalindromeList(makeList([]int{1, 2, 2, 1}))) // true
	fmt.Println(isPalindromeList(makeList([]int{1, 2})))        // false
	fmt.Println(isPalindromeList(makeList([]int{1})))           // true
	fmt.Println(isPalindromeList(makeList([]int{1, 0, 0, 1}))) // true
}
