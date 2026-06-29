package main

import "fmt"

// LeetCode 206 - Reverse Linked List
//
// Problem:
//   Given the head of a singly linked list, reverse the list and return the
//   new head.
//
// Example:
//   Input:  1 -> 2 -> 3 -> 4 -> 5
//   Output: 5 -> 4 -> 3 -> 2 -> 1
//
//   Explanation:
//   Iterate once, redirecting each node's Next pointer to its predecessor.
//   prev=nil, curr=1: curr.Next=nil, prev=1, curr=2
//   prev=1,   curr=2: curr.Next=1,   prev=2, curr=3  ... and so on.
//
// Pseudo code:
//   prev = nil, curr = head
//   while curr != nil:
//     save next; point curr.Next to prev; advance prev and curr
//   return prev (new head)

type ListNode struct {
	Val  int
	Next *ListNode
}

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

func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func printList(head *ListNode) {
	for head != nil {
		fmt.Printf("%d", head.Val)
		if head.Next != nil {
			fmt.Print(" -> ")
		}
		head = head.Next
	}
	fmt.Println()
}

func main() {
	printList(reverseList(makeList([]int{1, 2, 3, 4, 5}))) // 5 -> 4 -> 3 -> 2 -> 1
	printList(reverseList(makeList([]int{1, 2})))           // 2 -> 1
	printList(reverseList(makeList([]int{1})))              // 1
}
