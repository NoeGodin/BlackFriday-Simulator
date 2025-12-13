package mapgenerator

import "AI30_-_BlackFriday/pkg/utils"

func containsPos(slice []utils.IntVec2, x, y int) bool {
	for _, pos := range slice {
		if pos.X == x && pos.Y == y {
			return true
		}
	}
	return false
}

func samePositions(a, b []utils.IntVec2) bool {
	if len(a) != len(b) {
		return false
	}

	for _, posA := range a {
		found := false
		for _, posB := range b {
			if posA.X == posB.X && posA.Y == posB.Y {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func hasFalse(m map[utils.IntVec2]bool) bool {
	for _, v := range m {
		if !v {
			return true
		}
	}
	return false
}
