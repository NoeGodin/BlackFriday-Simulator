package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/logger"
	"AI30_-_BlackFriday/pkg/utils"
	"fmt"
	"math"
)

// FindNearestFreePosition finds the nearest free position around the given position
func FindNearestFreePosition(env *Environment, centerX, centerY float64) (float64, float64, bool) {
	// spiral research around center tile
	maxRadius := 10 // Limit search to a radius of 10 cells

	for radius := 1; radius <= maxRadius; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				// Only check the square perimeter (avoid re-checking the center)
				if utils.Abs(dx) != radius && utils.Abs(dy) != radius {
					continue
				}

				x := centerX + float64(dx)
				y := centerY + float64(dy)

				// position is within map bounds
				if x >= 0 && x < float64(env.Map.Width) && y >= 0 && y < float64(env.Map.Height) {
					if env.Map.IsWalkable(x, y) {
						logger.Debugf("Found free position (%.2f,%.2f) at radius %d", x, y, radius)
						return x, y, true
					}
				}
			}
		}
	}

	logger.Warnf("No free position found within radius %d", maxRadius)
	return 0, 0, false
}

func FindNearestElementPosition(env *Environment, a Agent, elementType constants.ElementType) (float64, float64, bool) {
	agentX, agentY := a.Coordinate().X, a.Coordinate().Y
	agentCoords := [2]float64{agentX, agentY}

	minDist := math.MaxFloat64
	nearestElement := [2]float64{-1, -1}

	switch elementType {
	case "shelf":
		agent, isClientAgent := a.(*ClientAgent)
		for k := range env.Map.ShelfData {
			if isClientAgent {
				_, alreadyVisited := agent.visitedShelves[k]
				if alreadyVisited {
					continue
				}
			}
			tempDist := math.Hypot(agentCoords[0]-k[0], agentCoords[1]-k[1])
			if minDist > tempDist {
				minDist = tempDist
				nearestElement = k
			}
		}

	case "C":
		for _, v := range env.Map.CheckoutZones {
			vx, vy := float64(v[0]), float64(v[1])
			tempDist := math.Hypot(agentCoords[0]-vx, agentCoords[1]-vy)
			if minDist > tempDist {
				minDist = tempDist
				nearestElement = [2]float64{vx, vy}
			}
		}
	case "D":
		for _, v := range env.Map.Doors {
			vx, vy := float64(v[0]), float64(v[1])
			tempDist := math.Hypot(agentCoords[0]-vx, agentCoords[1]-vy)
			if minDist > tempDist {
				minDist = tempDist
				nearestElement = [2]float64{vx, vy}
			}
		}
	default:
		logger.Warnf("This element (%s) cannot be used as a parameter for this function", elementType)
		return 0, 0, false
	}

	if nearestElement[0] == -1 && nearestElement[1] == -1 {
		logger.Warnf("No element (%s) position found", elementType)
		return 0, 0, false
	}

	fmt.Println(agentCoords, nearestElement)
	return nearestElement[0], nearestElement[1], true
}

func FindWalkablePositionNearbyElement(env *Environment, a Agent, elementType constants.ElementType) (float64, float64, bool) {
	nearestElementX, nearestElementY, res := FindNearestElementPosition(env, a, elementType)

	if res != true {
		logger.Warnf("No position found for this element (%s), cannot find nearest free position", elementType)
		return 0, 0, false
	}

	targetX, targetY, res := FindNearestFreePosition(env, nearestElementX, nearestElementY)
	if res != true {
		logger.Warnf("Cannot find nearest free position around element %s", elementType)
	}
	fmt.Printf("Walkable tile : [%.2f %.2f]\n", targetX, targetY)
	return targetX, targetY, res
}
