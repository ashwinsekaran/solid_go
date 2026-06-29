package main

import "fmt"

// LeetCode 126 - Word Ladder II
//
// Pseudo code:
//   BFS layer by layer from beginWord; build parent map (word -> list of predecessors)
//   stop once endWord is found (don't expand further — shortest path only)
//   DFS from endWord back to beginWord via parent map; reverse each path

func findLadders(beginWord string, endWord string, wordList []string) [][]string {
	wordSet := make(map[string]bool)
	for _, w := range wordList {
		wordSet[w] = true
	}
	if !wordSet[endWord] {
		return nil
	}
	parents := make(map[string][]string)
	layer := map[string]bool{beginWord: true}
	found := false
	for !found && len(layer) > 0 {
		nextLayer := map[string]bool{}
		for word := range layer {
			wordSet[word] = false
		}
		for word := range layer {
			bs := []byte(word)
			for i := 0; i < len(bs); i++ {
				orig := bs[i]
				for c := byte('a'); c <= byte('z'); c++ {
					if c == orig {
						continue
					}
					bs[i] = c
					next := string(bs)
					if wordSet[next] {
						nextLayer[next] = true
						parents[next] = append(parents[next], word)
						if next == endWord {
							found = true
						}
					}
					bs[i] = orig
				}
			}
		}
		layer = nextLayer
	}
	if !found {
		return nil
	}
	var result [][]string
	var dfs func(word string, path []string)
	dfs = func(word string, path []string) {
		if word == beginWord {
			tmp := make([]string, len(path))
			copy(tmp, path)
			for l, r := 0, len(tmp)-1; l < r; l, r = l+1, r-1 {
				tmp[l], tmp[r] = tmp[r], tmp[l]
			}
			result = append(result, tmp)
			return
		}
		for _, p := range parents[word] {
			dfs(p, append(path, p))
		}
	}
	dfs(endWord, []string{endWord})
	return result
}

func main() {
	paths := findLadders("hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"})
	for _, p := range paths {
		fmt.Println(p)
	}
	// [hit hot dot dog cog]
	// [hit hot lot log cog]

	fmt.Println(findLadders("hit", "cog", []string{"hot", "dot", "dog", "lot", "log"})) // nil
}
