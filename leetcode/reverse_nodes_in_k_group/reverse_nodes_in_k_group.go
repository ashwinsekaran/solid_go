package main

import "fmt"

// LeetCode 25 - Reverse Nodes in k-Group
//
// Pseudo code:
//   check if at least k nodes remain; if not return head as-is
//   reverse exactly k nodes iteratively
//   old head.Next = reverseKGroup(remaining, k)
//   return new head

type ListNode struct {
	Val  int
	Next *ListNode
}

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
	printList(reverseKGroup(makeList([]int{1, 2, 3, 4, 5}), 2)) // 2 -> 1 -> 4 -> 3 -> 5
	printList(reverseKGroup(makeList([]int{1, 2, 3, 4, 5}), 3)) // 3 -> 2 -> 1 -> 4 -> 5
}
