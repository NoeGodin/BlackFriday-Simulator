package Map

import "AI30_-_BlackFriday/pkg/constants"

type ElementType = constants.ElementType

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
	Items      []Item   `json:"items"`
	Categories []string `json:"categories"`
}

type Map struct {
	Width         int
	Height        int
	Doors         [][2]int
	CheckoutZones [][2]int
	Walls         [][2]int
	ShelfData     map[[2]int]Shelf
	ShelfChars    map[[2]int]string
	freeCells     [][2]int
}

func (element *Element) Type() ElementType {
	return element.elementType
}

type StockData struct {
	Stocks map[string]Shelf `json:"stocks"`
}
