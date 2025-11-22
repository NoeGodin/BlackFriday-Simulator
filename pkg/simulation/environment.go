package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	Map "AI30_-_BlackFriday/pkg/map"
	"fmt"
	"math"
)

type PickRequest struct {
	Agt Agent
	//any est temporaire
	ResponseChannel any
}

type MoveRequest struct {
	Agt Agent
	//temporaire : structure à définir
	ResponseChannel chan bool
}

type Environment struct {
	Map      *Map.Map
	Clients  []*ClientAgent
	pickChan chan PickRequest
	moveChan chan MoveRequest
}

func NewEnvironment(mapData *Map.Map) *Environment {
	pickChan := make(chan PickRequest)
	moveChan := make(chan MoveRequest)
	return &Environment{Map: mapData, Clients: make([]*ClientAgent, 0), pickChan: pickChan, moveChan: moveChan}
}
func (env *Environment) AddClient(agtId string, syncChan chan int) Agent {
	client := NewClientAgent(agtId, env, env.moveChan, env.pickChan, syncChan)
	env.Clients = append(env.Clients, client)
	return client
}

// demande pour prendre un objet (peut être réfusé si l'objet n'est plus dispo)
func (env *Environment) pickRequest() {
	// for {
	// 	select {
	// 	case pickRequest := <-env.PickChannel:

	// 	}
	// }
}
func (env *Environment) isCollision(agt Agent) bool {
	coords := agt.DryRunMove()
	//check de la collision avec les eléments collisables
	for _, walls := range env.Map.GetCollisables() {
		offsetX := math.Abs(float64(walls[0]) - coords.X)
		offsetY := math.Abs(float64(walls[1]) - coords.Y)

		if offsetX < constants.AgentToEnvironmentHitbox && offsetY < constants.AgentToEnvironmentHitbox {
			return true
		}
	}
	return false
}

func (env *Environment) checkAgentCollisions(agt Agent) []*ClientAgent {
	coords := agt.DryRunMove()
	collidingAgents := make([]*ClientAgent, 0)

	for _, neighbor := range env.Clients {
		if agt.ID() == neighbor.ID() {
			continue
		}
		offsetX := math.Abs(neighbor.coordinate.X - coords.X)
		offsetY := math.Abs(neighbor.coordinate.Y - coords.Y)

		if offsetX < constants.AgentToAgentHitbox && offsetY < constants.AgentToAgentHitbox {
			collidingAgents = append(collidingAgents, neighbor)
		}
	}
	return collidingAgents
}

// demande pour bouger (peut être refuser si une personne où un objet n'est plus dispo)
func (env *Environment) moveRequest() {
	for moveRequest := range env.moveChan {
		// Check wall/shelf collisions (solid objects)
		isWallCollision := env.isCollision(moveRequest.Agt)
		if isWallCollision {
			// Block movement if hitting walls/shelves
			moveRequest.ResponseChannel <- true
			continue
		}
		collidingAgents := env.checkAgentCollisions(moveRequest.Agt)
		fmt.Println(collidingAgents)
		//TODO: ici on pourrait gérer mieux la collision en regardant l'etat de la var collidingAgents
		moveRequest.Agt.Move()
		moveRequest.ResponseChannel <- false // Movement always succeeds unless wall collision
	}
}

func (env *Environment) Start() {
	go env.pickRequest()
	go env.moveRequest()
}
