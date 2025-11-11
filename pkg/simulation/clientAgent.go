package Simulation

import (
	"log"
)

type Direction struct {
	X int
	Y int
}

var (
	NORTH = Direction{X: 0, Y: -1}
	EAST  = Direction{X: 1, Y: 0}
	SOUTH = Direction{X: 0, Y: 1}
	WEST  = Direction{X: -1, Y: 0}
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
	Direction  Direction
	pickChan   chan PickRequest
	moveChan   chan MoveRequest

	syncChan chan int
	//temporaire
	moveChanResponse chan bool
	//rajouter un type action ?

}

func NewClientAgent(id string, env *Environment, moveChan chan MoveRequest, pickChan chan PickRequest, syncChan chan int) *ClientAgent {
	return &ClientAgent{
		id:    AgentID(id),
		Speed: 0.05, env: env,
		coordinate:       Coordinate{X: 5, Y: 5},
		Direction:        NORTH,
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
	ag.coordinate.X += float64(ag.Direction.X) * ag.Speed
	ag.coordinate.Y += float64(ag.Direction.Y) * ag.Speed
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

func (ag *ClientAgent) Percept() {

}

func (ag *ClientAgent) Deliberate() {
}

func (ag *ClientAgent) Act() {
	ag.moveChan <- MoveRequest{Agt: ag, ResponseChannel: ag.moveChanResponse}
	<-ag.moveChanResponse
}
