package Map

import (
	"AI30_-_BlackFriday/pkg/constants"
	"math/rand"
)

func NewMap(width int, height int) *Map {
	m := &Map{
		Width:         width,
		Height:        height,
		Doors:         make([][2]float64, 0),
		CheckoutZones: make([][2]float64, 0),
		Walls:         make([][2]float64, 0),
		ShelfData:     make(map[[2]float64]Shelf),
		ShelfChars:    make(map[[2]float64]string),
	}

	m.buildFreeCells()
	return m
}

func (m *Map) buildFreeCells() {
	total := m.Width * m.Height
	m.freeCells = make([][2]float64, 0, total)

	for x := 0.0; x < float64(m.Width); x++ {
		for y := 0.0; y < float64(m.Height); y++ {
			if !m.IsWall(x, y) &&
				!m.IsCheckout(x, y) &&
				!m.IsShelf(x, y) &&
				!m.IsDoor(x, y) {
				m.freeCells = append(m.freeCells, [2]float64{x, y})
			}
		}
	}
}

func (m *Map) AddDoor(x, y float64) {
	m.Doors = append(m.Doors, [2]float64{x, y})
	m.buildFreeCells()
}

func (m *Map) AddCheckoutZone(x, y float64) {
	m.CheckoutZones = append(m.CheckoutZones, [2]float64{x, y})
	m.buildFreeCells()
}

func (m *Map) AddWall(x, y float64) {
	m.Walls = append(m.Walls, [2]float64{x, y})
	m.buildFreeCells()
}

func containsCoordinate(coordinates [][2]float64, x, y float64) bool {
	for _, coord := range coordinates {
		if coord[0] == x && coord[1] == y {
			return true
		}
	}
	return false
}

func (m *Map) IsWall(x, y float64) bool {
	return containsCoordinate(m.Walls, x, y)
}

func (m *Map) IsCheckout(x, y float64) bool {
	return containsCoordinate(m.CheckoutZones, x, y)
}

func (m *Map) IsShelf(x, y float64) bool {
	_, exists := m.ShelfData[[2]float64{x, y}]
	return exists
}

func (m *Map) IsDoor(x, y float64) bool {
	return containsCoordinate(m.Doors, x, y)
}

func (m *Map) GetCollisables() [][2]float64 {
	total := len(m.CheckoutZones) + len(m.Doors) + len(m.ShelfData) + len(m.Walls)
	collisables := make([][2]float64, 0, total)

	collisables = append(collisables, m.CheckoutZones...)
	collisables = append(collisables, m.Doors...)
	collisables = append(collisables, m.Walls...)

	for pos := range m.ShelfData {
		collisables = append(collisables, pos)
	}

	return collisables
}

func (m *Map) GetElementType(x, y float64) constants.ElementType {
	if m.IsWall(x, y) {
		return constants.WALL
	}
	if m.IsCheckout(x, y) {
		return constants.CHECKOUT
	}
	if m.IsShelf(x, y) {
		return constants.SHELF
	}
	if m.IsDoor(x, y) {
		return constants.DOOR
	}
	return constants.VOID
}

func (m *Map) GetShelfData(x, y float64) (Shelf, bool) {
	shelf, exists := m.ShelfData[[2]float64{x, y}]
	return shelf, exists
}

func (m *Map) SetShelfData(x, y float64, shelf Shelf) {
	m.ShelfData[[2]float64{x, y}] = shelf
	m.buildFreeCells()
}

func (m *Map) GetAllShelvesData() map[[2]float64]Shelf {
	return m.ShelfData
}

func (m *Map) GetShelfCharacter(x, y float64) string {
	if char, exists := m.ShelfChars[[2]float64{x, y}]; exists {
		return char
	}
	return ""
}

func (m *Map) GetRandomFreeCoordinate() (float64, float64, bool) {
	if len(m.freeCells) == 0 {
		return 0, 0, false
	}

	idx := rand.Intn(len(m.freeCells))
	coord := m.freeCells[idx]
	return coord[0], coord[1], true
}

func (m *Map) IsWalkable(x, y float64) bool {
	return !m.IsWall(x, y) &&
		!m.IsCheckout(x, y) &&
		!m.IsShelf(x, y) &&
		!m.IsDoor(x, y)
}

func (m *Map) IsValidAndWalkable(x, y float64) bool {
	return x >= 0 && x < float64(m.Width) && y >= 0 && y < float64(m.Height) && m.IsWalkable(x, y)
}
