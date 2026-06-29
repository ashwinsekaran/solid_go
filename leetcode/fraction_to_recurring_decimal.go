package leetcode

import (
	"strconv"
	"strings"
)

// LeetCode 166 - Fraction to Recurring Decimal
//
// Pseudo code:
//   handle sign separately
//   integer part = numerator / denominator
//   remainder = numerator % denominator
//   if remainder == 0: return integer part
//   decimal part: simulate long division
//     track remainder -> position in decimal string
//     if remainder seen before: insert '(' at that position and append ')'
//     else: record position, decimal += remainder*10/denominator, remainder = remainder*10%denominator
//   return sign + intPart + "." + decimalPart

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
