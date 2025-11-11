package Graphics

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
)

var (
	hudFont        font.Face
	hudInitialized bool
)

// Global vars
var (
	agtImg      *ebiten.Image
	wallImg     *ebiten.Image
	groundImg   *ebiten.Image
	doorImg     *ebiten.Image
	itemImg     *ebiten.Image
	checkoutImg *ebiten.Image
)

func initHUD() {
	if hudInitialized {
		return
	}

	var err error
	agtImg, _, err = ebitenutil.NewImageFromFile("assets/agt.png")
	if err != nil {
		log.Printf("Warning: Could not load agt.png: %v", err)
	}
	wallImg, _, err = ebitenutil.NewImageFromFile("assets/wall.png")
	if err != nil {
		log.Printf("Warning: Could not load wall.png: %v", err)
	}

	groundImg, _, err = ebitenutil.NewImageFromFile("assets/ground.png")
	if err != nil {
		log.Printf("Warning: Could not load ground.png: %v", err)
	}

	doorImg, _, err = ebitenutil.NewImageFromFile("assets/door.png")
	if err != nil {
		log.Printf("Warning: Could not load door.png: %v", err)
	}

	itemImg, _, err = ebitenutil.NewImageFromFile("assets/item.png")
	if err != nil {
		log.Printf("Warning: Could not load item.png: %v", err)
	}

	checkoutImg, _, err = ebitenutil.NewImageFromFile("assets/checkout.png")
	if err != nil {
		log.Printf("Warning: Could not load checkout.png: %v", err)
	}

	hudInitialized = true
}

func (g *Game) DrawHUD(screen *ebiten.Image) {
	initHUD()
}
