package main

import "fmt"

// LeetCode 815 - Bus Routes
//
// Problem:
//   You are given an array routes where routes[i] is a list of bus stops for
//   bus i. You start at source and want to reach target. You can board any
//   bus at any stop it visits. Return the minimum number of buses you need to
//   take, or -1 if it's impossible.
//
// Example:
//   Input:  routes = [[1,2,7],[3,6,7]], source = 1, target = 6
//   Output: 2
//
//   Explanation:
//   Board bus 0 at stop 1, ride to stop 7.
//   Board bus 1 at stop 7, ride to stop 6.
//   That is 2 buses total.
//
// Pseudo code:
//   build stop -> list of routes map
//   BFS by route (not individual stop)
//   start: enqueue all routes passing through source
//   for each route: visit all its stops; if target found return count
//   enqueue unvisited routes reachable from each stop

func numBusesToDestination(routes [][]int, source int, target int) int {
	if source == target {
		return 0
	}
	stopToRoutes := make(map[int][]int)
	for i, route := range routes {
		for _, stop := range route {
			stopToRoutes[stop] = append(stopToRoutes[stop], i)
		}
	}
	visitedRoutes := make([]bool, len(routes))
	visitedStops := make(map[int]bool)
	queue := []int{}
	for _, r := range stopToRoutes[source] {
		if !visitedRoutes[r] {
			visitedRoutes[r] = true
			queue = append(queue, r)
		}
	}
	buses := 1
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			for _, stop := range routes[queue[i]] {
				if stop == target {
					return buses
				}
				if visitedStops[stop] {
					continue
				}
				visitedStops[stop] = true
				for _, r := range stopToRoutes[stop] {
					if !visitedRoutes[r] {
						visitedRoutes[r] = true
						queue = append(queue, r)
					}
				}
			}
		}
		queue = queue[size:]
		buses++
	}
	return -1
}

func main() {
	fmt.Println(numBusesToDestination([][]int{{1, 2, 7}, {3, 6, 7}}, 1, 6))  // 2
	fmt.Println(numBusesToDestination([][]int{{1, 2, 7}}, 1, 1))              // 0
	fmt.Println(numBusesToDestination([][]int{{7, 12}, {4, 5, 15}, {6}, {15, 19}, {9, 12, 13}}, 15, 12)) // -1
}
