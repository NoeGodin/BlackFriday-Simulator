package main

import (
	"AI30_-_BlackFriday/pkg/utils"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	nbMaps := 5
	width := 30
	height := 20
	nbDoors := 30
	nbShelves := 25
	nbCashiers := 3
	nbWalls := 100

	mapLayout := make([][]string, height)
	for y := range height {
		mapLayout[y] = make([]string, width)
	}
	
	err := os.MkdirAll("maps/generated_maps", os.ModePerm)
	if err != nil {
		fmt.Printf("unable to create directories: %v\n", err)
		return
	}

	for i := range nbMaps {
		cleanMapLayout(mapLayout)

		fillRow(mapLayout, 0)
		fillRow(mapLayout, len(mapLayout)-1)
		fillColumn(mapLayout, 0)
		fillColumn(mapLayout, len(mapLayout[0])-1)
		generateDoors(mapLayout, nbDoors)
		generateWalls(mapLayout, nbWalls+nbCashiers+nbShelves)
		
		generateShelves(mapLayout, nbShelves)
		generateCashiers(mapLayout, nbCashiers)
		removeAllDoorsWallsSurrounding(mapLayout)

		mapStr := writeMap(mapLayout)	
		err = os.WriteFile("maps/generated_maps/map" + fmt.Sprint(i+1) + ".txt", []byte(mapStr), 0644)
		if err != nil {
			fmt.Printf("unable to write file: %v", err)
		}
	}
}

func cleanMapLayout(mapLayout [][]string) {
	for y := range mapLayout {
		for x := range mapLayout[y] {
			mapLayout[y][x] = ""
		}
	}
}

// Create a column of walls
func fillColumn(mapLayout [][]string, x int) {
	for i := range len(mapLayout) {
		mapLayout[i][x] = "W"
	}
}

// Create a row of walls
func fillRow(mapLayout [][]string, y int) {
	for i := range len(mapLayout[y]) {
		mapLayout[y][i] = "W"
	}
}

// Convert map layout from [][]string to a string for file
func writeMap(mapLayout [][]string) string {
	mapStr := ""

	for y := range len(mapLayout) {
		for x := range len(mapLayout[y]) {
			if(mapLayout[y][x] == "") {
				mapStr += " "
			} 
			mapStr += mapLayout[y][x]
		}
		mapStr += "\n"
	}

	return mapStr
}

func generateDoors(mapLayout [][]string, nbDoors int) {
	for i := 0; i < nbDoors; i++ {
		// Choose between vertical or horizontal walls 
		isVertical := rand.Intn(2)
		if isVertical == 0 {
			// Cannot be the 2 lasts and firsts walls (avoid out of range + corner + first/last element)
			posY := rand.Intn(len(mapLayout) - 4) + 2
			isFirstColumn := rand.Intn(2)
			
			// Choose in which column the door will be generated
			if isFirstColumn == 0 {
				if isAlreadyDoor(mapLayout, 0, posY) {
					i--
					continue
				}

				mapLayout[posY][0] = "D"
			} else {
				if isAlreadyDoor(mapLayout, len(mapLayout[0]) - 1, posY) {
					i--
					continue
				}

				mapLayout[posY][len(mapLayout[0]) - 1] = "D"
			}
		} else {
			// Cannot be the 2 lasts and firsts walls (avoid out of range + corner + first/last element)
			posX := rand.Intn(len(mapLayout[0]) - 4) + 2
			isFirstRow := rand.Intn(2)

			// Choose in which row the door will be generated
			if isFirstRow == 0 {
				if isAlreadyDoor(mapLayout, posX, 0) {
					i--
					continue
				}

				mapLayout[0][posX] = "D"
			} else {
				if isAlreadyDoor(mapLayout, posX, len(mapLayout) - 1) {
					i--
					continue
				}

				mapLayout[len(mapLayout) - 1][posX] = "D"
			}
		}
	}
}

func isAlreadyDoor(mapLayout [][]string, x, y int) bool {
	return mapLayout[y][x] == "D"
}

// Remove the walls around an x, y position
func removeWallAround(mapLayout [][]string, x, y int) {
	if y != 0 {
		if mapLayout[y-1][x] == "W" {
			mapLayout[y-1][x] = ""
		}
	}
	if y != len(mapLayout) - 1 {
		if mapLayout[y+1][x] == "W" {
			mapLayout[y+1][x] = ""
		}
	}

	if x != 0 {
		if mapLayout[y][x-1] == "W" {
			mapLayout[y][x-1] = ""
		}
	}
	if x != len(mapLayout[0]) - 1 {
		if mapLayout[y][x+1] == "W" {
			mapLayout[y][x+1] = ""
		}
	}
}

