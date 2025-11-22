package Simulation

import (
	"AI30_-_BlackFriday/pkg/logger"
	"AI30_-_BlackFriday/pkg/pathfinding"
	"AI30_-_BlackFriday/pkg/utils"
	"math"
)

type MovementManager struct {
	agent *ClientAgent
}

func NewMovementManager(agent *ClientAgent) *MovementManager {
	return &MovementManager{agent: agent}
}

// GenerateNewDestination generates a new random destination for the agent
func (mm *MovementManager) GenerateNewDestination() {
	targetX, targetY, found := mm.agent.env.Map.GetRandomFreeCoordinate()
	if !found {
		logger.Warnf("Agent %s: No free coordinate found for destination", mm.agent.id)
		return
	}

	currentX, currentY := mm.agent.coordinate.ToInt()

	logger.Debugf("Agent %s: Current position (%d,%d) - Wall:%t Shelf:%t Checkout:%t Door:%t",
		mm.agent.id, currentX, currentY,
		mm.agent.env.Map.IsWall(currentX, currentY),
		mm.agent.env.Map.IsShelf(currentX, currentY),
		mm.agent.env.Map.IsCheckout(currentX, currentY),
		mm.agent.env.Map.IsDoor(currentX, currentY))

	logger.Debugf("Agent %s: Target position (%d,%d) - Wall:%t Shelf:%t Checkout:%t Door:%t",
		mm.agent.id, targetX, targetY,
		mm.agent.env.Map.IsWall(targetX, targetY),
		mm.agent.env.Map.IsShelf(targetX, targetY),
		mm.agent.env.Map.IsCheckout(targetX, targetY),
		mm.agent.env.Map.IsDoor(targetX, targetY))

	if targetX == currentX && targetY == currentY {
		logger.Debugf("Agent %s: Generated same position as current, retrying", mm.agent.id)
		mm.GenerateNewDestination()
		return
	}

	logger.Debugf("Agent %s: Generating path from (%d,%d) to (%d,%d)", mm.agent.id, currentX, currentY, targetX, targetY)

	path, pathFound := pathfinding.FindPath(mm.agent.env.Map, currentX, currentY, targetX, targetY)
	if !pathFound {
		logger.Warnf("Agent %s: No path found to destination (%d,%d)", mm.agent.id, targetX, targetY)
		mm.agent.hasDestination = false
		return
	}

	logger.Debugf("Agent %s: Path found successfully!", mm.agent.id)

	mm.agent.currentPath = path
	mm.agent.targetX = targetX
	mm.agent.targetY = targetY
	mm.agent.hasDestination = true

	logger.Debugf("Agent %s: New destination set to (%d,%d) with %d waypoints", mm.agent.id, targetX, targetY, len(path.GetWaypoints()))
}

// FollowPath makes the agent follow its current path
func (mm *MovementManager) FollowPath() {

	// Check wayoint reached
	if waypoints := mm.agent.currentPath.GetWaypoints(); len(waypoints) > 0 {
		nextWaypoint := waypoints[0]
		dx := mm.agent.coordinate.X - float64(nextWaypoint.X)
		dy := mm.agent.coordinate.Y - float64(nextWaypoint.Y)
		distance := math.Sqrt(dx*dx + dy*dy) //TODO: could change with util function but used for coordinate type

		// if under value consider it reached
		//TODO: could use a constant variable but only used here and 0.6 work well
		if distance < 0.6 {
			logger.Debugf("Agent %s: Reached waypoint (%d,%d) with distance %.2f",
				mm.agent.id, nextWaypoint.X, nextWaypoint.Y, distance)
			mm.agent.currentPath.RemoveFirstWaypoint()
			if mm.agent.currentPath.IsComplete() {
				logger.Debugf("Agent %s: Reached final destination (%d,%d)", mm.agent.id, mm.agent.targetX, mm.agent.targetY)
				mm.agent.hasDestination = false
				mm.agent.currentPath = nil
				return
			}
		}
	}

	nextWaypoint, hasNext := mm.agent.currentPath.GetNextWaypoint()

	if !hasNext {
		logger.Debugf("Agent %s: Reached destination (%d,%d)", mm.agent.id, mm.agent.targetX, mm.agent.targetY)
		mm.agent.hasDestination = false
		mm.agent.currentPath = nil
		return
	}

	//direction next waypoint
	dx := float64(nextWaypoint.X) - mm.agent.coordinate.X
	dy := float64(nextWaypoint.Y) - mm.agent.coordinate.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance > 0 {
		mm.agent.dx = dx / distance
		mm.agent.dy = dy / distance
		logger.Debugf("Agent %s: Moving towards waypoint (%d,%d), direction: (%.2f,%.2f)",
			mm.agent.id, nextWaypoint.X, nextWaypoint.Y, mm.agent.dx, mm.agent.dy)
	} else {
		mm.agent.dx = 0
		mm.agent.dy = 0
	}
}

// CalculateDirection calculates agent direction based on dx, dy
func (mm *MovementManager) CalculateDirection() utils.Direction {
	dx, dy := mm.agent.dx, mm.agent.dy

	if dx == 0 && dy == 0 || dx == 0 && dy > 0 {
		return utils.SOUTH
	}
	if dx > 0 && dy == 0 {
		return utils.EAST
	}
	if dx < 0 && dy == 0 {
		return utils.WEST
	}
	if dx == 0 && dy < 0 {
		return utils.NORTH
	}
	// nord-est
	if dx > 0 && dy < 0 {
		return utils.NORTH
	}
	//nord-ouest
	if dx < 0 && dy < 0 {
		return utils.NORTH
	}
	//sud-est
	if dx > 0 && dy > 0 {
		return utils.SOUTH
	}
	//sud-ouest
	if dx < 0 && dy > 0 {
		return utils.SOUTH
	}
	return utils.SOUTH
}
