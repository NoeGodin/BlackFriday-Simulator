package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	Hud "AI30_-_BlackFriday/pkg/hud"
	Simulation "AI30_-_BlackFriday/pkg/simulation"
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

type WalkAnimation [constants.DIRECTIONS][constants.FRAME_COUNT]*ebiten.Image

// Global vars
var (
	wallCeiling        *ebiten.Image
	wallImg            *ebiten.Image
	groundImg          *ebiten.Image
	doorImg            *ebiten.Image
	itemImg            *ebiten.Image
	itemAlmostFullImg  *ebiten.Image
	itemHalfEmptyImg   *ebiten.Image
	itemAlmostEmptyImg *ebiten.Image
	itemEmptyImg       *ebiten.Image
	checkoutImg        *ebiten.Image
	targetImg          *ebiten.Image
	BaseFrameImgs      WalkAnimation
	WalkFrameImgs      = make(map[Simulation.AgenType]WalkAnimation)
)

func initTexture() {
	if textureInitialized {
		return
	}

	var err error

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

	itemEmptyImg, _, err = ebitenutil.NewImageFromFile("assets/item_empty.png")
	if err != nil {
		log.Printf("Warning: Could not load item_empty.png: %v", err)
	}

	targetImg, _, err = ebitenutil.NewImageFromFile("assets/target.png")
	if err != nil {
		log.Printf("Warning: Could not load target.png: %v", err)
	}

	checkoutImg, _, err = ebitenutil.NewImageFromFile("assets/checkout.png")
	if err != nil {
		log.Printf("Warning: Could not load checkout.png: %v", err)
	}

	BaseFrameImgs, err = initAnimation("assets/walk.png")
	if err != nil {
		log.Fatalf("Error: Could not load walk.png: %v", err)
	}
	WalkFrameImgs[Simulation.CLIENT] = BaseFrameImgs

	GuardAnimation, err := initAnimation("assets/walk_guard.png")
	if err != nil {
		log.Fatalf("Error: Could not load walk_guard.png: %v", err)
	}
	WalkFrameImgs[Simulation.GUARD] = GuardAnimation

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

func initAnimation(path string) (WalkAnimation, error) {
	walk, _, err := ebitenutil.NewImageFromFile(path)
	var anim WalkAnimation

	if err != nil {
		return anim, err
	}
	for dir := range constants.DIRECTIONS {
		for f := range constants.FRAME_COUNT {
			sx := f * constants.CELL_SIZE
			sy := dir * constants.CELL_SIZE
			sub := walk.SubImage(image.Rect(sx, sy, sx+constants.CELL_SIZE, sy+constants.CELL_SIZE)).(*ebiten.Image)
			anim[dir][f] = sub
		}
	}
	return anim, nil
}

func getWalkAnimation(agtType Simulation.AgenType) WalkAnimation {
	walkAnimation, ok := WalkFrameImgs[agtType]
	if !ok {
		return BaseFrameImgs
	}
	return walkAnimation
}
func (g *Game) DrawTexture(screen *ebiten.Image) {
	initTexture()
}
