package leetcode

// LeetCode 58 - Length of Last Word
//
// Pseudo code:
//   i = len(s) - 1
//   skip trailing spaces
//   count chars until space or start
//   return count

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
