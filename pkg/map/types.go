package Map

import (
	"AI30_-_BlackFriday/pkg/constants"
	"fmt"
)

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
	Width           int
	Height          int
	Doors           [][2]float64
	CheckoutZones   [][2]float64
	Walls           [][2]float64
	ShelfData       map[[2]float64]Shelf
	ShelfChars      map[[2]float64]string
	freeCells       [][2]float64
	Items           []Item
	AvailableProfit float64
}

func (element *Element) Type() ElementType {
	return element.elementType
}

type StockData struct {
	Stocks map[string]Shelf `json:"stocks"`
}

func (i *Item) String() string {
	return fmt.Sprintf("Name: %s, Price: %.2f, Reduction: %.2f, Attractiveness: %.2f, Quantity: %d\n", i.Name, i.Price, i.Reduction, i.Attractiveness, i.Quantity)
}
