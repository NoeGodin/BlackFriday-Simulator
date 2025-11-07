package Simulation


type AgentID string

type Agent interface {
	Start()
	Percept(* Environment)
	Deliberate()
	Act(* Environment)
	ID() AgentID
}
