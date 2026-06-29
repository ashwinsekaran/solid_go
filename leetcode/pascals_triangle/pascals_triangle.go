package main

import "fmt"

// LeetCode 118 - Pascal's Triangle
//
// Problem:
//   Given an integer numRows, return the first numRows of Pascal's triangle.
//   Each number is the sum of the two numbers directly above it.
//
// Example:
//   Input:  numRows = 5
//   Output: [[1],[1,1],[1,2,1],[1,3,3,1],[1,4,6,4,1]]
//
//   Explanation:
//       1
//      1 1
//     1 2 1
//    1 3 3 1
//   1 4 6 4 1
//   Each interior value = left-parent + right-parent from the row above.
//
// Pseudo code:
//   result = [[1]]
//   for i from 1 to numRows-1:
//     prev = result[i-1]
//     row = [1] + [prev[j-1]+prev[j] for each interior j] + [1]
//     result.append(row)
//   return result

func generate(numRows int) [][]int {
	result := [][]int{{1}}
	for i := 1; i < numRows; i++ {
		prev := result[i-1]
		row := []int{1}
		for j := 1; j < len(prev); j++ {
			row = append(row, prev[j-1]+prev[j])
		}
		row = append(row, 1)
		result = append(result, row)
	}
	return result
}

func main() {
	fmt.Println("numRows=5:", generate(5))
	fmt.Println("numRows=1:", generate(1))
	fmt.Println("numRows=3:", generate(3))
}
