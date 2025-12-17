package Hud

import (
	Map "AI30_-_BlackFriday/pkg/map"
	Simulation "AI30_-_BlackFriday/pkg/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

var FONT font.Face

type HUD struct {
	TargetPositionX, TargetPositionY float64
	PaddingX, PaddingY               int
	HudWidth, HudHeight              int
	HudBg                            *ebiten.Image
	Lines                            []string

	selectedElement *Map.ElementType
	selectedAgent   Simulation.Agent
	exists 			bool
	shelf			*Map.Shelf

	hidden 			  bool
	DisplayAgentPaths bool
}

func NewHud() *HUD {
	return &HUD{
		TargetPositionX: 10,
		TargetPositionY: 10,
		PaddingX:        10,
		PaddingY:        5,
		hidden:          true,
		DisplayAgentPaths: false,
	}
}

func (h *HUD) Hidden() bool {
	return h.hidden
}

func (h *HUD) ToggleHidden() {
	h.hidden = !h.hidden
}
