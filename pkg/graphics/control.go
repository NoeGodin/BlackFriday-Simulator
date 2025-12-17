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
	envMap := g.Simulation.Env.Map
	dx, dy := 0.0, 0.0
	agentDir := 0
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		dy = -1
		agentDir = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		dy = 1
		agentDir = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		dx = -1
		agentDir = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		dx = 1
		agentDir = 1
	}

	if dx != 0 || dy != 0 {
		targetX := g.Hud.TargetPositionX + dx
		targetY := g.Hud.TargetPositionY + dy

		shelf, exists := envMap.GetShelfData(targetX, targetY)
		elementType := envMap.GetElementType(targetX, targetY)

		g.Hud.SetSelection(
			targetX,
			targetY,
			&elementType,
			g.getHUDAdjacentAgent(agentDir),
			&shelf,
			exists,
		)
	}

	// HANDLE MOUSE CLICK FOR DEBUGGING
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.handleMouseClick()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.Hud.ToggleHidden()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.Hud.DisplayAgentPaths = !g.Hud.DisplayAgentPaths
	}
}

func (g *Game) handleMouseClick() {
	mouseX, mouseY := ebiten.CursorPosition()

	// convert screen to map coordinates
	margin := 20.0
	mapX := float64((mouseX - int(margin) + g.CameraX) / constants.CELL_SIZE)
	mapY := float64((mouseY - int(margin) + g.CameraY) / constants.CELL_SIZE)
	envMap := g.Simulation.Env.Map
	if mapX >= 0 && mapX < float64(envMap.Width) && mapY >= 0 && mapY < float64(envMap.Height) {
		elementType := envMap.GetElementType(mapX, mapY)

		items, exists := envMap.GetShelfData(mapX, mapY)
		agt := g.isMouseClickedOnAgent(mapX, mapY)
		g.Hud.SetSelection(mapX, mapY, &elementType, agt, &items, exists)

		fmt.Printf("=== DEBUG CLICK ===\n")
		fmt.Printf("Position: (%.0f, %.0f)\n", mapX, mapY)

		fmt.Printf("Element Type: %s\n", elementType)

		if elementType == constants.SHELF {
			shelfChar := envMap.GetShelfCharacter(mapX, mapY)
			fmt.Printf("Shelf Zone at (%.0f, %.0f) - Shelf Type: '%s'\n", mapX, mapY, shelfChar)
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

func (g *Game) isMouseClickedOnAgent(mapX, mapY float64) Simulation.Agent {
	for _, a := range g.Simulation.Agents() {
		dx := math.Abs(a.Coordinate().X - mapX)
		dy := math.Abs(a.Coordinate().Y - mapY)

		if dx < 0.5 && dy < 0.5 {
			return a
		}
	}

	return nil
}

func (g *Game) getHUDAdjacentAgent(direction int) Simulation.Agent {
	currentAgt := g.Hud.GetSelectedAgent()
	if currentAgt == nil {
		return nil
	}

	agents := g.Simulation.Agents()
	n := len(agents)
	if n == 0 {
		return nil
	}

	for i, agt := range agents {
		if agt == currentAgt {
			nextIndex := (i + direction + n) % n
			return agents[nextIndex]
		}
	}

	return agents[0]
}

func (g *Game) ensureValidHUDAgent() {
	agents := g.Simulation.Agents()

	if len(agents) == 0 {
		g.Hud.SetSelectedAgent(nil)
		return
	}

	current := g.Hud.GetSelectedAgent()
	if current == nil {
		g.Hud.SetSelectedAgent(agents[0])
		return
	}

	for _, agt := range agents {
		if agt == current {
			return
		}
	}

	g.Hud.SetSelectedAgent(agents[0])
}