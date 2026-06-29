package main

import "fmt"

// LeetCode 291 - Word Pattern II
//
// Pseudo code:
//   backtracking with two maps: char->word and word->char
//   if char already mapped: check s starts with that word, recurse on rest
//   else: try every prefix of remaining s as a new word
//         skip if word already used; bind char<->word; recurse; unbind

func wordPatternMatch(pattern string, s string) bool {
	charToWord := make(map[byte]string)
	usedWords := make(map[string]bool)
	var bt func(pi, si int) bool
	bt = func(pi, si int) bool {
		if pi == len(pattern) && si == len(s) {
			return true
		}
		if pi == len(pattern) || si == len(s) {
			return false
		}
		c := pattern[pi]
		if word, ok := charToWord[c]; ok {
			if si+len(word) > len(s) || s[si:si+len(word)] != word {
				return false
			}
			return bt(pi+1, si+len(word))
		}
		for end := si + 1; end <= len(s); end++ {
			word := s[si:end]
			if usedWords[word] {
				continue
			}
			charToWord[c] = word
			usedWords[word] = true
			if bt(pi+1, end) {
				return true
			}
			delete(charToWord, c)
			usedWords[word] = false
		}
		return false
	}
	return bt(0, 0)
}

func main() {
	fmt.Println(wordPatternMatch("abab", "redblueredblue")) // true
	fmt.Println(wordPatternMatch("aaaa", "asdasdasdasd"))   // true
	fmt.Println(wordPatternMatch("aabb", "xyzxyz"))         // false
	fmt.Println(wordPatternMatch("ab", "aa"))               // false
}
