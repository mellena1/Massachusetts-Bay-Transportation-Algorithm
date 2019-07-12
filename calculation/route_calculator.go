package calculation

var numberOfRoutes = 0

// FindBestRoute finds the fastest route to traverse every stop, every stop must have an edge to every other stop
func FindBestRoute(stops []string) ([]string, int) {
	route := make([]string, 0)
	return findBestRouteHelper(route, stops)
}

func findBestRouteHelper(curRoute []string, stopsLeft []string) ([]string, int) {
	if len(stopsLeft) == 0 {
		numberOfRoutes++
		return curRoute, findRouteTime(curRoute)
	}

	var bestRoute []string
	bestTime := 999999

	for i := range stopsLeft {
		route, time := findBestRouteHelper(append(curRoute, stopsLeft[i]), removeIndex(i, stopsLeft))
		if time < bestTime {
			bestRoute = route
			bestTime = time
		}
	}

	return bestRoute, bestTime
}

func removeIndex(index int, list []string) []string {
	list[index] = list[len(list)-1]
	return list[:len(list)-1]
}

func findRouteTime(route []string) int {
	time := 0
	for i := 0; i < len(route)-1; i++ {
		time += findEdgeTime(route[i], route[i+1])
	}
	return time
}

func findEdgeTime(stopA string, stopB string) int {
	return 1 // TODO make this meaningful
}
