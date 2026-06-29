package leetcode

import "sort"

// LeetCode 1710 - Maximum Units on a Truck
//
// Pseudo code:
//   sort boxTypes by units per box descending
//   total = 0
//   for each [count, units] in boxTypes:
//     take = min(truckSize, count)
//     total += take * units
//     truckSize -= take
//     if truckSize == 0: break
//   return total

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
