package leetcode

import "math"

// LeetCode 8 - String to Integer (atoi)
//
// Pseudo code:
//   skip leading whitespace
//   read optional sign (+ or -)
//   read digits, accumulate result
//   clamp to [INT_MIN, INT_MAX]
//   return sign * result

func myAtoi(s string) int {
	i, n := 0, len(s)
	for i < n && s[i] == ' ' {
		i++
	}
	sign := 1
	if i < n && (s[i] == '-' || s[i] == '+') {
		if s[i] == '-' {
			sign = -1
		}
		i++
	}
	result := 0
	for i < n && s[i] >= '0' && s[i] <= '9' {
		digit := int(s[i] - '0')
		if result > (math.MaxInt32-digit)/10 {
			if sign == 1 {
				return math.MaxInt32
			}
			return math.MinInt32
		}
		result = result*10 + digit
		i++
	}
	return sign * result
}
