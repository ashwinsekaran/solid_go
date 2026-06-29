package leetcode

// LeetCode 815 - Bus Routes
//
// Pseudo code:
//   build map: stop -> list of routes that pass through it
//   BFS by route (not by stop)
//   queue starts with all routes passing through source stop
//   for each route, visit all stops; if target found return transfers+1
//   enqueue new routes from unvisited stops
//   return -1 if unreachable

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
			route := queue[i]
			for _, stop := range routes[route] {
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
