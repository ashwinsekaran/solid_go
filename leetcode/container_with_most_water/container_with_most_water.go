package main

import "fmt"

// LeetCode 11 - Container With Most Water
//
// Pseudo code:
//   two pointers l=0, r=end
//   area = min(height[l], height[r]) * (r - l)
//   advance the pointer with the shorter height
//   return max area seen

func maxArea(height []int) int {
	l, r, maxWater := 0, len(height)-1, 0
	for l < r {
		h := height[l]
		if height[r] < h {
			h = height[r]
		}
		if water := h * (r - l); water > maxWater {
			maxWater = water
		}
		if height[l] < height[r] {
			l++
		} else {
			r--
		}
	}
	return maxWater
}

func main() {
	fmt.Println(maxArea([]int{1, 8, 6, 2, 5, 4, 8, 3, 7})) // 49
	fmt.Println(maxArea([]int{1, 1}))                        // 1
	fmt.Println(maxArea([]int{4, 3, 2, 1, 4}))              // 16
}
