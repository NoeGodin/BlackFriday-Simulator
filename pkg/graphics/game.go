package Graphics

import (
	Hud "AI30_-_BlackFriday/pkg/hud"
	Simulation "AI30_-_BlackFriday/pkg/simulation"

	"github.com/ebitenui/ebitenui"
)

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	AgentAnimator             *AgentAnimator
	ShelfAnimator             *ShelfAnimator
	Simulation                *Simulation.Simulation
	Hud                       Hud.HUD
	UI                        *ebitenui.UI
}

func NewGame(screenWidth, screenHeight int, simu *Simulation.Simulation) *Game {
	game := &Game{
		ScreenWidth:   screenWidth,
		ScreenHeight:  screenHeight,
		CameraX:       0,
		CameraY:       0,
		AgentAnimator: NewAgentAnimator(),
		ShelfAnimator: NewShelfAnimator(),
		Simulation:    simu,
		Hud:           *Hud.NewHud(),
	}

	game.UI = createUI(game)

	return game
}
