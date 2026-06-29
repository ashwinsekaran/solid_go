package main

import (
	"fmt"
	"sort"
)

// LeetCode 791 - Custom Sort String
//
// Problem:
//   order is a string with no duplicates defining a custom character order.
//   Rearrange the characters of s so that they match the order in order.
//   Characters not in order can appear anywhere.
//
// Example:
//   Input:  order = "cba", s = "abcd"
//   Output: "cbad"  (or "dcba" — any order where c before b before a)
//
//   Explanation:
//   'c','b','a' appear in order; 'd' is not in order so it can go anywhere.
//
// Pseudo code:
//   assign rank to each char in order string
//   sort s by rank (chars not in order get rank 0, appear first/last in ties)
//   return sorted string

func customSortString(order string, s string) string {
	rank := make(map[rune]int)
	for i, c := range order {
		rank[c] = i + 1
	}
	bs := []rune(s)
	sort.Slice(bs, func(i, j int) bool {
		return rank[bs[i]] < rank[bs[j]]
	})
	return string(bs)
}

func main() {
	fmt.Println(customSortString("cba", "abcd"))   // dcba or similar
	fmt.Println(customSortString("bcafg", "abcd")) // bca d
}
