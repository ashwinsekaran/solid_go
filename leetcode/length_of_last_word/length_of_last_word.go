package main

import "fmt"

// LeetCode 58 - Length of Last Word
//
// Problem:
//   Given a string s consisting of words and spaces, return the length of the
//   last word. A word is a maximal substring of non-space characters.
//
// Example 1:
//   Input:  s = "Hello World"
//   Output: 5
//   Explanation: Last word is "World" which has length 5.
//
// Example 2:
//   Input:  s = "   fly me   to   the moon  "
//   Output: 4
//   Explanation: Last word is "moon"; trailing spaces are ignored.
//
// Pseudo code:
//   scan from end, skip trailing spaces
//   count non-space characters until space or start

func lengthOfLastWord(s string) int {
	i := len(s) - 1
	for i >= 0 && s[i] == ' ' {
		i--
	}
	count := 0
	for i >= 0 && s[i] != ' ' {
		count++
		i--
	}
	return count
}

func main() {
	fmt.Println(lengthOfLastWord("Hello World"))              // 5
	fmt.Println(lengthOfLastWord("   fly me   to   the moon  ")) // 4
	fmt.Println(lengthOfLastWord("luffy is still joyboy"))    // 6
}
