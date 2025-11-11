package main

import (
	"log"
	"os"

	Graphics "AI30_-_BlackFriday/pkg/graphics"
	Map "AI30_-_BlackFriday/pkg/map"
	Simulation "AI30_-_BlackFriday/pkg/simulation"

	"github.com/hajimehoshi/ebiten/v2"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700

func main() {
	//NOTE : I guess its for my setup (working on mac) maybe not needed for others idk :/
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Black Friday Simulator")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	mapData, err := Map.LoadMapFromFile("maps/store/layout.txt")
	if err != nil {
		panic("Error loading map: " + err.Error())
	}
	simu := Simulation.NewSimulation(0, mapData)

	// Appel de la fonction
	simu.AddClient("agent1")
	game := Graphics.NewGame(SCREEN_WIDTH, SCREEN_HEIGHT, simu)
	game.Simulation.Run()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
