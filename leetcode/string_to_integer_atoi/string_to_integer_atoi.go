package main

import (
	"fmt"
	"math"
)

// LeetCode 8 - String to Integer (atoi)
//
// Pseudo code:
//   skip leading whitespace
//   read optional sign
//   read digits, build integer
//   clamp to [INT32_MIN, INT32_MAX] on overflow

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

func main() {
	fmt.Println(myAtoi("42"))              // 42
	fmt.Println(myAtoi("   -42"))          // -42
	fmt.Println(myAtoi("4193 with words")) // 4193
	fmt.Println(myAtoi("words and 987"))   // 0
	fmt.Println(myAtoi("-91283472332"))    // -2147483648
	fmt.Println(myAtoi("2147483648"))      // 2147483647
}
