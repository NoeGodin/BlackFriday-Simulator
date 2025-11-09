package Map

func NewMap(width int, height int) *Map {
	return &Map{
		Width:         width,
		Height:        height,
		Doors:         make([][2]int, 0),
		CheckoutZones: make([][2]int, 0),
		ProductZones:  make([][2]int, 0),
		Walls:         make([][2]int, 0),
		ProductData:   make(map[[2]int][]Item),
	}
}

func (m *Map) AddDoor(x, y int) {
	m.Doors = append(m.Doors, [2]int{x, y})
}

func (m *Map) AddCheckoutZone(x, y int) {
	m.CheckoutZones = append(m.CheckoutZones, [2]int{x, y})
}

func (m *Map) AddProductZone(x, y int) {
	m.ProductZones = append(m.ProductZones, [2]int{x, y})
}

func (m *Map) AddWall(x, y int) {
	m.Walls = append(m.Walls, [2]int{x, y})
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

func (m *Map) IsProductZone(x, y int) bool {
	return containsCoordinate(m.ProductZones, x, y)
}

func (m *Map) IsDoor(x, y int) bool {
	return containsCoordinate(m.Doors, x, y)
}

func (m *Map) GetElementType(x, y int) ElementType {
	if m.IsWall(x, y) {
		return WALL
	}
	if m.IsCheckout(x, y) {
		return CHECKOUT
	}
	if m.IsProductZone(x, y) {
		return SHELF
	}
	if m.IsDoor(x, y) {
		return DOOR
	}
	return VOID
}

func (m *Map) GetProductData(x, y int) ([]Item, bool) {
	items, exists := m.ProductData[[2]int{x, y}]
	return items, exists
}

func (m *Map) SetProductData(x, y int, items []Item) {
	m.ProductData[[2]int{x, y}] = items
}

func (m *Map) GetAllProductsInZone() map[[2]int][]Item {
	return m.ProductData
}
