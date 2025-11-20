package pathfinding

import (
	"AI30_-_BlackFriday/pkg/utils"
)

// reconstructs the path from the final node
func ReconstructPath(node *Node) []utils.IntCoordinate {
	var path []utils.IntCoordinate
	current := node

	for current != nil {
		path = append([]utils.IntCoordinate{{X: current.X, Y: current.Y}}, path...)
		current = current.Parent
	}

	return path
}
