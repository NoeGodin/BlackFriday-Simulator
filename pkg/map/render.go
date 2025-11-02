package Map

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// TODO : on pourrait faire plus petit, genre 10 car j'ai peur que ça fasse paté si la map devient grande, on verra
// Si on change il faudra aussi changer les assets car sinon ils vont être déformée vu que reduit
const CellSize = 32

// Draw image at the right position with the right scale
func drawImageAt(screen *ebiten.Image, img *ebiten.Image, x, y float32) {
	if img == nil {
		return
	}
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(float64(CellSize)/float64(img.Bounds().Dx()),
		float64(CellSize)/float64(img.Bounds().Dy()))
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, options)
}

func (g *Game) Update() error {
	//TODO : should i put that here ? i don't really know...
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.CameraY -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.CameraY += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.CameraX -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.CameraX += 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.DrawMap(screen)
	g.DrawHUD(screen)
}

func (g *Game) DrawMap(screen *ebiten.Image) {
	margin := 20
	offsetX := margin
	offsetY := margin

	for i := 0; i <= g.Map.Width; i++ {
		x := float32(i*CellSize - g.CameraX + offsetX)
		vector.StrokeLine(
			screen,
			x,
			float32(offsetY),
			x,
			float32(g.Map.Height*CellSize+offsetY),
			1,
			color.Gray{Y: 240},
			false,
		)
	}

	for i := 0; i <= g.Map.Height; i++ {
		y := float32(i*CellSize - g.CameraY + offsetY)
		vector.StrokeLine(
			screen,
			float32(offsetX),
			y,
			float32(g.Map.Width*CellSize+offsetX),
			y,
			1,
			color.Gray{Y: 240},
			false,
		)
	}

	// DRAW THE GROUND FIRST
	for y := 0; y < g.Map.Height; y++ {
		for x := 0; x < g.Map.Width; x++ {
			drawX := float32(x*CellSize - g.CameraX + offsetX)
			drawY := float32(y*CellSize - g.CameraY + offsetY)

			drawImageAt(screen, groundImg, drawX, drawY)
		}
	}

	// DRAW EVERY OTHER ELEMENTS, SAME LOGIC IN A WAY BUT GOOD PRACTICE TO KEEP BOTH
	for y := 0; y < g.Map.Height; y++ {
		for x := 0; x < g.Map.Width; x++ {
			element := g.Map.Grid[y][x]
			if element == nil || element.Type() == VOID {
				continue
			}

			drawX := float32(x*CellSize - g.CameraX + offsetX)
			drawY := float32(y*CellSize - g.CameraY + offsetY)

			//If image not exist, will not render !
			switch element.Type() {
			case WALL:
				drawImageAt(screen, wallImg, drawX, drawY)
			case ITEM:
				drawImageAt(screen, itemImg, drawX, drawY)
			case DOOR:
				drawImageAt(screen, doorImg, drawX, drawY)
			case CHECKOUT:
				drawImageAt(screen, checkoutImg, drawX, drawY)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// adapt the window size
	mapWidth := g.Map.Width * CellSize
	mapHeight := g.Map.Height * CellSize

	// better with this :)
	margin := 20
	return mapWidth + margin*2, mapHeight + margin*2
}

func NewGame(screenWidth, screenHeight int) *Game {
	mapData, _ := LoadMapFromFile("maps/blackfriday_store")
	//TODO : not catching any errors, i don't know what to do in case of error
	return &Game{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		CameraX:      0,
		CameraY:      0,
		Map:          *mapData,
	}
}
