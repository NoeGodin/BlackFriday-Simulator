package mapgenerator

import (
	"AI30_-_BlackFriday/pkg/utils"
	"math/rand"
)

// Create a column of walls
func (m *MapLayout) FillColumn(x int) {
	for i := range m.height {
		m.mapLayout[i][x] = "W"
	}
}

// Create a row of walls
func (m *MapLayout) FillRow(y int) {
	for i := range m.width {
		m.mapLayout[y][i] = "W"
	}
}

func (m *MapLayout) GenerateDoors(nbDoors int) {
	for i := 0; i < nbDoors; i++ {
		// Choose between vertical or horizontal walls 
		isVertical := rand.Intn(2)
		if isVertical == 0 {
			// Cannot be the 2 lasts and firsts walls (avoid out of range + corner + first/last element)
			posY := rand.Intn(m.height - 4) + 2
			isFirstColumn := rand.Intn(2)
			
			// Choose in which column the door will be generated
			if isFirstColumn == 0 {
				if m.isAlreadyDoor(0, posY) {
					i--
					continue
				}

				m.mapLayout[posY][0] = "D"
			} else {
				if m.isAlreadyDoor(m.width - 1, posY) {
					i--
					continue
				}

				m.mapLayout[posY][m.width - 1] = "D"
			}
		} else {
			// Cannot be the 2 lasts and firsts walls (avoid out of range + corner + first/last element)
			posX := rand.Intn(len(m.mapLayout[0]) - 4) + 2
			isFirstRow := rand.Intn(2)

			// Choose in which row the door will be generated
			if isFirstRow == 0 {
				if m.isAlreadyDoor(posX, 0) {
					i--
					continue
				}

				m.mapLayout[0][posX] = "D"
			} else {
				if m.isAlreadyDoor(posX, m.height - 1) {
					i--
					continue
				}

				m.mapLayout[m.height - 1][posX] = "D"
			}
		}
	}
}

func (m *MapLayout) isAlreadyDoor(x, y int) bool {
	return m.mapLayout[y][x] == "D"
}

// Remove the walls around an x, y position
func (m *MapLayout) RemoveWallAround(x, y int) {
	if y != 0 {
		if m.mapLayout[y-1][x] == "W" {
			m.mapLayout[y-1][x] = ""
		}
	}
	if y != m.height - 1 {
		if m.mapLayout[y+1][x] == "W" {
			m.mapLayout[y+1][x] = ""
		}
	}

	if x != 0 {
		if m.mapLayout[y][x-1] == "W" {
			m.mapLayout[y][x-1] = ""
		}
	}
	if x != m.width - 1 {
		if m.mapLayout[y][x+1] == "W" {
			m.mapLayout[y][x+1] = ""
		}
	}
}

// Check if (x, y) is near to a door (in the 8 coordinates around)
func (m *MapLayout) isCloseToDoor(x, y int) bool {
	flag := false

	// "Direct" directions
	if y != 0 {
		flag = flag || m.mapLayout[y-1][x] == "D"
	}
	if y != m.height - 1 {
		flag = flag || m.mapLayout[y+1][x] == "D"
	}

	if x != 0 {
		flag = flag || m.mapLayout[y][x-1] == "D"
	}
	if x != m.width - 1 {
		flag = flag || m.mapLayout[y][x+1] == "D"
	}

	// Diagonals
	if y != 0 && x != 0 {
		flag = flag || m.mapLayout[y-1][x-1] == "D"
	}
	if y != m.height - 1 && x != 0 {
		flag = flag || m.mapLayout[y+1][x-1] == "D"
	}
	if y != 0 && x != m.width - 1 {
		flag = flag || m.mapLayout[y-1][x+1] == "D"
	}
	if y != m.height - 1 && x != m.width - 1 {
		flag = flag || m.mapLayout[y+1][x+1] == "D"
	}

	return flag 
}

// Takes walls positions and replace the wall to a shelf
func (m *MapLayout) GenerateShelves(nbShelves int) {
	for i := range nbShelves {
		letter := string(rune('a' + i))
		wallPositions := m.getAllWallsPos()
		randomPos := wallPositions[rand.Intn(len(wallPositions) - 1)]
		m.mapLayout[int(randomPos.Y)][int(randomPos.X)] = letter
	}
}

// Takes walls positions and replace the wall to a cashier
func (m *MapLayout) GenerateCashiers(nbCashiers int) {
	for range nbCashiers {
		wallPositions := m.getAllWallsPos()
		randomPos := wallPositions[rand.Intn(len(wallPositions)-1)]
		m.mapLayout[int(randomPos.Y)][int(randomPos.X)] = "C"
	}
}

// Generate randomly walls and try to avoid creating blocked zones
func (m *MapLayout) GenerateWalls(nbWalls int) {
    for i := 0; i < nbWalls; i++ {
        x := rand.Intn(m.width - 2) + 1
        y := rand.Intn(m.height - 2) + 1

        if m.canPlaceObstacle(x, y) {
            m.mapLayout[y][x] = "W"
        } else {
            i--
        }
    }
}

func (m *MapLayout) canPlaceObstacle(x, y int) bool {
    if m.mapLayout[y][x] != "" {
		return false
	}
	
    if m.isBlockingCorridor(x, y) || m.isCloseToDoor(x, y) {
        return false
    }

    freeSides := 0
    if y > 0 && m.mapLayout[y-1][x] != "W" {
		freeSides++
	}
    if y < m.height - 1 && m.mapLayout[y+1][x] != "W" {
		freeSides++
	}
    if x > 0 && m.mapLayout[y][x-1] != "W" { 
		freeSides++
	}
    if x < m.width - 1 && m.mapLayout[y][x+1] != "W" {
		freeSides++
	}

    return freeSides >= 2
}

func (m *MapLayout) isBlockingCorridor(x, y int) bool {
    dirs := [][]int{
        {0, 0}, {1, 0}, {0, 1}, {1, 1},
        {-1,0},{0,-1},{-1,-1},{-1,1},{1,-1},
    }

    for _, d := range dirs {
        cx := x + d[0]
        cy := y + d[1]

        if cx < 0 || cy < 0 || cx+1 >= m.width || cy+1 >= m.height {
            continue
        }

        if (m.mapLayout[cy][cx] == "W" || (cx == x && cy == y)) &&
           (m.mapLayout[cy][cx+1] == "W" || cx+2 >= m.width) &&
           (m.mapLayout[cy+1][cx] == "W" || cy+2 >= m.height) &&
           (m.mapLayout[cy+1][cx+1] == "W" || cy+2 >= m.height || cx+2 >= m.width) {
            return true
        }
    }

    return false
}

func (m *MapLayout) getAllWallsPos() []utils.Vec2 {
	var allWallsPos []utils.Vec2

	for y := 1; y < m.height-1; y++ {
		for x := 1; x < m.width-1; x++ {
			if m.mapLayout[y][x] == "W" {
				allWallsPos = append(allWallsPos, utils.Vec2{X: float64(x), Y: float64(y)})
			}
		}
	}

	return allWallsPos
}

// Loop in every case and if it is a door, remove the walls around it
func (m *MapLayout) RemoveAllDoorsWallsSurrounding() {
	for y := range m.height {
		for x := range m.width {
			if m.mapLayout[y][x] == "D" {
				m.RemoveWallAround(x, y)
			}
		}
	}
}