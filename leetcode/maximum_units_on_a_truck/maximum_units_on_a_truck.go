package main

import (
	"fmt"
	"sort"
)

// LeetCode 1710 - Maximum Units on a Truck
//
// Problem:
//   You have a truck of a given size. boxTypes[i] = [numberOfBoxes, unitsPerBox].
//   You can choose any boxes to put on the truck as long as the total number
//   of boxes does not exceed truckSize. Return the maximum total units.
//
// Example:
//   Input:  boxTypes = [[1,3],[2,2],[3,1]], truckSize = 4
//   Output: 8
//
//   Explanation:
//   Take 1 box of type-1 (3 units) + 2 boxes of type-2 (4 units) + 1 box of
//   type-3 (1 unit) = 8.  We pick highest-unit boxes first.
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
