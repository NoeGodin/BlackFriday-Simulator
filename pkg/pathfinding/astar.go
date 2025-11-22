package pathfinding

import (
	"AI30_-_BlackFriday/pkg/logger"
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/utils"
	"container/heap"
)

// FindPath path between two points with A*
func FindPath(worldMap *Map.Map, startX, startY, targetX, targetY int) (*Path, bool) {
	logger.Debugf("FindPath called from (%d,%d) to (%d,%d)", startX, startY, targetX, targetY)

	// If agent is stuck in a non-walkable position, find nearest walkable position first
	//TODO: on pourrait remove à la fin du projet, car ce sera plus possible d'avoir coord invalide
	if !worldMap.IsValidAndWalkable(startX, startY) {
		logger.Debugf("FindPath: Start position (%d,%d) is invalid or not walkable, finding nearest walkable position", startX, startY)
		nearestX, nearestY, found := findNearestWalkablePosition(worldMap, startX, startY)
		if !found {
			logger.Debugf("FindPath: No walkable position found near (%d,%d)", startX, startY)
			return nil, false
		}
		logger.Debugf("FindPath: Using nearest walkable position (%d,%d) instead of (%d,%d)", nearestX, nearestY, startX, startY)
		startX, startY = nearestX, nearestY
	}
	if !worldMap.IsValidAndWalkable(targetX, targetY) {
		logger.Debugf("FindPath: Target position (%d,%d) is invalid or not walkable", targetX, targetY)
		return nil, false
	}

	// A* algorithm
	logger.Debugf("FindPath: Calling A* algorithm...")
	waypoints, found := AStar(worldMap, startX, startY, targetX, targetY)

	if !found {
		logger.Debugf("FindPath: No path found")
		return nil, false
	}

	logger.Debugf("FindPath: Created path with %d waypoints", len(waypoints))
	return &Path{
		Waypoints: waypoints,
		Target:    utils.IntCoordinate{X: targetX, Y: targetY},
	}, true
}

func AStar(worldMap *Map.Map, startX, startY, targetX, targetY int) ([]utils.IntCoordinate, bool) {
	logger.Debugf("A*: Finding path from (%d,%d) to (%d,%d)", startX, startY, targetX, targetY)

	// Initialize
	openSet := &PriorityQueue{}
	heap.Init(openSet)
	closedSet := make(map[[2]int]bool)
	nodeMap := make(map[[2]int]*Node)

	startNode := &Node{
		X:      startX,
		Y:      startY,
		GScore: 0, //c'est long mais f score c'est juste la distance entre 2 points
		FScore: utils.Coordinate{X: float64(startX), Y: float64(startY)}.Distance(utils.Coordinate{X: float64(targetX), Y: float64(targetY)}),
		Parent: nil,
	}

	heap.Push(openSet, startNode)
	nodeMap[[2]int{startX, startY}] = startNode

	// Main A* algorithm loop
	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Node)
		currentKey := [2]int{current.X, current.Y}

		// is destination reached
		if current.X == targetX && current.Y == targetY {
			logger.Debugf("A*: Path found!")
			return ReconstructPath(current), true
		}

		closedSet[currentKey] = true

		// Explore neighbors
		processNeighbors(current, targetX, targetY, worldMap, openSet, closedSet, nodeMap)
	}

	logger.Debugf("A*: No path found")
	return nil, false
}

// processNeighbors processes all neighbors of a node
func processNeighbors(current *Node, targetX, targetY int, worldMap *Map.Map,
	openSet *PriorityQueue, closedSet map[[2]int]bool, nodeMap map[[2]int]*Node) {

	directions := utils.GetMovementDirections()

	for _, dir := range directions {
		nx, ny := current.X+dir[0], current.Y+dir[1]

		// Validity checks
		if !worldMap.IsValidAndWalkable(nx, ny) {
			continue
		}

		// For diagonal movement, check that both intermediate cells are walkable
		// to prevent cutting through corners
		if dir[0] != 0 && dir[1] != 0 {
			// Check horizontal intermediate cell
			if !worldMap.IsValidAndWalkable(current.X+dir[0], current.Y) {
				continue
			}
			// Check vertical intermediate cell
			if !worldMap.IsValidAndWalkable(current.X, current.Y+dir[1]) {
				continue
			}
		}

		neighborKey := [2]int{nx, ny}
		if closedSet[neighborKey] {
			continue
		}

		// Calculate movement cost
		moveCost := utils.CalculateMovementCost(dir[0], dir[1])
		tentativeG := current.GScore + moveCost

		// Handle neighbor
		handleNeighbor(nx, ny, tentativeG, current, targetX, targetY,
			openSet, nodeMap, neighborKey)
	}
}

// handleNeighbor handles adding or updating a neighbor
func handleNeighbor(nx, ny int, tentativeG float64, current *Node,
	targetX, targetY int, openSet *PriorityQueue,
	nodeMap map[[2]int]*Node, neighborKey [2]int) {

	neighbor, exists := nodeMap[neighborKey]
	if !exists {
		neighbor = &Node{
			X:      nx,
			Y:      ny,
			GScore: tentativeG,
			FScore: tentativeG + utils.Coordinate{X: float64(nx), Y: float64(ny)}.Distance(utils.Coordinate{X: float64(targetX), Y: float64(targetY)}),
			Parent: current,
		}
		nodeMap[neighborKey] = neighbor
		heap.Push(openSet, neighbor)
	} else if tentativeG < neighbor.GScore {
		// Update existing node with better path
		neighbor.GScore = tentativeG
		neighbor.FScore = tentativeG + utils.Coordinate{X: float64(nx), Y: float64(ny)}.Distance(utils.Coordinate{X: float64(targetX), Y: float64(targetY)})
		neighbor.Parent = current
		openSet.Update(neighbor, neighbor.FScore)
	}
}

// findNearestWalkablePosition finds the nearest walkable position
func findNearestWalkablePosition(worldMap *Map.Map, startX, startY int) (int, int, bool) {
	queue := [][2]int{{startX, startY}}
	visited := make(map[[2]int]bool)
	visited[[2]int{startX, startY}] = true

	// Search in expanding squares
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		x, y := current[0], current[1]

		if worldMap.IsValidAndWalkable(x, y) {
			return x, y, true
		}

		// Add neighbors to queue
		for _, dir := range [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			nx, ny := x+dir[0], y+dir[1]
			key := [2]int{nx, ny}

			if !visited[key] && nx >= 0 && ny >= 0 && nx < worldMap.Width && ny < worldMap.Height {
				visited[key] = true
				queue = append(queue, [2]int{nx, ny})
			}
		}
	}

	return 0, 0, false
}
