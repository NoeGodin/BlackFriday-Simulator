package Map

import (
	"os"
	"strconv"
	"strings"
)

func LoadMapFromFile(filename string) (*Map, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadMapFromString(string(content))
}

// w = WALL
// i = ITEM
// c = CHECKOUT
// d = DOOR
// " "= VOID
func LoadMapFromString(content string) (*Map, error) {
	lines := strings.Split(content, "\n")

	mapLines := []string{} // contains every map lines
	itemData := []string{} // contains item data (price, attractiveness)
	inItemSection := false // when reading item

	for _, line := range lines {
		trimmed := strings.TrimSpace(line) //TODO: maybe change name of var ? idk if its clear

		if trimmed == "---" {
			inItemSection = true
			continue
		}

		if inItemSection && len(trimmed) > 0 {
			itemData = append(itemData, trimmed)
		} else if len(line) > 0 {
			mapLines = append(mapLines, line)
		}
	}

	height := len(mapLines)
	width := 0
	if height > 0 {
		width = len(mapLines[0]) //fixing width based on first line
	}

	m := NewMap(width, height)
	itemPositions := [][]int{}

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
			case 'i':
				element = &Element{elementType: ITEM}
				itemPositions = append(itemPositions, []int{x, y})
			case 'd':
				element = &Element{elementType: DOOR}
			case 'c':
				element = &Element{elementType: CHECKOUT}
			case ' ':
				element = &Element{elementType: VOID}
			default:
				element = &Element{elementType: VOID}
			}

			if element != nil {
				m.Grid[y][x] = element
			}
		}
	}

	for i, pos := range itemPositions {
		x, y := pos[0], pos[1]

		// Default value
		price := 10.0
		attractiveness := 0.5

		//custom value
		if i < len(itemData) {
			parts := strings.Split(itemData[i], ",")
			if len(parts) >= 2 {
				if p, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64); err == nil {
					price = p
				}
				if a, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					attractiveness = a
				}
			}
		}

		if x < width && y < height {
			item := NewItem(price, 0.0, attractiveness)
			m.Grid[y][x] = item
		}
	}

	return m, nil
}
