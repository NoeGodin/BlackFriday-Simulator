package utils

// Abs Absolute value
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetMovementDirections all possible movemnt
func GetMovementDirections() [][2]int {
	return [][2]int{
		{0, -1}, {1, 0}, {0, 1}, {-1, 0}, // Nord, Est, Sud, Ouest
		{1, -1}, {1, 1}, {-1, 1}, {-1, -1}, // Diagonales
	}
}

// CalculateMovementCost movement cost
func CalculateMovementCost(dx, dy int) float64 {
	if dx != 0 && dy != 0 {
		return 1.414 // sqrt(2) for diagonals
	}
	return 1.0 // Horizontal/vertical movement
}
