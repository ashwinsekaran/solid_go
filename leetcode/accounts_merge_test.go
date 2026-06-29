package leetcode

import (
	"sort"
	"reflect"
	"testing"
)

func TestAccountsMerge(t *testing.T) {
	accounts := [][]string{
		{"John", "johnsmith@mail.com", "john_newyork@mail.com"},
		{"John", "johnsmith@mail.com", "john00@mail.com"},
		{"Mary", "mary@mail.com"},
		{"John", "johnnybravo@mail.com"},
	}
	want := [][]string{
		{"John", "john00@mail.com", "john_newyork@mail.com", "johnsmith@mail.com"},
		{"John", "johnnybravo@mail.com"},
		{"Mary", "mary@mail.com"},
	}
	got := accountsMerge(accounts)

	normalize := func(res [][]string) {
		for _, acc := range res {
			sort.Strings(acc[1:])
		}
		sort.Slice(res, func(i, j int) bool {
			if res[i][0] != res[j][0] {
				return res[i][0] < res[j][0]
			}
			return res[i][1] < res[j][1]
		})
	}
	normalize(got)
	normalize(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("accountsMerge = %v, want %v", got, want)
	}
}
