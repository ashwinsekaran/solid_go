package main

import (
	"fmt"
	"sort"
)

// LeetCode 1710 - Maximum Units on a Truck
//
// Pseudo code:
//   sort boxTypes by units per box descending
//   greedily take as many of the highest-unit boxes as truck allows
//   return total units

func maximumUnits(boxTypes [][]int, truckSize int) int {
	sort.Slice(boxTypes, func(i, j int) bool {
		return boxTypes[i][1] > boxTypes[j][1]
	})
	total := 0
	for _, box := range boxTypes {
		take := box[0]
		if take > truckSize {
			take = truckSize
		}
		total += take * box[1]
		truckSize -= take
		if truckSize == 0 {
			break
		}
	}
	return total
}

func main() {
	fmt.Println(maximumUnits([][]int{{1, 3}, {2, 2}, {3, 1}}, 4))          // 8
	fmt.Println(maximumUnits([][]int{{5, 10}, {2, 5}, {4, 7}, {3, 9}}, 10)) // 91
}