// Check if (x, y) is near to a door (in the 8 coordinates around)
func isCloseToDoor(mapLayout [][]string, x, y int) bool {
	flag := false
	if y != 0 {
		flag = flag || mapLayout[y-1][x] == "D"
	}
	if y != len(mapLayout) - 1 {
		flag = flag || mapLayout[y+1][x] == "D"
	}

	if x != 0 {
		flag = flag || mapLayout[y][x-1] == "D"
	}
	if x != len(mapLayout[0]) - 1 {
		flag = flag || mapLayout[y][x+1] == "D"
	}

	if y != 0 && x != 0 {
		flag = flag || mapLayout[y-1][x-1] == "D"
	}
	if y != len(mapLayout) - 1 && x != 0 {
		flag = flag || mapLayout[y+1][x-1] == "D"
	}
	if y != 0 && x != len(mapLayout[0]) - 1 {
		flag = flag || mapLayout[y-1][x+1] == "D"
	}
	if y != len(mapLayout) - 1 && x != len(mapLayout[0]) - 1 {
		flag = flag || mapLayout[y+1][x+1] == "D"
	}

	return flag 
}

// Takes walls positions and replace the wall to a shelf
func generateShelves(mapLayout [][]string, nbShelves int) {
	for i := 0; i < nbShelves; i++ {
		letter := string(rune('a' + i))
		wallPositions := getAllWallsPos(mapLayout)
		randomPos := wallPositions[rand.Intn(len(wallPositions)-1)]
		mapLayout[int(randomPos.Y)][int(randomPos.X)] = letter
	}
}

// Takes walls positions and replace the wall to a cashier
func generateCashiers(mapLayout [][]string, nbCashiers int) {
	for i := 0; i < nbCashiers; i++ {
		wallPositions := getAllWallsPos(mapLayout)
		randomPos := wallPositions[rand.Intn(len(wallPositions)-1)]
		mapLayout[int(randomPos.Y)][int(randomPos.X)] = "C"
	}
}

// Generate randomly walls and try to avoid creating blocked zones
func generateWalls(mapLayout [][]string, nbWalls int) {
    for i := 0; i < nbWalls; i++ {
        x := rand.Intn(len(mapLayout[0]) - 2) + 1
        y := rand.Intn(len(mapLayout) - 2) + 1

        if canPlaceObstacle(mapLayout, x, y) {
            mapLayout[y][x] = "W"
        } else {
            i--
        }
    }
}

func canPlaceObstacle(mapLayout [][]string, x, y int) bool {
    if mapLayout[y][x] != "" {
		return false
	}
	
    if isBlockingCorridor(mapLayout, x, y) || isCloseToDoor(mapLayout, x, y) {
        return false
    }

    freeSides := 0
    if y > 0 && mapLayout[y-1][x] != "W" {
		freeSides++
	}
    if y < len(mapLayout) - 1 && mapLayout[y+1][x] != "W" {
		freeSides++
	}
    if x > 0 && mapLayout[y][x-1] != "W" { 
		freeSides++
	}
    if x < len(mapLayout[0]) - 1 && mapLayout[y][x+1] != "W" {
		freeSides++
	}

    return freeSides >= 2
}

func isBlockingCorridor(mapLayout [][]string, x, y int) bool {
    dirs := [][]int{
        {0, 0}, {1, 0}, {0, 1}, {1, 1},
        {-1,0},{0,-1},{-1,-1},{-1,1},{1,-1},
    }

    for _, d := range dirs {
        cx := x + d[0]
        cy := y + d[1]

        if cx < 0 || cy < 0 || cx+1 >= len(mapLayout[0]) || cy+1 >= len(mapLayout) {
            continue
        }

        if (mapLayout[cy][cx] == "W" || (cx == x && cy == y)) &&
           (mapLayout[cy][cx+1] == "W" || cx+2 >= len(mapLayout[0])) &&
           (mapLayout[cy+1][cx] == "W" || cy+2 >= len(mapLayout)) &&
           (mapLayout[cy+1][cx+1] == "W" || cy+2 >= len(mapLayout) || cx+2 >= len(mapLayout[0])) {
            return true
        }
    }

    return false
}

func getAllWallsPos(mapLayout [][]string) []utils.Vec2 {
	var allWallsPos []utils.Vec2

	for y := 1; y < len(mapLayout)-1; y++ {
		for x := 1; x < len(mapLayout[0])-1; x++ {
			if mapLayout[y][x] == "W" {
				allWallsPos = append(allWallsPos, utils.Vec2{X: float64(x), Y: float64(y)})
			}
		}
	}

	return allWallsPos
}

func removeAllDoorsWallsSurrounding(mapLayout [][]string) {
	for y := range len(mapLayout) {
		for x := range len(mapLayout[0]) {
			if mapLayout[y][x] == "D" {
				removeWallAround(mapLayout, x, y)
			}
		}
	}
}