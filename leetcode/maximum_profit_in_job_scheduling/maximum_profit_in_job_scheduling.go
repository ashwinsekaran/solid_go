package main

import (
	"fmt"
	"sort"
)

// LeetCode 1235 - Maximum Profit in Job Scheduling
//
// Problem:
//   You have n jobs. Job i runs from startTime[i] to endTime[i] and earns profit[i].
//   You can only do one job at a time (no overlapping). Return the max total profit.
//
// Example:
//   Input:  startTime=[1,2,3,3], endTime=[3,4,5,6], profit=[50,10,40,70]
//   Output: 120
//
//   Explanation:
//   Job 0: [1,3) earns 50
//   Job 3: [3,6) earns 70
//   Total: 50 + 70 = 120
//   (Cannot combine jobs 0+2 as job 2 starts at 3 which overlaps nothing,
//   but job 3 earns more than job 2.)
//
// Pseudo code:
//   sort jobs by end time
//   dp[i] = max profit considering first i jobs
//   for each job i: binary search for latest job j with end <= start[i]
//   dp[i] = max(dp[i-1], dp[j] + profit[i])

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
		if v := dp[lo] + p; v > dp[i] {
			dp[i] = v
		}
	}
	return dp[n]
}

func main() {
	fmt.Println(jobScheduling([]int{1, 2, 3, 3}, []int{3, 4, 5, 6}, []int{50, 10, 40, 70}))       // 120
	fmt.Println(jobScheduling([]int{1, 2, 3, 4, 6}, []int{3, 5, 10, 6, 9}, []int{20, 20, 100, 70, 60})) // 150
	fmt.Println(jobScheduling([]int{1, 1, 1}, []int{2, 3, 4}, []int{5, 6, 4}))                     // 6
}
