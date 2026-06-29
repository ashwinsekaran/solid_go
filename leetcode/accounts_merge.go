package leetcode

import "sort"

// LeetCode 721 - Accounts Merge
//
// Pseudo code:
//   Union-Find: group emails belonging to same account
//   emailToOwner map: email -> account name
//   for each account, union all emails with the first email
//   group emails by their root representative
//   sort each group, prepend account name
//   return result

func accountsMerge(accounts [][]string) [][]string {
	parent := make(map[string]string)
	var find func(x string) string
	find = func(x string) string {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(x, y string) {
		parent[find(x)] = find(y)
	}
	emailToName := make(map[string]string)
	for _, acc := range accounts {
		name := acc[0]
		for _, email := range acc[1:] {
			if _, ok := parent[email]; !ok {
				parent[email] = email
			}
			emailToName[email] = name
			union(acc[1], email)
		}
	}
	groups := make(map[string][]string)
	for email := range parent {
		root := find(email)
		groups[root] = append(groups[root], email)
	}
	var result [][]string
	for root, emails := range groups {
		sort.Strings(emails)
		merged := append([]string{emailToName[root]}, emails...)
		result = append(result, merged)
	}
	return result
}
