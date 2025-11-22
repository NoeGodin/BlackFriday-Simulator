package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/logger"
	"math"
)

type StuckDetector struct {
	agent *ClientAgent
}

func NewStuckDetector(agent *ClientAgent) *StuckDetector {
	return &StuckDetector{agent: agent}
}

// DetectAndResolve detects if agent is stuck
func (sd *StuckDetector) DetectAndResolve() {
	dx := sd.agent.coordinate.X - sd.agent.lastPosition.X
	dy := sd.agent.coordinate.Y - sd.agent.lastPosition.Y
	distanceMoved := math.Sqrt(dx*dx + dy*dy)

	//TODO:peut-être adapté le threshold avec base speed mais en même temps s'il est stuck...
	if distanceMoved < constants.StuckDistanceThreshold && sd.agent.hasDestination {
		sd.agent.stuckCounter++
		logger.Debugf("Agent %s: Potentially stuck (counter: %d, distance moved: %.3f)",
			sd.agent.id, sd.agent.stuckCounter, distanceMoved)
	} else {
		sd.agent.stuckCounter = 0
	}

	if sd.agent.stuckCounter > constants.StuckCounterThreshold {
		sd.resolveStuckState()
	}

	// Update last position
	sd.agent.lastPosition = sd.agent.coordinate
}

func (sd *StuckDetector) resolveStuckState() {
	currentX, currentY := sd.agent.coordinate.ToInt()

	// Check if agent is on non-walkable tile (shelf, wall, etc.)
	if !sd.agent.env.Map.IsWalkable(currentX, currentY) {
		logger.Warnf("Agent %s: Stuck on non-walkable tile (%d,%d)! Relocating to nearest free position", sd.agent.id, currentX, currentY)
		sd.relocateAgent()
	} else {
		logger.Warnf("Agent %s: Stuck detected! Regenerating destination", sd.agent.id)
	}

	//TODO: changer la logique, mais si on lui remet le même path il va se re stuck
	sd.agent.hasDestination = false
	sd.agent.currentPath = nil

	sd.agent.stuckCounter = 0
}

// relocateAgent moves the agent to a free position
func (sd *StuckDetector) relocateAgent() {
	currentX, currentY := sd.agent.coordinate.ToInt()

	newX, newY, found := FindNearestFreePosition(sd.agent.env, currentX, currentY)
	if found {
		logger.Infof("Agent %s: Relocating from (%d,%d) to (%d,%d)", sd.agent.id, currentX, currentY, newX, newY)
		sd.agent.coordinate.X = float64(newX)
		sd.agent.coordinate.Y = float64(newY)
		sd.agent.lastPosition = sd.agent.coordinate
	} else {
		// If no free position found, use GetRandomFreeCoordinate as last resort
		//TODO: faudrait vraiment changer ça mais pour l'instant ça suffira
		newX, newY, found = sd.agent.env.Map.GetRandomFreeCoordinate()
		if found {
			logger.Warnf("Agent %s: No nearby free position found, teleporting to random position (%d,%d)", sd.agent.id, newX, newY)
			sd.agent.coordinate.X = float64(newX)
			sd.agent.coordinate.Y = float64(newY)
			sd.agent.lastPosition = sd.agent.coordinate
		}
	}
}
