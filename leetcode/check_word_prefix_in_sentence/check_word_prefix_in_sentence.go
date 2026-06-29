package main

import (
	"fmt"
	"strings"
)

// LeetCode 1455 - Check If a Word Occurs As a Prefix of Any Word in a Sentence
//
// Problem:
//   Given a sentence and a searchWord, return the 1-indexed position of the
//   first word in the sentence that has searchWord as a prefix. Return -1 if none.
//
// Example 1:
//   Input:  sentence = "i love eating burger", searchWord = "burg"
//   Output: 4
//   Explanation: "burger" (word 4) starts with "burg".
//
// Example 2:
//   Input:  sentence = "this problem is an easy problem", searchWord = "pro"
//   Output: 2
//   Explanation: "problem" at position 2 starts with "pro" (position 6 also does,
//   but we return the first match).
//
// Pseudo code:
//   split sentence into words
//   return 1-indexed position of first word that starts with searchWord
//   return -1 if none found

func isPrefixOfWord(sentence string, searchWord string) int {
	for i, w := range strings.Fields(sentence) {
		if strings.HasPrefix(w, searchWord) {
			return i + 1
		}
	}
	return -1
}

func main() {
	fmt.Println(isPrefixOfWord("i love eating burger", "burg"))          // 4
	fmt.Println(isPrefixOfWord("this problem is an easy problem", "pro")) // 2
	fmt.Println(isPrefixOfWord("i am tired", "you"))                      // -1
}
