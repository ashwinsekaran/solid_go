package main

import "fmt"

// LeetCode 205 - Isomorphic Strings
//
// Problem:
//   Two strings s and t are isomorphic if the characters in s can be replaced
//   to get t. Every character in s maps to exactly one character in t, and
//   no two characters map to the same character.
//
// Example 1:
//   Input:  s = "egg", t = "add"
//   Output: true
//   Explanation: 'e' -> 'a', 'g' -> 'd'
//
// Example 2:
//   Input:  s = "foo", t = "bar"
//   Output: false
//   Explanation: 'o' would need to map to both 'a' and 'r'
//
// Pseudo code:
//   maintain s->t and t->s character maps
//   for each pair (cs, ct):
//     if mapping conflicts in either direction: return false
//   return true

func isIsomorphic(s string, t string) bool {
	sToT := make(map[byte]byte)
	tToS := make(map[byte]byte)
	for i := 0; i < len(s); i++ {
		cs, ct := s[i], t[i]
		if mapped, ok := sToT[cs]; ok && mapped != ct {
			return false
		}
		if mapped, ok := tToS[ct]; ok && mapped != cs {
			return false
		}
		sToT[cs] = ct
		tToS[ct] = cs
	}
	return true
}

func main() {
	fmt.Println(isIsomorphic("egg", "add"))   // true
	fmt.Println(isIsomorphic("foo", "bar"))   // false
	fmt.Println(isIsomorphic("paper", "title")) // true
	fmt.Println(isIsomorphic("badc", "baba")) // false
}
