package Simulation

import (
	Map "AI30_-_BlackFriday/pkg/map"
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

// demande pour bouger (peut être refuser si une personne où un objet n'est plus dispo)
func (env *Environment) moveRequest() {
	for moveRequest := range env.moveChan {
		moveRequest.Agt.Move()
		moveRequest.ResponseChannel <- true
	}
}

func (env *Environment) Start() {
	go env.pickRequest()
	go env.moveRequest()
}
