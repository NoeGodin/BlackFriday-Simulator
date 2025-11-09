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
	X int
	Y int
}
type ClientAgent struct {
	id         AgentID
	Speed      float32
	env        *Environment
	Coordinate Coordinate
	Direction  Direction
	viewChan   chan ViewRequest
	pickChan   chan PickRequest
	moveChan   chan MoveRequest
	syncChan   chan int

	//rajouter un type action ?

}

func NewClientAgent(id string, env *Environment, viewChan chan ViewRequest, moveChan chan MoveRequest, pickChan chan PickRequest, syncChan chan int) *ClientAgent {
	return &ClientAgent{AgentID(id), 1.0, env, Coordinate{X: 5, Y: 5}, NORTH, viewChan, pickChan, moveChan, syncChan}
}

func (ag *ClientAgent) ID() AgentID {
	return ag.id
}

func (ag *ClientAgent) Start() {
	log.Printf("%s starting...\n", ag.id)

	go func() {
		env := ag.env
		var step int
		for {
			step = <-ag.syncChan
			perception := <-ag.viewChan

			ag.Percept(env)
			ag.Deliberate()
			ag.Act(env)
			step++
			ag.syncChan <- step
		}
	}()
}

func (ag *ClientAgent) Percept(env *Environment) {

}

func (ag *ClientAgent) Deliberate() {
}

func (ag *ClientAgent) Act(env *Environment) {

}
