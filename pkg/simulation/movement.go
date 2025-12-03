package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
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

	currentX, currentY := mm.agent.coordinate.X, mm.agent.coordinate.Y

	logger.Debugf("Agent %s: Current position (%.2f,%.2f) - Wall:%t Shelf:%t Checkout:%t Door:%t",
		mm.agent.id, currentX, currentY,
		mm.agent.env.Map.IsWall(currentX, currentY),
		mm.agent.env.Map.IsShelf(currentX, currentY),
		mm.agent.env.Map.IsCheckout(currentX, currentY),
		mm.agent.env.Map.IsDoor(currentX, currentY))

	logger.Debugf("Agent %s: Target position (%.2f,%.2f) - Wall:%t Shelf:%t Checkout:%t Door:%t",
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

	logger.Debugf("Agent %s: Generating path from (%.2f,%.2f) to (%.2f,%.2f)", mm.agent.id, currentX, currentY, targetX, targetY)

	path, pathFound := pathfinding.FindPath(mm.agent.env.Map, currentX, currentY, targetX, targetY)
	if !pathFound {
		logger.Warnf("Agent %s: No path found to destination (%.2f,%.2f)", mm.agent.id, targetX, targetY)
		mm.agent.hasDestination = false
		return
	}

	logger.Debugf("Agent %s: Path found successfully!", mm.agent.id)

	mm.agent.currentPath = path
	mm.agent.moveTargetX = targetX
	mm.agent.moveTargetY = targetY
	mm.agent.hasDestination = true

	logger.Debugf("Agent %s: New destination set to (%.2f,%.2f) with %d waypoints", mm.agent.id, targetX, targetY, len(path.GetWaypoints()))
}

func (mm *MovementManager) SetDestination(targetX, targetY float64) {
	currentX := mm.agent.coordinate.X
	currentY := mm.agent.coordinate.Y
	if targetX == currentX && targetY == currentY {
		mm.agent.hasDestination = false
		return
	}

	path, pathFound := pathfinding.FindPath(mm.agent.env.Map, currentX, currentY, targetX, targetY)
	if !pathFound {
		mm.agent.hasDestination = false
		return
	}

	mm.agent.currentPath = path
	mm.agent.moveTargetX = targetX
	mm.agent.moveTargetY = targetY
	mm.agent.hasDestination = true
}

// FollowPath makes the agent follow its current path
func (mm *MovementManager) FollowPath() {

	// Check wayoint reached
	if waypoints := mm.agent.currentPath.GetWaypoints(); len(waypoints) > 0 {
		nextWaypoint := waypoints[0]
		dx := mm.agent.coordinate.X - nextWaypoint.X
		dy := mm.agent.coordinate.Y - nextWaypoint.Y
		distance := math.Sqrt(dx*dx + dy*dy) //TODO: could change with util function but used for coordinate type

		// if under value consider it reached
		if distance < constants.WaypointReachedThreshold {
			logger.Debugf("Agent %s: Reached waypoint (%.2f,%.2f) with distance %.2f",
				mm.agent.id, nextWaypoint.X, nextWaypoint.Y, distance)
			mm.agent.currentPath.RemoveFirstWaypoint()
			if mm.agent.currentPath.IsComplete() {
				logger.Debugf("Agent %s: Reached final destination (%.2f,%.2f)", mm.agent.id, mm.agent.moveTargetX, mm.agent.moveTargetY)
				mm.agent.hasDestination = false
				mm.agent.currentPath = nil
				return
			}
		}
	}

	nextWaypoint, hasNext := mm.agent.currentPath.GetNextWaypoint()

	if !hasNext {
		logger.Debugf("Agent %s: Reached destination (%.2f,%.2f)", mm.agent.id, mm.agent.moveTargetX, mm.agent.moveTargetY)
		mm.agent.hasDestination = false
		mm.agent.currentPath = nil
		return
	}

	//direction next waypoint
	dx := nextWaypoint.X - mm.agent.coordinate.X
	dy := nextWaypoint.Y - mm.agent.coordinate.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance > 0 {
		mm.agent.dx = dx / distance
		mm.agent.dy = dy / distance
		mm.agent.desiredVelocity.X = mm.agent.dx * mm.agent.Speed
		mm.agent.desiredVelocity.Y = mm.agent.dy * mm.agent.Speed

		logger.Debugf("Agent %s: Moving towards waypoint (%.2f,%.2f), direction: (%.2f,%.2f)",
			mm.agent.id, nextWaypoint.X, nextWaypoint.Y, mm.agent.dx, mm.agent.dy)
	} else {
		mm.agent.dx = 0
		mm.agent.dy = 0
		mm.agent.desiredVelocity.X = 0
		mm.agent.desiredVelocity.Y = 0
	}
}

// CalculateDirection calculates agent direction based on dx, dy
func (mm *MovementManager) CalculateDirection() utils.Direction {
	dx, dy := mm.agent.dx, mm.agent.dy

	if dx == 0 && dy == 0 || dx == 0 && dy > 0 {
		return utils.SOUTH
	}
	if dx == 0 && dy < 0 {
		return utils.NORTH
	}
	if dx > 0 && dy == 0 {
		return utils.EAST
	}
	if dx < 0 && dy == 0 {
		return utils.WEST
	}

	// diagonal movement
	absDx := math.Abs(dx)
	absDy := math.Abs(dy)

	if absDx > absDy {
		if dx > 0 {
			return utils.EAST
		} else {
			return utils.WEST
		}
	} else {
		if dy > 0 {
			return utils.SOUTH
		} else {
			return utils.NORTH
		}
	}
}
