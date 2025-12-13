package main

import (
	mapgenerator "AI30_-_BlackFriday/pkg/map-generator"
	"AI30_-_BlackFriday/pkg/utils"
	"fmt"
	"os"
	"sync"
)

func main() {
	nbMaps := 5
	width := 30
	height := 20
	nbDoors := 5
	nbShelves := 25
	nbCashiers := 1
	// Rough estimation: We can go up to 33% of walls in a map (n * m) size, above there are too much of conflicts
	nbWalls := 170

	mapLayouts := make([]mapgenerator.MapLayout, 0, nbMaps)
	for range nbMaps {
		mapLayouts = append(mapLayouts, mapgenerator.NewMapLayout(width, height))
	}
	var wg sync.WaitGroup
	
	err := os.MkdirAll("maps/generated_maps", os.ModePerm)
	if err != nil {
		fmt.Printf("unable to create directories: %v\n", err)
		return
	}

	err = utils.CopyFile("maps/store/stocks.json", "maps/generated_maps/stocks.json")
	if err != nil {
		fmt.Println("Unable to copy maps/store/stocks.json :", err)
	} else {
		fmt.Println("Copy of maps/store/stocks.json was successfull")
	}

	for i, mapLayout := range mapLayouts {
		wg.Add(1)

		go func(mapLayout mapgenerator.MapLayout, i int) {
			defer wg.Done()

			for {
				mapLayout.CleanMapLayout()
				mapLayout.FillRow(0)
				mapLayout.FillRow(height-1)
				mapLayout.FillColumn(0)
				mapLayout.FillColumn(width-1)
				mapLayout.GenerateDoors(nbDoors)
				mapLayout.GenerateWalls(nbWalls+nbCashiers+nbShelves)
				
				mapLayout.GenerateShelves(nbShelves)
				mapLayout.GenerateCashiers(nbCashiers)
				mapLayout.RemoveAllDoorsWallsSurrounding()
	
				mapStr := mapLayout.ToString()
	
				if mapLayout.IsMapValid() {
					fmt.Println("OK")
					err = os.WriteFile("maps/generated_maps/map" + fmt.Sprint(i+1) + ".txt", []byte(mapStr), 0644)
					if err != nil {
						fmt.Printf("unable to write file: %v", err)
					}
					return
				} else {
					fmt.Println("Not valid map, regenerating...")
				}
			}
		}(mapLayout, i)
	}

	wg.Wait()
}
