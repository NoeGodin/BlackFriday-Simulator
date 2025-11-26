package pathfinding

import (
	"AI30_-_BlackFriday/pkg/utils"
)

// reconstructs the path from the final node
func ReconstructPath(node *Node) []utils.Coordinate {
	var path []utils.Coordinate
	current := node

	for current != nil {
		path = append([]utils.Coordinate{{X: current.X, Y: current.Y}}, path...)
		current = current.Parent
	}

	return path
}
