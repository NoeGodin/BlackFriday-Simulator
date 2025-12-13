package main

import (
	mapgenerator "AI30_-_BlackFriday/pkg/map-generator"
	"AI30_-_BlackFriday/pkg/utils"
	"fmt"
	"os"
)

func main() {
	nbMaps := 5
	width := 30
	height := 20
	nbDoors := 5
	nbShelves := 25
	nbCashiers := 3
	nbWalls := 140

	mapLayout := mapgenerator.NewMapLayout(width, height)
	
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


	for i := 0; i < nbMaps; i++ {
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
		} else {
			fmt.Println("Not valid map, regenerating...")
			i--
			continue
		}
		
		err = os.WriteFile("maps/generated_maps/map" + fmt.Sprint(i+1) + ".txt", []byte(mapStr), 0644)
		if err != nil {
			fmt.Printf("unable to write file: %v", err)
		}
	}
}
