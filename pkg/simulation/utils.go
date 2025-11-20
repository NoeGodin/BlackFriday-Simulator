package Simulation

import (
	"AI30_-_BlackFriday/pkg/logger"
	"AI30_-_BlackFriday/pkg/utils"
)

// FindNearestFreePosition finds the nearest free position around the given position
func FindNearestFreePosition(env *Environment, centerX, centerY int) (int, int, bool) {
	// spiral research around center tile
	maxRadius := 10 // Limit search to a radius of 10 cells

	for radius := 1; radius <= maxRadius; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				// Only check the square perimeter (avoid re-checking the center)
				if utils.Abs(dx) != radius && utils.Abs(dy) != radius {
					continue
				}

				x := centerX + dx
				y := centerY + dy

				// position is within map bounds
				if x >= 0 && x < env.Map.Width && y >= 0 && y < env.Map.Height {
					if env.Map.IsWalkable(x, y) {
						logger.Debugf("Found free position (%d,%d) at radius %d", x, y, radius)
						return x, y, true
					}
				}
			}
		}
	}

	logger.Warnf("No free position found within radius %d", maxRadius)
	return 0, 0, false
}
