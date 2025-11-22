package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	Simulation "AI30_-_BlackFriday/pkg/simulation"
	"AI30_-_BlackFriday/pkg/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimationState struct {
	animationFRAME_COUNT *[constants.DIRECTIONS][constants.FRAME_COUNT]*ebiten.Image
	step                 int
}
type AgentAnimator struct {
	agentStates map[Simulation.AgentID]*AnimationState
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
	state, ok := animator.agentStates[agt.ID()]
	if !ok {
		// si on implémente de nouveaux spites on pourrait décider aléatoirement de quel sprite attribuer à l'agent
		newState := &AnimationState{animationFRAME_COUNT: &WalkFrameImgs, step: 1}
		animator.agentStates[agt.ID()] = newState
		return newState.animationFRAME_COUNT[direction][0]
	}
	frame := (state.step / constants.FRAME_DURATION) % constants.FRAME_COUNT
	image := state.animationFRAME_COUNT[direction][frame]
	state.step++

	if state.step >= constants.FRAME_DURATION*constants.FRAME_COUNT {
		state.step = 0
	}
	return image
}

func NewAgentAnimator() *AgentAnimator {
	return &AgentAnimator{agentStates: make(map[Simulation.AgentID]*AnimationState)}
}
