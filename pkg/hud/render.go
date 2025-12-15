package Hud

import (
	"AI30_-_BlackFriday/pkg/constants"
	Map "AI30_-_BlackFriday/pkg/map"
	Simulation "AI30_-_BlackFriday/pkg/simulation"
	"fmt"
	"strings"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (h *HUD) GetSelectedAgent() Simulation.Agent {
	if h.hidden {
		return nil
	}
	
	return h.selectedAgent
}

func (h *HUD) SetSelectedElement(posX, posY float64, elementType *Map.ElementType, shelf *Map.Shelf, exists bool) {
	h.selectedElement = elementType
	h.TargetPositionX = posX
	h.TargetPositionY = posY
	h.exists = exists
	h.shelf = shelf
}

func (h *HUD) SetSelectedAgent(agent Simulation.Agent) {
	h.selectedAgent = agent
	if agent == nil {
		return
	}
	
	h.TargetPositionX = agent.Coordinate().X
	h.TargetPositionY = agent.Coordinate().Y
}

func (h *HUD) SetSelection(posX, posY float64, elementType *Map.ElementType, agent Simulation.Agent, shelf *Map.Shelf, exists bool) {
	h.clearSelection()
	h.hidden = false

	msg := ""
	if agent == nil {
		h.SetSelectedElement(posX, posY, elementType, shelf, exists)
		msg = h.getElementSelectionMessage()
	} else {
		h.SetSelectedAgent(agent)
		msg = h.getAgentSelectionMessage()
	}

	h.prepareRender(msg)
}

// If an agent is selected, we refresh its position
func (h *HUD) Update() {
	if h.selectedAgent != nil {
		h.TargetPositionX = h.selectedAgent.Coordinate().X
		h.TargetPositionY = h.selectedAgent.Coordinate().Y
		msg := h.getAgentSelectionMessage()
		h.prepareRender(msg)
	}

	if h.selectedElement != nil {
		msg := h.getElementSelectionMessage()
		h.prepareRender(msg)
	}	
}

// Determine the width and height the background based on the text
func (h *HUD) prepareRender(msg string) {
	lines := strings.Split(msg, "\n")
	h.Lines = lines
	lineHeight := FONT.Metrics().Height.Ceil()
	maxWidth := 0
	for _, line := range lines {
		bounds := text.BoundString(FONT, line)
		width := bounds.Max.X + 1
		if width > maxWidth {
			maxWidth = width
		}
	}
	h.HudWidth = maxWidth + h.PaddingX*2
	h.HudHeight = len(lines)*lineHeight + h.PaddingY*2
	h.HudBg = ebiten.NewImage(h.HudWidth, h.HudHeight)
	h.HudBg.Fill(color.RGBA{0, 0, 0, 180})
}

func (h *HUD) clearSelection() {
	h.selectedElement = nil
	h.selectedAgent = nil
	h.shelf = nil
}

func (h *HUD) getAgentSelectionMessage() string {
	return h.selectedAgent.GetDisplayData()
}

func (h *HUD) getElementSelectionMessage() string {
	msg := fmt.Sprintf("Position: (%d, %d)\n", int(h.TargetPositionX), int(h.TargetPositionY))
	msg += fmt.Sprintf("Element Type: %s\n", *h.selectedElement)
	
	if *h.selectedElement == constants.SHELF {
		if h.exists {
			msg += fmt.Sprintf("Shelf Stock (%d items):\n", len(h.shelf.Items))
			for i, item := range h.shelf.Items {
				msg += fmt.Sprintf("  [%d] %s - Price: %.2f, Quantity: %d, Reduction: %.2f%%, Attractiveness: %.2f\n",
					i+1, item.Name, item.Price, item.Quantity, item.Reduction*100, item.Attractiveness)
			}
		} else {
			msg += "No stock data available\n"
		}
	}

	return msg
}
