package utils

import "math"

// Abs Absolute value
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CalculateMovementCost movement cost
func CalculateMovementCost(dx, dy float64) float64 {
	if dx != 0 && dy != 0 {
		return 1.414 // sqrt(2) for diagonals
	}
	return 1.0 // Horizontal/vertical movement
}

func EuclideanDistance(c1, c2 [2]float64) float64 {
	dx := (c1[0] - c2[0])
	dy := (c1[1] - c2[1])
	return math.Sqrt(dx*dx + dy*dy)
}
