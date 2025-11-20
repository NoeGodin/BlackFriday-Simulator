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

	m.LoadStockData(stocks, string(content))
	m.buildFreeCells()
	return m, nil
}

// W = WALL
// a-z, 0-9 = SHELF (mapped to specific products)
// C = CHECKOUT
// D = DOOR
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

	// Parsing the map
	for y, line := range mapLines {
		for x, c := range line {
			if x >= width || y >= height {
				continue
			}

			switch c {
			case 'W':
				m.AddWall(x, y)
			case 'D':
				m.AddDoor(x, y)
			case 'C':
				m.AddCheckoutZone(x, y)
			default:
				// Shelves are handled in LoadStockData, no need to add them here
			}
		}
	}

	return m, nil
}

func (m *Map) LoadStockData(stockData StockData, layoutContent string) {
	lines := strings.Split(layoutContent, "\n")
	mapLines := []string{}
	for _, line := range lines {
		if len(line) > 0 {
			mapLines = append(mapLines, line)
		}
	}

	// Map product to character in the layout
	for y, line := range mapLines {
		for x, c := range line {
			if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
				charKey := string(c)
				position := [2]int{x, y}
				m.ShelfChars[position] = charKey
				// Store shelf data
				if shelf, exists := stockData.Stocks[charKey]; exists {
					m.ShelfData[position] = shelf
				}
			}
		}
	}
}
