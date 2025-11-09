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
	Width         int
	Height        int
	Doors         [][2]int
	CheckoutZones [][2]int
	ProductZones  [][2]int
	Walls         [][2]int
	ProductData   map[[2]int][]Item
}

func (element *Element) Type() ElementType {
	return element.elementType
}

type StockData struct {
	Stocks [][]Item `json:"stocks"`
}
