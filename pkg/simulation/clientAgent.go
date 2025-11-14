package Simulation

import (
	"log"
	"math/rand"
)

type Direction int

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

type Coordinate struct {
	X float64
	Y float64
}
type ClientAgent struct {
	id         AgentID
	Speed      float64
	env        *Environment
	coordinate Coordinate
	dx, dy     float64
	pickChan   chan PickRequest
	moveChan   chan MoveRequest

	syncChan chan int
	//temporaire
	moveChanResponse chan bool
	//rajouter un type action ?

}

func NewClientAgent(id string, env *Environment, moveChan chan MoveRequest, pickChan chan PickRequest, syncChan chan int) *ClientAgent {
	return &ClientAgent{
		id:               AgentID(id),
		Speed:            BASE_AGENT_SPEED,
		env:              env,
		coordinate:       Coordinate{X: 5, Y: 5},
		dx:               0,
		dy:               0,
		pickChan:         pickChan,
		moveChan:         moveChan,
		syncChan:         syncChan,
		moveChanResponse: make(chan bool),
	}
}

func (ag *ClientAgent) ID() AgentID {
	return ag.id
}

func (ag *ClientAgent) Move() {
	ag.coordinate.X += ag.dx * ag.Speed
	ag.coordinate.Y += ag.dy * ag.Speed
}
func (ag *ClientAgent) DryRunMove() Coordinate {
	coordinate := ag.coordinate
	coordinate.X += ag.dx * ag.Speed
	coordinate.Y += ag.dy * ag.Speed
	return ag.coordinate
}
func (ag *ClientAgent) Start() {
	log.Printf("%s starting...\n", ag.id)

	go func() {
		var step int
		for {
			step = <-ag.syncChan
			// perception := <-ag.viewChan

			ag.Percept()
			ag.Deliberate()
			ag.Act()
			step++
			ag.syncChan <- step
		}
	}()
}

func (ag *ClientAgent) Coordinate() Coordinate {
	return ag.coordinate
}

func (ag *ClientAgent) Direction() Direction {
	if ag.dx == 0 && ag.dy == 0 || ag.dx == 0 && ag.dy > 0 {
		return SOUTH
	}
	if ag.dx > 0 && ag.dy == 0 {
		return EAST
	}
	if ag.dx < 0 && ag.dy == 0 {
		return WEST
	}
	if ag.dx == 0 && ag.dy < 0 {
		return NORTH
	}
	// nord-est
	if ag.dx > 0 && ag.dy < 0 {
		return NORTH
	}
	//nord-ouest
	if ag.dx < 0 && ag.dy < 0 {
		return NORTH
	}
	//sud-est
	if ag.dx > 0 && ag.dy > 0 {
		return SOUTH
	}
	//sd-ouest
	if ag.dx < 0 && ag.dy > 0 {
		return SOUTH
	}
	return SOUTH
}
func (ag *ClientAgent) Percept() {

}

func (ag *ClientAgent) Deliberate() {
	ag.dx = rand.Float64()*2 - 1
	ag.dy = rand.Float64()*2 - 1
}

func (ag *ClientAgent) Act() {
	ag.moveChan <- MoveRequest{Agt: ag, ResponseChannel: ag.moveChanResponse}
	<-ag.moveChanResponse
}
