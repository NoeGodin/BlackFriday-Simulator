package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/logger"
	"AI30_-_BlackFriday/pkg/pathfinding"
	"AI30_-_BlackFriday/pkg/utils"
	"context"
	"sync"
)

type AgentID string

type AgenType int

const (
	CLIENT AgenType = iota
	GUARD
)

type AgentBehavior interface {
	Percept()
	Deliberate()
	Act()
}

type Agent interface {
	AgentBehavior
	utils.ClickableEntity

	Start()
	ID() AgentID
	Coordinate() *utils.Vec2
	DesiredVelocity() *utils.Vec2
	Velocity() *utils.Vec2
	HasSpawned() bool
	GetCurrentPath() *pathfinding.Path
	Direction() utils.Direction
	Move()
	Speed() float64
	DryRunMove() utils.Vec2
	VisionManager() VisionManager
	Type() AgenType
}

type BaseAgent struct {
	MovableEntity
	utils.ClickableEntity

	id            AgentID
	agType        AgenType
	env           *Environment
	agentBehavior AgentBehavior

	moveChan  chan MoveRequest
	syncChan  chan int
	startChan chan StartRequest
	exitChan  chan ExitRequest
	stopCtx   context.Context
	stopWg    *sync.WaitGroup

	moveChanResponse  chan bool
	startChanResponse chan bool
	exitChanResponse  chan bool

	currentPath              *pathfinding.Path
	moveTargetX, moveTargetY float64
	hasDestination           bool
	hasSpawned               bool

	stuckCounter  int
	stuckDetector *StuckDetector
	visionManager *VisionManager
}

func NewBaseAgent(
	id string,
	pos [2]float64,
	env *Environment,
	moveChan chan MoveRequest,
	syncChan chan int,
	startChan chan StartRequest,
	exitChan chan ExitRequest,
	agtType AgenType,
	stopCtx context.Context,
	stopWg *sync.WaitGroup) *BaseAgent {

	agent := &BaseAgent{
		id:     AgentID(id),
		env:    env,
		agType: agtType,
		MovableEntity: MovableEntity{
			coordinate:   utils.Vec2{X: pos[0], Y: pos[1]},
			speed:        constants.BASE_AGENT_SPEED,
			dx:           0,
			dy:           0,
			lastPosition: utils.Vec2{X: pos[0], Y: pos[1]},
		},
		moveChan: moveChan,
		syncChan: syncChan,

		startChanResponse: make(chan bool),
		exitChanResponse:  make(chan bool),
		stopCtx:           stopCtx,
		stopWg:            stopWg,

		moveChanResponse: make(chan bool),
		hasDestination:   false,
		hasSpawned:       false,
		stuckCounter:     0,
		startChan:        startChan,
		exitChan:         exitChan,
	}
	agent.movementManager = NewMovementManager(agent)
	agent.stuckDetector = NewStuckDetector(agent)
	agent.visionManager = NewVisionManager(agent)

	return agent
}

func (ag *BaseAgent) ID() AgentID {
	return ag.id
}

func (ag *BaseAgent) HasSpawned() bool {
	return ag.hasSpawned
}
func (ag *BaseAgent) Type() AgenType {
	return ag.agType
}

func (ag *BaseAgent) Start() {
	if ag.startChan != nil {
		ag.startChan <- StartRequest{Agt: ag, ResponseChannel: ag.startChanResponse}
		ag.hasSpawned = <-ag.startChanResponse
	} else {
		ag.hasSpawned = true
	}

	logger.Infof("Agent %s starting at position (%.1f, %.1f)", ag.id, ag.coordinate.X, ag.coordinate.Y)

	go func() {
		defer ag.stopWg.Done()
		var step int
		for {
			select {
			case <-ag.stopCtx.Done():
				return
			case <-ag.exitChanResponse:
				logger.Infof("Agent %s finished", ag.id)
				return
			default:
				step = <-ag.syncChan
				// perception := <-ag.viewChan

				ag.Percept()
				ag.Deliberate()
				ag.Act()
				step++
				ag.syncChan <- step
			}
		}
	}()
}

func (ag *BaseAgent) Percept() {
	ag.agentBehavior.Percept()
}

func (ag *BaseAgent) Deliberate() {
	ag.agentBehavior.Deliberate()

}

func (ag *BaseAgent) Act() {
	ag.agentBehavior.Act()
}

func (ag *BaseAgent) GetCurrentPath() *pathfinding.Path {
	return ag.currentPath
}

func (ag *BaseAgent) VisionManager() VisionManager {
	return *ag.visionManager
}

func (ag *BaseAgent) processPathMovement() {
	// If path, follow it
	if ag.currentPath != nil && !ag.currentPath.IsComplete() {
		ag.movementManager.FollowPath()
	} else {
		// Stop movement
		ag.dx = 0
		ag.dy = 0
	}
}
