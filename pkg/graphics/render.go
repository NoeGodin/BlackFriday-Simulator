package Graphics

import (
	Map "AI30_-_BlackFriday/pkg/map"

	"github.com/hajimehoshi/ebiten/v2"
	//Line drawing debugging dependency
	//"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// TODO : on pourrait faire plus petit, genre 10 car j'ai peur que ça fasse paté si la map devient grande, on verra
// Si on change il faudra aussi changer les assets car sinon ils vont être déformée vu que reduit
const CellSize = 32

func (g *Game) mapToDrawCoords(mapX, mapY, offsetX, offsetY int) (float32, float32) {
	drawX := float32(mapX*CellSize - g.CameraX + offsetX)
	drawY := float32(mapY*CellSize - g.CameraY + offsetY)
	return drawX, drawY
}

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
	g.HandleInput()
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

	// DRAW GRID LINES USEFUL FOR DEBUGGING
	/*for i := range g.Map.Width {
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

	for i := range g.Map.Height {
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
	}*/

	// DRAW THE GROUND FIRST
	for y := range g.Map.Height {
		for x := range g.Map.Width {
			drawX, drawY := g.mapToDrawCoords(x, y, offsetX, offsetY)
			drawImageAt(screen, groundImg, drawX, drawY)
		}
	}

	// DRAW EVERY OTHER ELEMENTS, SAME LOGIC IN A WAY BUT GOOD PRACTICE TO KEEP BOTH
	for y := range g.Map.Height {
		for x := range g.Map.Width {
			elementType := g.Map.GetElementType(x, y)
			if elementType == Map.VOID {
				continue
			}

			drawX, drawY := g.mapToDrawCoords(x, y, offsetX, offsetY)

			//If image not exist, will not render !
			switch elementType {
			case Map.WALL:
				drawImageAt(screen, wallImg, drawX, drawY)
			case Map.SHELF:
				drawImageAt(screen, itemImg, drawX, drawY)
			case Map.DOOR:
				drawImageAt(screen, doorImg, drawX, drawY)
			case Map.CHECKOUT:
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
	mapData, err := Map.LoadMapFromFile("maps/store/layout.txt")
	if err != nil {
		panic("Error loading map: " + err.Error())
	}

	return &Game{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		CameraX:      0,
		CameraY:      0,
		Map:          *mapData,
	}
}
