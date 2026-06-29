package main

import "fmt"

// LeetCode 118 - Pascal's Triangle
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
