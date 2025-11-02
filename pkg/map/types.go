package Map

type ElementType string

type MapElement interface {
	Type() ElementType
}
type Element struct {
	elementType ElementType
}

type Item struct {
	Element
	Price          float64
	Reduction      float64
	Attractiveness float64
}

type Map struct {
	Width  int
	Height int
	Grid   [][]MapElement
}

func (element *Element) Type() ElementType {
	return element.elementType
}

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	Map                       Map
}
