package main

import (
	"AI30_-_BlackFriday/pkg/constants"
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/shopping"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type StockData struct {
	Stocks map[string]struct {
		Items      []Map.Item `json:"items"`
		Categories []string   `json:"categories"`
	} `json:"stocks"`
}

func main() {
	stockFile, err := os.Open("maps/store/stocks.json")
	if err != nil {
		log.Fatalf("Error opening file stocks.json: %v", err)
	}
	defer stockFile.Close()

	var stockData StockData
	if err := json.NewDecoder(stockFile).Decode(&stockData); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Extract all items
	var allItems []Map.Item
	for shelfID, shelf := range stockData.Stocks {
		fmt.Printf("Shelf %s: %d items\n", shelfID, len(shelf.Items))
		for _, item := range shelf.Items {
			allItems = append(allItems, item)
		}
	}

	const (
		seed     = 12345
		filePath = "maps/store/shopping_lists.json"
	)

	numAgents := constants.NUMBER_OF_CLIENTS

	fmt.Printf("Generate %d cours list with seed %d from %d items stocks...\n",
		numAgents, seed, len(allItems))

	err = shopping.GenerateShoppingListsFile(filePath, numAgents, allItems, seed)
	if err != nil {
		log.Fatalf("Error generating file: %v", err)
	}

	fmt.Printf("File generated successfully: %s\n", filePath)
}
