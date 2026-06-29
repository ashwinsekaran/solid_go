package main

import "fmt"

// LeetCode 9 - Palindrome Number
//
// Problem:
//   Given an integer x, return true if it reads the same forward and backward.
//   Negative numbers are never palindromes.
//
// Example 1:
//   Input:  121   Output: true
//   Explanation: 121 reads the same left-to-right and right-to-left.
//
// Example 2:
//   Input:  -121  Output: false
//   Explanation: Reads -121 forward but 121- backward.
//
// Example 3:
//   Input:  10    Output: false
//   Explanation: Reads 10 forward but 01 backward (leading zero).
//
// Pseudo code:
//   negative numbers and multiples of 10 (except 0) are not palindromes
//   reverse the second half digit by digit
//   compare first half == reversed second half

func isPalindrome(x int) bool {
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}
	rev := 0
	for x > rev {
		rev = rev*10 + x%10
		x /= 10
	}
	return x == rev || x == rev/10
}

func main() {
	fmt.Println(isPalindrome(121))  // true
	fmt.Println(isPalindrome(-121)) // false
	fmt.Println(isPalindrome(10))   // false
	fmt.Println(isPalindrome(0))    // true
	fmt.Println(isPalindrome(1221)) // true
}
