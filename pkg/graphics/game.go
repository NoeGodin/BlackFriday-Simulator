package Graphics

import (
	Simulation "AI30_-_BlackFriday/pkg/simulation"
	Hud "AI30_-_BlackFriday/pkg/Hud"
)

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	AgentAnimator             *AgentAnimator
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
		Simulation:    simu,
	}
}