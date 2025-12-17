package Simulation

import "fmt"

type GuardAgent struct {
	*BaseAgent
}

func NewGuardAgent(id string, pos [2]float64, env *Environment, moveChan chan MoveRequest, startChan chan StartRequest, exitChan chan ExitRequest, syncChan chan int, agentIndex int) *GuardAgent {

	agent := &GuardAgent{
		BaseAgent: NewBaseAgent(id, pos, env, moveChan, syncChan, startChan, exitChan, GUARD),
	}
	agent.agentBehavior = &GuardAgentBehavior{ag: agent}

	return agent
}

type GuardAgentBehavior struct {
	ag *GuardAgent
}

func (va *GuardAgent) canSee(agt Agent) bool {
	return va.visionManager.areCoordinatesIntersectingFOV(*agt.Coordinate())
}
func (vb *GuardAgentBehavior) Percept() {
	ag := vb.ag
	ag.visionManager.UpdateFOV(ag.dx, ag.dy)
}

func (vb *GuardAgentBehavior) Deliberate() {
	ag := vb.ag
	ag.stuckDetector.DetectAndResolve()

	// If no destination, generate a new one
	if !ag.hasDestination || ag.currentPath == nil || ag.currentPath.IsComplete() {
		ag.movementManager.GenerateNewDestination()
	}

	// If path, follow it
	if ag.currentPath != nil && !ag.currentPath.IsComplete() {
		ag.movementManager.FollowPath()
	} else {
		// Stop movement
		ag.dx = 0
		ag.dy = 0
		ag.desiredVelocity.X = 0
		ag.desiredVelocity.Y = 0
	}

}

func (vb *GuardAgentBehavior) Act() {
	ag := vb.ag
	ag.moveChan <- MoveRequest{Agt: ag, ResponseChannel: ag.moveChanResponse}
	<-ag.moveChanResponse

}

func (vb *GuardAgent) GetDisplayData() string {
	msg := fmt.Sprintf("Agent guard: %s", vb.id)
	return msg
}