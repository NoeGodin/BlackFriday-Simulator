package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	Hud "AI30_-_BlackFriday/pkg/hud"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	textureFont        font.Face
	textureInitialized bool
)

// Global vars
var (
	agtImg        *ebiten.Image
	wallCeiling   *ebiten.Image
	wallImg       *ebiten.Image
	groundImg     *ebiten.Image
	doorImg       *ebiten.Image
	itemImg       *ebiten.Image
	checkoutImg   *ebiten.Image
	targetImg     *ebiten.Image
	WalkFrameImgs [constants.DIRECTIONS][constants.FRAME_COUNT]*ebiten.Image
)

func initTexture() {
	if textureInitialized {
		return
	}

	var err error
	agtImg, _, err = ebitenutil.NewImageFromFile("assets/agt.png")
	if err != nil {
		log.Printf("Warning: Could not load agt.png: %v", err)
	}

	wallImg, _, err = ebitenutil.NewImageFromFile("assets/wall_front.png")
	if err != nil {
		log.Printf("Warning: Could not load wall_front.png: %v", err)
	}
	
	wallCeiling, _, err = ebitenutil.NewImageFromFile("assets/wall_ceiling.png")
	if err != nil {
		log.Printf("Warning: Could not load wall_ceiling.png: %v", err)
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

	targetImg, _, err = ebitenutil.NewImageFromFile("assets/target.png")
	if err != nil {
		log.Printf("Warning: Could not load target.png: %v", err)
	}

	checkoutImg, _, err = ebitenutil.NewImageFromFile("assets/checkout.png")
	if err != nil {
		log.Printf("Warning: Could not load checkout.png: %v", err)
	}

	fontBytes, err := os.ReadFile("assets/fonts/Monaco.ttf")
	if err != nil {
		panic(err)
	}

	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	Hud.FONT, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    16,
		DPI:     96,
		Hinting: font.HintingFull,
	})

	if err != nil {
		panic(err)
	}

	textureInitialized = true
}

func initAnimation() {
	walk, _, err := ebitenutil.NewImageFromFile("assets/walk.png")
	if err != nil {
		log.Printf("Warning: Could not load walk.png: %v", err)
	}
	for dir := range constants.DIRECTIONS {
		for f := range constants.FRAME_COUNT {
			sx := f * constants.CELL_SIZE
			sy := dir * constants.CELL_SIZE
			sub := walk.SubImage(image.Rect(sx, sy, sx+constants.CELL_SIZE, sy+constants.CELL_SIZE)).(*ebiten.Image)
			WalkFrameImgs[dir][f] = sub
		}
	}
}

func (g *Game) DrawTexture(screen *ebiten.Image) {
	initTexture()
	initAnimation()
}