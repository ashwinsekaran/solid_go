package main

import "fmt"

// LeetCode 9 - Palindrome Number
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
