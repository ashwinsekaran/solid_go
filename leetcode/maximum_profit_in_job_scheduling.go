package leetcode

import "sort"

// LeetCode 1235 - Maximum Profit in Job Scheduling
//
// Pseudo code:
//   sort jobs by end time
//   dp[i] = max profit using first i jobs
//   for each job i (1-indexed):
//     binary search for latest job j whose end <= start[i]
//     dp[i] = max(dp[i-1], dp[j] + profit[i])
//   return dp[n]

func jobScheduling(startTime []int, endTime []int, profit []int) int {
	n := len(startTime)
	jobs := make([][3]int, n)
	for i := 0; i < n; i++ {
		jobs[i] = [3]int{endTime[i], startTime[i], profit[i]}
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i][0] < jobs[j][0] })

	dp := make([]int, n+1)
	ends := make([]int, n+1)
	for i, j := range jobs {
		ends[i+1] = j[0]
	}

	for i := 1; i <= n; i++ {
		start := jobs[i-1][1]
		p := jobs[i-1][2]
		// find latest job ending <= start
		lo, hi := 0, i-1
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if ends[mid] <= start {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		dp[i] = dp[i-1]
		if dp[lo]+p > dp[i] {
			dp[i] = dp[lo] + p
		}
	}
	return dp[n]
}
