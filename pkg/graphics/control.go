package Graphics

import (
	Map "AI30_-_BlackFriday/pkg/map"
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
		elementType := g.Map.GetElementType(mapX, mapY)

		fmt.Printf("=== DEBUG CLICK ===\n")
		fmt.Printf("Position: (%d, %d)\n", mapX, mapY)

		fmt.Printf("Element Type: %s\n", elementType)

		if elementType == Map.SHELF {
			fmt.Printf("Product Zone at (%d, %d)\n", mapX, mapY)
			if items, exists := g.Map.GetProductData(mapX, mapY); exists {
				fmt.Printf("Shelf Stock (%d items):\n", len(items))
				for i, item := range items {
					fmt.Printf("  [%d] %s - Price: %.2f, Quantity: %d, Reduction: %.2f%%, Attractiveness: %.2f\n",
						i+1, item.Name, item.Price, item.Quantity, item.Reduction*100, item.Attractiveness)
				}
			} else {
				fmt.Printf("No stock data available\n")
			}
		}
		fmt.Printf("==================\n")
	}
}
