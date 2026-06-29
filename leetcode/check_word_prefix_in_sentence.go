package leetcode

import "strings"

// LeetCode 1455 - Check If a Word Occurs As a Prefix of Any Word in a Sentence
//
// Pseudo code:
//   split sentence into words
//   for each word at index i:
//     if word starts with searchWord: return i+1 (1-indexed)
//   return -1

func isPrefixOfWord(sentence string, searchWord string) int {
	words := strings.Fields(sentence)
	for i, w := range words {
		if strings.HasPrefix(w, searchWord) {
			return i + 1
		}
	}
	return -1
}
