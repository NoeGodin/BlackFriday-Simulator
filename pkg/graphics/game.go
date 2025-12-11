package Graphics

import (
	Hud "AI30_-_BlackFriday/pkg/hud"
	Simulation "AI30_-_BlackFriday/pkg/simulation"
)

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	AgentAnimator             *AgentAnimator
	ShelfAnimator             *ShelfAnimator
	Simulation                *Simulation.Simulation
	Hud                       Hud.HUD
}

func NewGame(screenWidth, screenHeight int, simu *Simulation.Simulation) *Game {

	return &Game{
		ScreenWidth:   screenWidth,
		ScreenHeight:  screenHeight,
		CameraX:       0,
		CameraY:       0,
		AgentAnimator: NewAgentAnimator(),
		ShelfAnimator: NewShelfAnimator(),
		Simulation:    simu,
		Hud:           *Hud.NewHud(),
	}
}
