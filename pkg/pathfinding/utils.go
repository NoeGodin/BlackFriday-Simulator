package pathfinding

import (
	"AI30_-_BlackFriday/pkg/utils"
)

// reconstructs the path from the final node
func ReconstructPath(node *Node) []utils.Vec2 {
	var path []utils.Vec2
	current := node

	for current != nil {
		path = append([]utils.Vec2{{X: current.X, Y: current.Y}}, path...)
		current = current.Parent
	}

	return path
}
