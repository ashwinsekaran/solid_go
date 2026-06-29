package main

import "fmt"

// LeetCode 58 - Length of Last Word
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
