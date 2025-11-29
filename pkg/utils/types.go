package utils

import "math"

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

type Vec2 struct {
	X float64
	Y float64
}

func (c Vec2) ToInt() (int, int) {
	return int(c.X), int(c.Y)
}

// Distance between two coordinates
func (c Vec2) Distance(other Vec2) float64 {
	dx := c.X - other.X
	dy := c.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type IntVec2 struct {
	X int
	Y int
}

// ToFloat Converts an integer coordinate to float coordinate
func (ic IntVec2) ToFloat() Vec2 {
	return Vec2{X: float64(ic.X), Y: float64(ic.Y)}
}
