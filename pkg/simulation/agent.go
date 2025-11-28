package Simulation

import (
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/utils"
)

type AgentID string

type Agent interface {
	Start()
	Percept()
	Deliberate()
	Act()
	ID() AgentID
	Coordinate() utils.Vec2
	DesiredVelocity() *utils.Vec2
	Velocity() *utils.Vec2

	Direction() utils.Direction
	Move()
	DryRunMove() utils.Vec2
	ShoppingList() []Map.Item
	VisionManager() VisionManager
}
