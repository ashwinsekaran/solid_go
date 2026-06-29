package main

import (
	"fmt"
	"strconv"
	"strings"
)

// LeetCode 166 - Fraction to Recurring Decimal
//
// Problem:
//   Given two integers numerator and denominator, return the fraction as a string.
//   If the decimal part is repeating, enclose the repeating part in parentheses.
//
// Example 1:
//   Input:  numerator = 1, denominator = 2
//   Output: "0.5"
//
// Example 2:
//   Input:  numerator = 1, denominator = 3
//   Output: "0.(3)"
//   Explanation: 1/3 = 0.3333... the '3' repeats forever.
//
// Example 3:
//   Input:  numerator = 4, denominator = 333
//   Output: "0.(012)"
//   Explanation: 4/333 = 0.012012012... the '012' block repeats.
//
//   Key insight: track each remainder's position in the decimal string.
//   When a remainder repeats, the digits from that position onward recur.
//
// Pseudo code:
//   handle sign and zero separately
//   integer part = num / den
//   remainder = num % den; if 0 return early
//   long division: track remainder -> position seen
//   if remainder repeats: insert '(' at that position and append ')'

func fractionToDecimal(numerator int, denominator int) string {
	if numerator == 0 {
		return "0"
	}
	var sb strings.Builder
	if (numerator < 0) != (denominator < 0) {
		sb.WriteByte('-')
	}
	num := abs(numerator)
	den := abs(denominator)
	sb.WriteString(strconv.Itoa(num / den))
	remainder := num % den
	if remainder == 0 {
		return sb.String()
	}
	sb.WriteByte('.')
	seen := make(map[int]int)
	var dec strings.Builder
	for remainder != 0 {
		if pos, ok := seen[remainder]; ok {
			part := dec.String()
			return sb.String() + part[:pos] + "(" + part[pos:] + ")"
		}
		seen[remainder] = dec.Len()
		remainder *= 10
		dec.WriteString(strconv.Itoa(remainder / den))
		remainder %= den
	}
	return sb.String() + dec.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	fmt.Println(fractionToDecimal(1, 2))    // "0.5"
	fmt.Println(fractionToDecimal(2, 1))    // "2"
	fmt.Println(fractionToDecimal(4, 333))  // "0.(012)"
	fmt.Println(fractionToDecimal(1, 6))    // "0.1(6)"
	fmt.Println(fractionToDecimal(-1, 2))   // "-0.5"
	fmt.Println(fractionToDecimal(-50, 8))  // "-6.25"
}
