package leetcode

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindLadders(t *testing.T) {
	got := findLadders("hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"})
	want := [][]string{
		{"hit", "hot", "dot", "dog", "cog"},
		{"hit", "hot", "lot", "log", "cog"},
	}
	sortPaths := func(paths [][]string) {
		for _, p := range paths {
			sort.Strings(p)
		}
		sort.Slice(paths, func(i, j int) bool {
			return paths[i][0] < paths[j][0]
		})
	}
	sortPaths(got)
	sortPaths(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("findLadders = %v, want %v", got, want)
	}

	noPath := findLadders("hit", "cog", []string{"hot", "dot", "dog", "lot", "log"})
	if noPath != nil {
		t.Errorf("expected nil when no path, got %v", noPath)
	}
}
