package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	Simulation "AI30_-_BlackFriday/pkg/simulation"
	"fmt"
	"math"

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

	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.Hud.ToggleHidden()
	}
}

func (g *Game) handleMouseClick() {
	mouseX, mouseY := ebiten.CursorPosition()

	// convert screen to map coordinates
	margin := 20
	mapX := (mouseX - margin + g.CameraX) / constants.CELL_SIZE
	mapY := (mouseY - margin + g.CameraY) / constants.CELL_SIZE
	envMap := g.Simulation.Env.Map
	if mapX >= 0 && mapX < envMap.Width && mapY >= 0 && mapY < envMap.Height {
		elementType := envMap.GetElementType(mapX, mapY)

		items, exists := envMap.GetShelfData(mapX, mapY)
		agt := g.isMouseClickedOnAgent(mapX, mapY)
		g.Hud.SetSelection(float64(mapX), float64(mapY), &elementType, agt, items, exists)

		fmt.Printf("=== DEBUG CLICK ===\n")
		fmt.Printf("Position: (%d, %d)\n", mapX, mapY)

		fmt.Printf("Element Type: %s\n", elementType)

		if elementType == constants.SHELF {
			shelfChar := envMap.GetShelfCharacter(mapX, mapY)
			fmt.Printf("Shelf Zone at (%d, %d) - Shelf Type: '%s'\n", mapX, mapY, shelfChar)
			if shelf, exists := envMap.GetShelfData(mapX, mapY); exists {
				fmt.Printf("Shelf '%s' Categories: %v\n", shelfChar, shelf.Categories)
				fmt.Printf("Shelf '%s' Stock (%d items):\n", shelfChar, len(shelf.Items))
				for i, item := range shelf.Items {
					fmt.Printf("  [%d] %s - Price: %.2f, Quantity: %d, Reduction: %.2f%%, Attractiveness: %.2f\n",
						i+1, item.Name, item.Price, item.Quantity, item.Reduction*100, item.Attractiveness)
				}
			} else {
				fmt.Printf("No shelf data available\n")
			}
		}
		fmt.Printf("==================\n")
	}
}

func (g *Game) isMouseClickedOnAgent(mapX, mapY int) Simulation.Agent {
	for _, a := range g.Simulation.Agents() {
		dx := math.Abs(a.Coordinate().X - float64(mapX))
		dy := math.Abs(a.Coordinate().Y - float64(mapY))

		if dx < 0.5 && dy < 0.5 {
			return a
		}
	}

	return nil
}
