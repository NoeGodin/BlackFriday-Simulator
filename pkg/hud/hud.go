package Hud

import (
	Map "AI30_-_BlackFriday/pkg/map"

	"github.com/hajimehoshi/ebiten/v2"
)

type HUD struct {
	PositionX, PositionY float64
	PaddingX, PaddingY   int
	HudWidth, HudHeight  int
	HudBg                *ebiten.Image
	Lines           []string
	SelectedElement Map.MapElement
}

func NewHud() *HUD {
	return &HUD{
		PositionX: 10,
		PositionY: 10,
		PaddingX:  10,
		PaddingY:  5,
	}
}