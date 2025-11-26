package pathfinding

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/logger"
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/utils"
	"container/heap"
)

// FindPath path between two points with A*
func FindPath(worldMap *Map.Map, startX, startY, targetX, targetY float64) (*Path, bool) {
	logger.Debugf("FindPath called from (%.2f,%.2f) to (%.2f,%.2f)", startX, startY, targetX, targetY)

	// If agent is stuck in a non-walkable position, find nearest walkable position first
	//TODO: on pourrait remove à la fin du projet, car ce sera plus possible d'avoir coord invalide
	if !worldMap.IsValidAndWalkable(startX, startY) {
		logger.Debugf("FindPath: Start position (%.2f,%.2f) is invalid or not walkable, finding nearest walkable position", startX, startY)
		nearestX, nearestY, found := findNearestWalkablePosition(worldMap, startX, startY)
		if !found {
			logger.Debugf("FindPath: No walkable position found near (%.2f,%.2f)", startX, startY)
			return nil, false
		}
		logger.Debugf("FindPath: Using nearest walkable position (%.2f,%.2f) instead of (%.2f,%.2f)", nearestX, nearestY, startX, startY)
		startX, startY = nearestX, nearestY
	}
	if !worldMap.IsValidAndWalkable(targetX, targetY) {
		logger.Debugf("FindPath: Target position (%.2f,%.2f) is invalid or not walkable", targetX, targetY)
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
		waypoints: waypoints,
		target:    utils.Coordinate{X: targetX, Y: targetY},
	}, true
}

func AStar(worldMap *Map.Map, startX, startY, targetX, targetY float64) ([]utils.Coordinate, bool) {
	logger.Debugf("A*: Finding path from (%.2f,%.2f) to (%.2f,%.2f)", startX, startY, targetX, targetY)

	// Initialize
	openSet := &PriorityQueue{}
	heap.Init(openSet)
	closedSet := make(map[[2]float64]bool)
	nodeMap := make(map[[2]float64]*Node)

	startNode := &Node{
		X:      startX,
		Y:      startY,
		GScore: 0, //c'est long mais f score c'est juste la distance entre 2 points
		FScore: utils.Coordinate{X: float64(startX), Y: float64(startY)}.Distance(utils.Coordinate{X: float64(targetX), Y: float64(targetY)}),
		Parent: nil,
	}

	heap.Push(openSet, startNode)
	nodeMap[[2]float64{startX, startY}] = startNode

	// Main A* algorithm loop
	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Node)
		currentKey := [2]float64{current.X, current.Y}

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
func processNeighbors(current *Node, targetX, targetY float64, worldMap *Map.Map,
	openSet *PriorityQueue, closedSet map[[2]float64]bool, nodeMap map[[2]float64]*Node) {

	directions := constants.MovementDirections

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

		neighborKey := [2]float64{nx, ny}
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
func handleNeighbor(nx, ny float64, tentativeG float64, current *Node,
	targetX, targetY float64, openSet *PriorityQueue,
	nodeMap map[[2]float64]*Node, neighborKey [2]float64) {

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
func findNearestWalkablePosition(worldMap *Map.Map, startX, startY float64) (float64, float64, bool) {
	queue := [][2]float64{{startX, startY}}
	visited := make(map[[2]float64]bool)
	visited[[2]float64{startX, startY}] = true

	// Search in expanding squares
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		x, y := current[0], current[1]

		if worldMap.IsValidAndWalkable(x, y) {
			return x, y, true
		}

		// Add neighbors to queue
		for _, dir := range [][2]float64{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			nx, ny := x+dir[0], y+dir[1]
			key := [2]float64{nx, ny}

			if !visited[key] && nx >= 0 && ny >= 0 && nx < float64(worldMap.Width) && ny < float64(worldMap.Height) {
				visited[key] = true
				queue = append(queue, [2]float64{nx, ny})
			}
		}
	}

	return 0, 0, false
}
