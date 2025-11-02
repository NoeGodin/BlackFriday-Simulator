package main

import (
	Map "AI30_-_BlackFriday/pkg/map"
	"fmt"
)

func main() {
	item := Map.NewItem(19.99,0.1,0.8,Map.ITEM)
	m := Map.NewMap(5,5)
    m.Grid[0][0] = item
	fmt.Printf(string(m.Grid[0][0].Type()))
}
