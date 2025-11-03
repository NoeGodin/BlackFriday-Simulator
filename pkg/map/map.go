package Map

func NewShelf(items []Item) *Shelf {
	return &Shelf{
		Element: Element{elementType: SHELF},
		Items:   items,
	}
}

func NewMap(width int, height int) *Map {
	m := Map{
		Width:  width,
		Height: height,
		Grid:   make([][]MapElement, height),
	}
	for i := range m.Grid {
		m.Grid[i] = make([]MapElement, width)
		for y := range m.Grid[i] {
			void := Element{elementType: VOID}
			m.Grid[i][y] = &void
		}
	}
	return &m
}
