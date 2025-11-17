package Hud

import (
	Map "AI30_-_BlackFriday/pkg/map"
	"fmt"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func (h *HUD) Update(posX, posY int, elementType Map.ElementType, items []Map.Item, exists bool) {
	h.SelectedElement = elementType
	
	msg := fmt.Sprintf("Position: (%d, %d)\n", posX, posY)
	msg += fmt.Sprintf("Element Type: %s\n", elementType)

	if elementType == Map.SHELF {
		if exists {
			msg += fmt.Sprintf("Shelf Stock (%d items):\n", len(items))
			for i, item := range items {
				msg += fmt.Sprintf("  [%d] %s - Price: %.2f, Quantity: %d, Reduction: %.2f%%, Attractiveness: %.2f\n",
					i+1, item.Name, item.Price, item.Quantity, item.Reduction*100, item.Attractiveness)
			}
		} else {
			msg += fmt.Sprintf("No stock data available\n")
		}
	}

	h.DebugMsg = msg
	h.prepareRender()
}

// Determine the width and height the background based on the text
func (h *HUD) prepareRender() {
	lines := splitLines(h.DebugMsg)
	h.renderLines = lines

	lineHeight := FONT.Metrics().Height.Ceil()

	maxWidth := 0
	for _, line := range lines {
		bounds, _ := font.BoundString(FONT, line)
		width := (bounds.Max.X - bounds.Min.X).Ceil()
		if width > maxWidth {
			maxWidth = width
		}
	}

	h.HudWidth = maxWidth + h.PaddingX*2
	h.HudHeight = len(lines)*lineHeight + h.PaddingY*2

	h.HudBg = ebiten.NewImage(h.HudWidth, h.HudHeight)
	h.HudBg.Fill(color.RGBA{0, 0, 0, 180})
}

func (h *HUD) Draw(screen *ebiten.Image) {
	if h.HudBg == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(h.PositionX), float64(h.PositionY))
	screen.DrawImage(h.HudBg, op)

	y := h.PositionY + h.PaddingY + FONT.Metrics().Height.Ceil()
	for _, line := range h.renderLines {
		text.Draw(screen, line, FONT, h.PositionX+h.PaddingX, y, color.White)
		y += FONT.Metrics().Height.Ceil()
	}
}

func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
