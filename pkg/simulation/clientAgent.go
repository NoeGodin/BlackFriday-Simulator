package Simulation

import (
	"log"
	"math/rand"
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
	direction  Direction
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
		Speed: BASE_AGENT_SPEED, env: env,
		coordinate:       Coordinate{X: 5, Y: 5},
		direction:        NORTH,
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
	ag.coordinate.X += float64(ag.direction.X) * ag.Speed
	ag.coordinate.Y += float64(ag.direction.Y) * ag.Speed
}
func (ag *ClientAgent) DryRunMove() Coordinate {
	coordinate := ag.coordinate
	coordinate.X += float64(ag.direction.X) * ag.Speed
	coordinate.Y += float64(ag.direction.Y) * ag.Speed
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
	return ag.direction
}
func (ag *ClientAgent) Percept() {

}

func (ag *ClientAgent) Deliberate() {
	dir := rand.Intn(4)
	switch dir {
	case 0:
		ag.direction = EAST
	case 1:
		ag.direction = WEST
	case 2:
		ag.direction = NORTH
	case 3:
		ag.direction = SOUTH
	}
}

func (ag *ClientAgent) Act() {
	ag.moveChan <- MoveRequest{Agt: ag, ResponseChannel: ag.moveChanResponse}
	<-ag.moveChanResponse
}
