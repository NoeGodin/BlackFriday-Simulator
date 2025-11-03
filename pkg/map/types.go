package Map

type ElementType string

type MapElement interface {
	Type() ElementType
}
type Element struct {
	elementType ElementType
}

type Item struct {
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Reduction      float64 `json:"reduction"`
	Attractiveness float64 `json:"attractiveness"`
	Quantity       int     `json:"quantity"`
}

type Shelf struct {
	Element
	Items []Item
}

type Map struct {
	Width  int
	Height int
	Grid   [][]MapElement
}

func (element *Element) Type() ElementType {
	return element.elementType
}

type StockData struct {
	Stocks [][]Item `json:"stocks"`
}

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	Map                       Map
}
