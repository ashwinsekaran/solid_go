package main

import "fmt"

// LeetCode 205 - Isomorphic Strings
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
