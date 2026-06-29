package leetcode

import "strconv"

// LeetCode 443 - String Compression
//
// Pseudo code:
//   write = 0, i = 0
//   while i < len(chars):
//     char = chars[i], count = 0
//     while i < len(chars) and chars[i] == char:
//       i++, count++
//     chars[write++] = char
//     if count > 1:
//       for each digit in str(count): chars[write++] = digit
//   return write

func compress(chars []byte) int {
	write, i := 0, 0
	for i < len(chars) {
		ch := chars[i]
		count := 0
		for i < len(chars) && chars[i] == ch {
			i++
			count++
		}
		chars[write] = ch
		write++
		if count > 1 {
			for _, d := range strconv.Itoa(count) {
				chars[write] = byte(d)
				write++
			}
		}
	}
	return write
}
