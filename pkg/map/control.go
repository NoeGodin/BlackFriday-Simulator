package Map

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) HandleInput() {
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

	// HANDLE MOUS DLICK FOR DEBUGGING
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.handleMouseClick()
	}
}

func (g *Game) handleMouseClick() {
	mouseX, mouseY := ebiten.CursorPosition()

	// convert screen to map coordinates
	margin := 20
	mapX := (mouseX - margin + g.CameraX) / CellSize
	mapY := (mouseY - margin + g.CameraY) / CellSize

	if mapX >= 0 && mapX < g.Map.Width && mapY >= 0 && mapY < g.Map.Height {
		element := g.Map.Grid[mapY][mapX]

		fmt.Printf("=== DEBUG CLICK ===\n")
		fmt.Printf("Position: (%d, %d)\n", mapX, mapY)

		if element == nil {
			fmt.Printf("Element: nil\n")
		} else {
			fmt.Printf("Element Type: %s\n", element.Type())

			if element.Type() == SHELF {
				if shelf, ok := element.(*Shelf); ok {
					fmt.Printf("Shelf Stock (%d items):\n", len(shelf.Items))
					for i, item := range shelf.Items {
						fmt.Printf("  [%d] %s - Price: %.2f, Quantity: %d, Reduction: %.2f%%, Attractiveness: %.2f\n",
							i+1, item.Name, item.Price, item.Quantity, item.Reduction*100, item.Attractiveness)
					}
				}
			}
		}
		fmt.Printf("==================\n")
	}
}
