package Map

import (
	"math/rand"
)

func NewMap(width int, height int) *Map {
	m := &Map{
		Width:         width,
		Height:        height,
		Doors:         make([][2]int, 0),
		CheckoutZones: make([][2]int, 0),
		Walls:         make([][2]int, 0),
		ShelfData:     make(map[[2]int]Shelf),
		ShelfChars:    make(map[[2]int]string),
	}

	m.buildFreeCells()
	return m
}

func (m *Map) buildFreeCells() {
	total := m.Width * m.Height
	m.freeCells = make([][2]int, 0, total)

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if !m.IsWall(x, y) &&
				!m.IsCheckout(x, y) &&
				!m.IsShelf(x, y) &&
				!m.IsDoor(x, y) {
				m.freeCells = append(m.freeCells, [2]int{x, y})
			}
		}
	}
}

func (m *Map) AddDoor(x, y int) {
	m.Doors = append(m.Doors, [2]int{x, y})
	m.buildFreeCells()
}

func (m *Map) AddCheckoutZone(x, y int) {
	m.CheckoutZones = append(m.CheckoutZones, [2]int{x, y})
	m.buildFreeCells()
}

func (m *Map) AddWall(x, y int) {
	m.Walls = append(m.Walls, [2]int{x, y})
	m.buildFreeCells()
}

func containsCoordinate(coordinates [][2]int, x, y int) bool {
	for _, coord := range coordinates {
		if coord[0] == x && coord[1] == y {
			return true
		}
	}
	return false
}

func (m *Map) IsWall(x, y int) bool {
	return containsCoordinate(m.Walls, x, y)
}

func (m *Map) IsCheckout(x, y int) bool {
	return containsCoordinate(m.CheckoutZones, x, y)
}

func (m *Map) IsShelf(x, y int) bool {
	_, exists := m.ShelfData[[2]int{x, y}]
	return exists
}

func (m *Map) IsDoor(x, y int) bool {
	return containsCoordinate(m.Doors, x, y)
}

func (m *Map) GetCollisables() [][2]int {
	total := len(m.CheckoutZones) + len(m.Doors) + len(m.ShelfData) + len(m.Walls)
	collisables := make([][2]int, 0, total)

	collisables = append(collisables, m.CheckoutZones...)
	collisables = append(collisables, m.Doors...)
	collisables = append(collisables, m.Walls...)

	for pos := range m.ShelfData {
		collisables = append(collisables, pos)
	}

	return collisables
}

func (m *Map) GetElementType(x, y int) ElementType {
	if m.IsWall(x, y) {
		return WALL
	}
	if m.IsCheckout(x, y) {
		return CHECKOUT
	}
	if m.IsShelf(x, y) {
		return SHELF
	}
	if m.IsDoor(x, y) {
		return DOOR
	}
	return VOID
}

func (m *Map) GetShelfData(x, y int) (Shelf, bool) {
	shelf, exists := m.ShelfData[[2]int{x, y}]
	return shelf, exists
}

func (m *Map) SetShelfData(x, y int, shelf Shelf) {
	m.ShelfData[[2]int{x, y}] = shelf
	m.buildFreeCells()
}

func (m *Map) GetAllShelvesData() map[[2]int]Shelf {
	return m.ShelfData
}

func (m *Map) GetShelfCharacter(x, y int) string {
	if char, exists := m.ShelfChars[[2]int{x, y}]; exists {
		return char
	}
	return ""
}

func (m *Map) GetRandomFreeCoordinate() (int, int, bool) {
	if len(m.freeCells) == 0 {
		return 0, 0, false
	}

	idx := rand.Intn(len(m.freeCells))
	coord := m.freeCells[idx]
	return coord[0], coord[1], true
}

func (m *Map) IsWalkable(x, y int) bool {
	return !m.IsWall(x, y) &&
		!m.IsCheckout(x, y) &&
		!m.IsShelf(x, y) &&
		!m.IsDoor(x, y)
}
