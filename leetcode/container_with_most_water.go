package leetcode

// LeetCode 11 - Container With Most Water
//
// Pseudo code:
//   l = 0, r = len(height)-1, maxWater = 0
//   while l < r:
//     water = min(height[l], height[r]) * (r - l)
//     maxWater = max(maxWater, water)
//     if height[l] < height[r]: l++
//     else: r--
//   return maxWater

func maxArea(height []int) int {
	l, r, maxWater := 0, len(height)-1, 0
	for l < r {
		h := height[l]
		if height[r] < h {
			h = height[r]
		}
		water := h * (r - l)
		if water > maxWater {
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
