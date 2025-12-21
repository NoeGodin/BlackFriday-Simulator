package main

import (
	"log"
	"os"
	"runtime"

	Graphics "AI30_-_BlackFriday/pkg/graphics"
	"AI30_-_BlackFriday/pkg/logger"

	"github.com/hajimehoshi/ebiten/v2"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700

func main() {
	logger.InitLogger()

	//NOTE : I guess its for my setup (working on mac) maybe not needed for others idk :/
	if runtime.GOOS == "darwin" {
		os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	}

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Black Friday Simulator")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	logger.Info("Loading map...")

	game := Graphics.NewGame(SCREEN_WIDTH, SCREEN_HEIGHT)

	// launch simulation, generate repport at the end
	defer func() {
		if game.Simulation == nil {
			return
		}
		logger.Info("Exporting sales data...")
		if err := game.Simulation.Env.ExportSalesData(); err != nil {
			logger.Errorf("Error exporting sales data: %v", err)
		} else {
			logger.Info("Sales data exported successfully")
		}
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
