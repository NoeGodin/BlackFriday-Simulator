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

	// Parsing the map
	for y, line := range mapLines {
		for x, c := range line {
			if x >= width || y >= height {
				continue
			}

			switch c {
			case 'w':
				m.AddWall(x, y)
			case 's':
				m.AddProductZone(x, y)
			case 'd':
				m.AddDoor(x, y)
			case 'c':
				m.AddCheckoutZone(x, y)
			}
		}
	}

	return m, nil
}

func (m *Map) LoadStockData(stockData StockData) {
	//HOW Stocks are assigned to shelves
	//LEFT TO RIGHT TOP TO BOTTOM
	for i, productZone := range m.ProductZones {
		if i < len(stockData.Stocks) {
			position := [2]int{productZone[0], productZone[1]}
			m.ProductData[position] = stockData.Stocks[i]
		}
	}
}
