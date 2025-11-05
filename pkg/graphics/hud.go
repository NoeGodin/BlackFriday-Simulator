package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var (
	hudFont        font.Face
	hudMsg		   string
	hudInitialized bool
)

// Global vars
var (
	agtImg        *ebiten.Image
	wallImg       *ebiten.Image
	groundImg     *ebiten.Image
	doorImg       *ebiten.Image
	itemImg       *ebiten.Image
	checkoutImg   *ebiten.Image
	WalkFrameImgs [constants.DIRECTIONS][constants.FRAME_COUNT]*ebiten.Image
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

	hudFont = basicfont.Face7x13

	hudInitialized = true
}

func initAnimation() {
	walk, _, err := ebitenutil.NewImageFromFile("assets/walk.png")
	if err != nil {
		log.Printf("Warning: Could not load agt.png: %v", err)
	}
	for dir := 0; dir < constants.DIRECTIONS; dir++ {
		for f := 0; f < constants.FRAME_COUNT; f++ {
			sx := f * constants.CELL_SIZE
			sy := dir * constants.CELL_SIZE
			sub := walk.SubImage(image.Rect(sx, sy, sx+constants.CELL_SIZE, sy+constants.CELL_SIZE)).(*ebiten.Image)
			WalkFrameImgs[dir][f] = sub
		}
	}
}
func (g *Game) DrawHUD(screen *ebiten.Image) {
	initHUD()
	initAnimation()

	hudWidth := 200
    hudHeight := 60
    hudX := 470
    hudY := 10

    hudBg := ebiten.NewImage(hudWidth, hudHeight)
    hudBg.Fill(color.RGBA{0, 0, 0, 180})

    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(float64(hudX), float64(hudY))
    screen.DrawImage(hudBg, op)

    text.Draw(screen, hudMsg, hudFont, hudX+10, hudY+25, color.White)

}

func UpdateHUD(posX, posY int) {
	hudMsg = fmt.Sprintln("DEBUG")
    hudMsg += fmt.Sprintf("Position: (%d, %d)", posX, posY)
}