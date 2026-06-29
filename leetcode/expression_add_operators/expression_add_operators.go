package main

import (
	"fmt"
	"strconv"
)

// LeetCode 282 - Expression Add Operators
//
// Pseudo code:
//   backtracking over all splits of num string
//   track: current expression string, running value, last multiplied term
//   last term is needed to undo and redo multiplication when '*' is chosen
//   skip substrings with leading zeros (except "0" itself)
//   when entire string consumed and value == target: record expression

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
			if len(s) > 1 && s[0] == '0' { // no leading zeros
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

func main() {
	fmt.Println(addOperators("123", 6))   // [1+2+3 1*2*3]
	fmt.Println(addOperators("232", 8))   // [2*3+2 2+3*2]
	fmt.Println(addOperators("3456237490", 9191)) // []
	fmt.Println(addOperators("105", 5))   // [1*0+5 10-5]
}
