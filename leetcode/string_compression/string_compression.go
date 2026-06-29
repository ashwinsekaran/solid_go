package main

import (
	"fmt"
	"strconv"
)

// LeetCode 443 - String Compression
//
// Pseudo code:
//   write pointer tracks where to write next
//   group consecutive identical chars; write char then count (if > 1)
//   return write pointer (new length)

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

func main() {
	c1 := []byte{'a', 'a', 'b', 'b', 'c', 'c', 'c'}
	n1 := compress(c1)
	fmt.Printf("len=%d result=%q\n", n1, string(c1[:n1])) // len=6 result="a2b2c3"

	c2 := []byte{'a'}
	n2 := compress(c2)
	fmt.Printf("len=%d result=%q\n", n2, string(c2[:n2])) // len=1 result="a"

	c3 := []byte{'a', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b'}
	n3 := compress(c3)
	fmt.Printf("len=%d result=%q\n", n3, string(c3[:n3])) // len=4 result="ab12"
}
