package shopping

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"

	"AI30_-_BlackFriday/pkg/constants"
	Map "AI30_-_BlackFriday/pkg/map"
)

type PredefShoppingList struct {
	AgentID int        `json:"agentId"`
	Items   []Map.Item `json:"items"`
}

type ShoppingListLoader struct {
	lists []PredefShoppingList
}

func NewShoppingListLoader(filePath string) (*ShoppingListLoader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open shopping lists file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read shopping lists file: %w", err)
	}

	var lists []PredefShoppingList
	if err := json.Unmarshal(data, &lists); err != nil {
		return nil, fmt.Errorf("failed to parse shopping lists JSON: %w", err)
	}

	return &ShoppingListLoader{lists: lists}, nil
}

func (sl *ShoppingListLoader) GetShoppingList(agentIndex int) []Map.Item {
	if agentIndex < 0 || agentIndex >= len(sl.lists) {
		// if out of bound
		agentIndex = agentIndex % len(sl.lists)
	}
	return sl.lists[agentIndex].Items
}

func GenerateShoppingListsFile(filePath string, numAgents int, items []Map.Item, seed int64) error {
	rng := rand.New(rand.NewSource(seed))

	totalAttractiveness := 0.0
	for _, item := range items {
		totalAttractiveness += item.Attractiveness
	}

	var shoppingLists []PredefShoppingList

	for i := 0; i < numAgents; i++ {
		var agentShoppingList []Map.Item

		// Generate a shopping list by considering attractiveness
		numItemsInList := rng.Intn(constants.AGENT_MAX_SHOPPING_LIST) + 1

		for j := 0; j < numItemsInList; j++ {
			wantedItem := rng.Float64() * totalAttractiveness
			cumulative := 0.0

			for _, item := range items {
				cumulative += item.Attractiveness
				if wantedItem <= cumulative {
					// Crandom quantity between 1 and min(stock_total, MAX_QUANTITY_PER_ITEM)
					itemCopy := item
					maxQuantity := item.Quantity
					if constants.MAX_QUANTITY_PER_ITEM < maxQuantity {
						maxQuantity = constants.MAX_QUANTITY_PER_ITEM
					}
					if maxQuantity > 1 {
						itemCopy.Quantity = rng.Intn(maxQuantity) + 1
					} else {
						itemCopy.Quantity = 1
					}
					agentShoppingList = append(agentShoppingList, itemCopy)
					break
				}
			}
		}

		shoppingLists = append(shoppingLists, PredefShoppingList{
			AgentID: i,
			Items:   agentShoppingList,
		})
	}

	data, err := json.MarshalIndent(shoppingLists, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal shopping lists: %w", err)
	}

	return os.WriteFile(filePath, data, 0644)
}
