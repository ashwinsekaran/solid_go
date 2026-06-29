package main

import (
	"fmt"
	"sort"
)

// LeetCode 1366 - Rank Teams by Votes
//
// Pseudo code:
//   for each team count votes at each position
//   sort teams: compare position-by-position vote counts descending
//   tie-break alphabetically

func rankTeams(votes []string) string {
	n := len(votes[0])
	counts := make(map[byte][]int)
	for _, v := range votes[0] {
		counts[byte(v)] = make([]int, n)
	}
	for _, vote := range votes {
		for pos, c := range vote {
			counts[byte(c)][pos]++
		}
	}
	teams := []byte(votes[0])
	sort.Slice(teams, func(i, j int) bool {
		ci, cj := counts[teams[i]], counts[teams[j]]
		for k := 0; k < n; k++ {
			if ci[k] != cj[k] {
				return ci[k] > cj[k]
			}
		}
		return teams[i] < teams[j]
	})
	return string(teams)
}

func main() {
	fmt.Println(rankTeams([]string{"ABC", "ACB", "ABC", "ACB", "ACB"}))          // "ACB"
	fmt.Println(rankTeams([]string{"WXYZ", "XWYZ", "ZWYX"}))                     // "XWYZ"
	fmt.Println(rankTeams([]string{"BCA", "CAB", "CBA", "ABC", "ACB", "BAC"}))   // "ABC"
}
