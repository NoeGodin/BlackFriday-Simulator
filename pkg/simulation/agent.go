package Simulation

import "AI30_-_BlackFriday/pkg/utils"

type AgentID string

type Agent interface {
	Start()
	Percept()
	Deliberate()
	Act()
	ID() AgentID
	Coordinate() utils.Coordinate
	Direction() utils.Direction
	Move()
	DryRunMove() utils.Coordinate
}
