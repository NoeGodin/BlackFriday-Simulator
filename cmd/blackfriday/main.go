package main

import (
	"fmt"
	"log"
	"os"

	Graphics "AI30_-_BlackFriday/pkg/graphics"
	"AI30_-_BlackFriday/pkg/logger"
	Map "AI30_-_BlackFriday/pkg/map"
	Simulation "AI30_-_BlackFriday/pkg/simulation"

	"github.com/hajimehoshi/ebiten/v2"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700

func main() {
	logger.InitLogger()

	//NOTE : I guess its for my setup (working on mac) maybe not needed for others idk :/
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Black Friday Simulator")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	logger.Info("Loading map...")
	mapData, err := Map.LoadMapFromFile("maps/store/large_layout.txt")
	if err != nil {
		logger.Errorf("Error loading map: %s", err.Error())
		panic("Error loading map: " + err.Error())
	}

	logger.Info("Creating simulation...")
	simu := Simulation.NewSimulation(0, 100.0, mapData)

	logger.Info("Adding agents...")
	for i := 1; i <= 50; i++ {
		agentID := fmt.Sprintf("agent%d", i)
		simu.AddClient(agentID)
	}

	logger.Info("Starting game...")
	game := Graphics.NewGame(SCREEN_WIDTH, SCREEN_HEIGHT, simu)
	game.Simulation.Run()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
