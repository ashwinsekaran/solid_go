package main

import (
	"fmt"
	"strings"
)

// LeetCode 1455 - Check If a Word Occurs As a Prefix of Any Word in a Sentence
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
