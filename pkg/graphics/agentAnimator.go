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
)

var (
	defaultColorScale = ebiten.ColorScale{}
	WinItemColorScale = func() ebiten.ColorScale {
		var cs ebiten.ColorScale
		cs.Scale(0.8, 2.5, 0.8, 1.0)
		return cs
	}()
)

type AnimationState struct {
	animationFrames *[constants.DIRECTIONS][constants.FRAME_COUNT]*ebiten.Image
	step            int
	agentState      AgentState
	itemCount       int
	stateStart      time.Time
}
type AgentAnimator struct {
	agentStates map[Simulation.AgentID]*AnimationState
}

func (animator *AgentAnimator) agentState(id Simulation.AgentID) (*AnimationState, bool) {
	state, ok := animator.agentStates[id]
	if !ok {
		// si on implémente de nouveaux spites on pourrait décider aléatoirement de quel sprite attribuer à l'agent
		newState := &AnimationState{animationFrames: &WalkFrameImgs, step: 0, agentState: None, itemCount: 0}
		animator.agentStates[id] = newState
		return newState, true
	} else {
		return state, false
	}
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
	state, new := animator.agentState(agt.ID())
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
	client, ok := agt.(*Simulation.ClientAgent)
	if !ok {
		return &ebiten.ColorScale{}
	}

	state, _ := animator.agentState(agt.ID())
	currentQuantity := client.CalculateCartQuantity()

	if state.itemCount < currentQuantity {
		state.agentState = WinItem
		state.stateStart = time.Now()
	}
	state.itemCount = currentQuantity

	// Effet buff temporaire
	if state.agentState == WinItem && time.Since(state.stateStart) < (constants.AGENT_STATE_DURATION) {
		return &WinItemColorScale
	}

	state.agentState = None
	return &defaultColorScale
}
