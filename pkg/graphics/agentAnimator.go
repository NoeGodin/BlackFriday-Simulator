package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	Simulation "AI30_-_BlackFriday/pkg/simulation"
	"AI30_-_BlackFriday/pkg/utils"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type AgentState int

const (
	None AgentState = iota
	WinItem
	LooseItem
)

type AnimationState struct {
	animationFrames WalkAnimation
	step            int
	agentState      AgentState
	itemCount       int
	stateStart      time.Time
}
type AgentAnimator struct {
	agentStates map[Simulation.AgentID]*AnimationState
}

func (animator *AgentAnimator) agentState(agt Simulation.Agent) (*AnimationState, bool) {
	state, ok := animator.agentStates[agt.ID()]
	if !ok {
		// si on implémente de nouveaux spites on pourrait décider aléatoirement de quel sprite attribuer à l'agent
		newState := &AnimationState{animationFrames: getWalkAnimation(agt.Type()), step: 0, agentState: None, itemCount: 0}
		animator.agentStates[agt.ID()] = newState
		return newState, true
	} else {
		if state.animationFrames[0][0] == nil {
			state.animationFrames = getWalkAnimation(agt.Type())
		}
		return state, false
	}
}
func (animator *AgentAnimator) GetEmotion(agt Simulation.Agent) *ebiten.Image {
	if agt.Type() == Simulation.CLIENT {
		if client, ok := agt.(*Simulation.ClientAgent); ok {
			if client.Aggressiveness() >= constants.AGENT_AGGRESSIVENESS_TRESHOLD {
				return angryEmotion
			}
		}
	}
	return nil
}

// il faudra reset les steps lorsque l'agent n'est pas en déplacement
func (animator *AgentAnimator) AnimationFrame(agt Simulation.Agent) *ebiten.Image {
	direction := 0
	switch agt.Direction() {
	case utils.SOUTH:
		direction = 0
	case utils.NORTH:
		direction = 1

	case utils.EAST:
		direction = 2

	case utils.WEST:
		direction = 3

	}
	state, new := animator.agentState(agt)
	if new {
		return state.animationFrames[direction][0]
	}
	frame := (state.step / constants.FRAME_DURATION) % constants.FRAME_COUNT
	image := state.animationFrames[direction][frame]
	state.step++

	if state.step >= constants.FRAME_DURATION*constants.FRAME_COUNT {
		state.step = 0
	}

	return image
}

func NewAgentAnimator() *AgentAnimator {
	return &AgentAnimator{agentStates: make(map[Simulation.AgentID]*AnimationState)}
}

func (animator *AgentAnimator) getColorScale(agt Simulation.Agent) *ebiten.ColorScale {
	if agt.Type() != Simulation.CLIENT {
		return &ebiten.ColorScale{}
	}
	client, ok := agt.(*Simulation.ClientAgent)
	if !ok {
		return &ebiten.ColorScale{}
	}

	state, _ := animator.agentState(agt)
	currentQuantity := client.CalculateCartQuantity()

	if state.itemCount < currentQuantity {
		state.agentState = WinItem
		state.stateStart = time.Now()
	} else if state.itemCount > currentQuantity {
		state.agentState = LooseItem
		state.stateStart = time.Now()
	}
	state.itemCount = currentQuantity

	// Effet buff temporaire
	if time.Since(state.stateStart) < (constants.AGENT_STATE_DURATION) {
		if state.agentState == WinItem {
			return &WinItemColorScale
		} else if state.agentState == LooseItem {
			return &LooseItemColorScale
		}
	}

	state.agentState = None
	return &defaultColorScale
}
