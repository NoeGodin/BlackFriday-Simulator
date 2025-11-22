package pathfinding

import "AI30_-_BlackFriday/pkg/utils"

type Path struct {
	waypoints []utils.IntCoordinate
	target    utils.IntCoordinate
}

// node for A* algorithme
type Node struct {
	X, Y   int
	GScore float64
	FScore float64
	Parent *Node
	index  int // Index in priority queue
}

// priority queue for A* nodes
type PriorityQueue []*Node
