package leetcode

// LeetCode 291 - Word Pattern II
//
// Pseudo code:
//   backtracking:
//     if pattern and s both exhausted: return true
//     if either exhausted: return false
//     char = pattern[0]
//     if char already mapped: check if s starts with mapped word, recurse on rest
//     else: try all prefixes of s as candidate word
//       if word not already used as value:
//         map char -> word, mark word used
//         recurse on pattern[1:] and s[len(word):]
//         if success: return true
//         unmap char, unmark word
//   return false

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
