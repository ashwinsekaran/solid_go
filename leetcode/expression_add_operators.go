package leetcode

import "strconv"

// LeetCode 282 - Expression Add Operators
//
// Pseudo code:
//   backtracking: try every split of num string
//   at each position, try +, -, * between accumulated expression and next number
//   track current value and last multiplied term (for handling * precedence)
//   if num has leading zero (except "0" itself): skip
//   when entire num consumed and value == target: add expression to result

func addOperators(num string, target int) []string {
	result := []string{}
	var bt func(idx int, path string, val, last int)
	bt = func(idx int, path string, val, last int) {
		if idx == len(num) {
			if val == target {
				result = append(result, path)
			}
			return
		}
		for end := idx + 1; end <= len(num); end++ {
			s := num[idx:end]
			if len(s) > 1 && s[0] == '0' {
				break
			}
			n, _ := strconv.ParseInt(s, 10, 64)
			curr := int(n)
			if idx == 0 {
				bt(end, s, curr, curr)
			} else {
				bt(end, path+"+"+s, val+curr, curr)
				bt(end, path+"-"+s, val-curr, -curr)
				bt(end, path+"*"+s, val-last+last*curr, last*curr)
			}
		}
	}
	bt(0, "", 0, 0)
	return result
}
