package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"AI30_-_BlackFriday/pkg/constants"
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
	if runtime.GOOS == "darwin" {
		os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	}

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Black Friday Simulator")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	logger.Info("Loading map...")
	mapPath := "maps/store/large_layout.txt"
	mapData, err := Map.LoadMapFromFile(mapPath)
	if err != nil {
		logger.Errorf("Error loading map: %s", err.Error())
		panic("Error loading map: " + err.Error())
	}

	mapName := extractMapName(mapPath)
	shoppingListsPath := "maps/store/shopping_lists.json"
	logger.Info("Creating simulation...")
	simu := Simulation.NewSimulation(0, constants.TIC_DURATION, mapData, constants.DELTA_TIME, constants.AGENT_SEARCH_RADIUS, mapName, shoppingListsPath)

	logger.Info("Adding agents...")
	logger.Infof("Creating %d agents from config", constants.NUMBER_OF_AGENTS)
	for i := 1; i <= constants.NUMBER_OF_AGENTS; i++ {
		agentID := fmt.Sprintf("agent%d", i)
		simu.AddClient(agentID)
	}

	logger.Info("Starting game...")
	game := Graphics.NewGame(SCREEN_WIDTH, SCREEN_HEIGHT, simu)
	game.Simulation.Run()

	// launch simulation, generate repport at the end
	defer func() {
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

func extractMapName(mapPath string) string {
	filename := filepath.Base(mapPath)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	dir := filepath.Base(filepath.Dir(mapPath))

	if dir != "." && dir != "" {
		return fmt.Sprintf("%s_%s", dir, name)
	}

	return name
}
