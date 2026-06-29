package main

import (
	"fmt"
	"strconv"
	"strings"
)

// LeetCode 166 - Fraction to Recurring Decimal
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
