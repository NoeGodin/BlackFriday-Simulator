package Hud

import (
	Map "AI30_-_BlackFriday/pkg/map"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type HUD struct {
	PositionX, PositionY int
	PaddingX, PaddingY   int
	HudWidth, HudHeight  int
	HudBg                *ebiten.Image
	HudFont              font.Face

	DebugMsg        string
	renderLines     []string
	SelectedElement Map.MapElement
}
