package main

import (
	"log"
	"os"

	Graphics "AI30_-_BlackFriday/pkg/graphics"

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

	game := Graphics.NewGame(SCREEN_WIDTH, SCREEN_HEIGHT)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
