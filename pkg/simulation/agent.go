package Simulation

type AgentID string

type Agent interface {
	Start()
	Percept()
	Deliberate()
	Act()
	ID() AgentID
	Coordinate() Coordinate
	Move()
}
