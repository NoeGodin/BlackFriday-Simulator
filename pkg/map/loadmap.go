package Map

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func LoadMapFromFile(filename string) (*Map, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	m, err := LoadMapFromString(string(content))
	if err != nil {
		return nil, err
	}

	// Same dir
	stockFile := filepath.Join(filepath.Dir(filename), "stocks.json")
	stockData, err := os.ReadFile(stockFile)
	if err != nil {
		return nil, err
	}

	var stocks StockData
	err = json.Unmarshal(stockData, &stocks)
	if err != nil {
		return nil, err
	}

	m.LoadStockData(stocks)
	return m, nil
}

// w = WALL
// s = SHELF
// c = CHECKOUT
// d = DOOR
// " "= VOID
func LoadMapFromString(content string) (*Map, error) {
	lines := strings.Split(content, "\n")

	mapLines := []string{} // contains every map lines

	for _, line := range lines {
		if len(line) > 0 {
			mapLines = append(mapLines, line)
		}
	}

	height := len(mapLines)
	width := 0
	if height > 0 {
		width = len(mapLines[0]) //fixing width based on first line
	}

	m := NewMap(width, height)
	shelfPositions := [][]int{}

	// Parsing the map
	for y, line := range mapLines {
		for x, c := range line {
			if x >= width || y >= height {
				continue
			}

			var element *Element
			switch c {
			case 'w':
				element = &Element{elementType: WALL}
			case 's':
				element = &Element{elementType: SHELF}
				shelfPositions = append(shelfPositions, []int{x, y})
			case 'd':
				element = &Element{elementType: DOOR}
			case 'c':
				element = &Element{elementType: CHECKOUT}
			default:
				element = &Element{elementType: VOID}
			}

			if element != nil {
				m.Grid[y][x] = element
			}
		}
	}

	// Will populate after with JSON File
	for _, pos := range shelfPositions {
		x, y := pos[0], pos[1]
		if x < width && y < height {
			// will be filled by LoadStockData
			shelf := NewShelf([]Item{})
			m.Grid[y][x] = shelf
		}
	}

	return m, nil
}

func (m *Map) LoadStockData(stockData StockData) {
	shelfIndex := 0

	//HOW Stocks are assigned to shelves
	//LEFT TO RIGHT TOP TO BOTTOM
	for y := range m.Height {
		for x := range m.Width {
			if element := m.Grid[y][x]; element != nil && element.Type() == SHELF {
				if shelfIndex < len(stockData.Stocks) {
					if shelf, ok := element.(*Shelf); ok {
						shelf.Items = stockData.Stocks[shelfIndex]
					}
					shelfIndex++
				}
			}
		}
	}
}
