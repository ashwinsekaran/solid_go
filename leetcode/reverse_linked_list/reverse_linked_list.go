package main

import "fmt"

// LeetCode 206 - Reverse Linked List
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
