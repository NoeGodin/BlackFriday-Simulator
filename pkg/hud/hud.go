package Hud

import (
	Map "AI30_-_BlackFriday/pkg/map"
	Simulation "AI30_-_BlackFriday/pkg/simulation"

	"github.com/hajimehoshi/ebiten/v2"
)

type HUD struct {
	PositionX, PositionY float64
	PaddingX, PaddingY   int
	HudWidth, HudHeight  int
	HudBg                *ebiten.Image
	Lines                []string

	selectedElement Map.MapElement
	selectedAgent   Simulation.Agent

	hidden bool
}

func NewHud() *HUD {
	return &HUD{
		PositionX: 10,
		PositionY: 10,
		PaddingX:  10,
		PaddingY:  5,
		hidden:    true,
	}
}

func (h *HUD) Hidden() bool {
	return h.hidden
}

func (h *HUD) ToggleHidden() {
	h.hidden = !h.hidden
}