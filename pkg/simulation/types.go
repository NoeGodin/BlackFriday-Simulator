package Simulation

type AgentState int

const ( // enumeration
	StateWandering AgentState = iota
	StateMovingToShelf
	StatePickingItem
	StateMovingToCheckout
	StateCheckingOut
	StateMovingToExit
	StateLeaving
)

var agentStateName = map[AgentState]string{
	StateWandering:        "wandering",
	StateMovingToShelf:    "moving to shelf",
	StatePickingItem:      "picking item",
	StateMovingToCheckout: "moving to checkout",
	StateCheckingOut:      "checking out",
	StateMovingToExit:     "moving to exit",
	StateLeaving:          "leaving",
}

func (as AgentState) String() string {
	return agentStateName[as]
}

type ActionType int

const (
	ActionMove ActionType = iota
	ActionPick
	ActionWait
	ActionCheckout
	ActionExit
)

var actionTypeName = map[ActionType]string{
	ActionMove:     "moving",
	ActionPick:     "picking",
	ActionWait:     "waiting",
	ActionCheckout: "checking out",
	ActionExit:     "exiting",
}

func (at ActionType) String() string {
	return actionTypeName[at]
}