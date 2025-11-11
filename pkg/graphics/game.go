package Graphics

import (
	Simulation "AI30_-_BlackFriday/pkg/simulation"
)

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	Simulation                *Simulation.Simulation
}
