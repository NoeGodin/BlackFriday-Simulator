package Simulation

import (
	Map "AI30_-_BlackFriday/pkg/map"
)

type ViewRequest struct {
	Agt *Agent
	//any est temporaire
	ResponseChannel any
}

type PickRequest struct {
	Agt *Agent
	//any est temporaire
	ResponseChannel any
}

type MooveRequest struct {
	Agt *Agent
	//any est temporaire
	ResponseChannel any
}

type moveRequest struct {
	Agt *Agent
}
type Environment struct {
	Map          *Map.Map
	Clients      []Agent
	ViewChannel  chan ViewRequest
	PickChannel  chan ViewRequest
	MooveChannel chan ViewRequest
}

func NewEnvironment(width int, height int) *Environment {
	return &Environment{Map: Map.NewMap(width, height)}
}
func (env *Environment) AddClient(agt Agent) {
	env.Clients = append(env.Clients, agt)
}

// fournir à l'agent sa "vision" qst : donner toute la map ? une partie ? des points d'intérets ?
func (env *Environment) viewRequest() {
	for {
		select {
		case viewRequest := <-env.ViewChannel:

		}
	}
}

// demande pour prendre un objet (peut être réfusé si l'objet n'est plus dispo)
func (env *Environment) pickRequest() {
	for {
		select {
		case pickRequest := <-env.PickChannel:

		}
	}
}

// demande pour bouger (peut être refuser si une personne où un objet n'est plus dispo)
func (env *Environment) mooveRequest() {
	for {
		select {
		case mooveRequest := <-env.MooveChannel:

		}
	}
}

func (env *Environment) Start() {
	go env.viewRequest()
	go env.pickRequest()
	go env.mooveRequest()
}
