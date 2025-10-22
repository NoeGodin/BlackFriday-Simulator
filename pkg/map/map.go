package Map

func NewItem(price float64, reduction float64, attractiveness float64, typeItem ElementType) *Item {
	return &Item{
		Element:        Element{elementType: typeItem},
		Price:          price,
		Reduction:      reduction,
		Attractiveness: attractiveness,
	}
}

func NewMap(width int, height int) *Map {
	m := Map{
		Width:  width,
		Height: height,
		Grid:   make([][]MapElement, height),
	}
	for i := range m.Grid {
		m.Grid[i] = make([]MapElement, 5)
		for y := range m.Grid[i] {
			void := Element{elementType: VOID}
			m.Grid[i][y] = &void
		}
	}
	return &m
}