package utils

import "math"

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

type Coordinate struct {
	X float64
	Y float64
}

func (c Coordinate) ToInt() (int, int) {
	return int(c.X), int(c.Y)
}

// Distance between two coordinates
func (c Coordinate) Distance(other Coordinate) float64 {
	dx := c.X - other.X
	dy := c.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type IntCoordinate struct {
	X int
	Y int
}

// ToFloat Converts an integer coordinate to float coordinate
func (ic IntCoordinate) ToFloat() Coordinate {
	return Coordinate{X: float64(ic.X), Y: float64(ic.Y)}
}
