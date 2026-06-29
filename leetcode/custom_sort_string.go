package leetcode

import "sort"

// LeetCode 791 - Custom Sort String
//
// Pseudo code:
//   rank = map each char in order to its index
//   sort s characters: chars in order come first by rank, rest are 0
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
