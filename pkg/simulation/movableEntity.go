package Simulation

import (
	"AI30_-_BlackFriday/pkg/utils"
)

type MovableEntity struct {
	coordinate      utils.Vec2
	velocity        utils.Vec2
	desiredVelocity utils.Vec2
	movementManager *MovementManager
	dx, dy          float64
	speed           float64
	lastPosition    utils.Vec2
}

func (ag *MovableEntity) Move() {
	ag.coordinate.X += ag.velocity.X
	ag.coordinate.Y += ag.velocity.Y
}
func (ag *MovableEntity) DryRunMove() utils.Vec2 {
	coordinate := ag.coordinate
	coordinate.X += ag.velocity.X
	coordinate.Y += ag.velocity.Y
	return coordinate
}

func (ag *MovableEntity) Coordinate() *utils.Vec2 {
	return &ag.coordinate
}

func (ag *MovableEntity) DesiredVelocity() *utils.Vec2 {
	return &ag.desiredVelocity
}

func (ag *MovableEntity) Velocity() *utils.Vec2 {
	return &ag.velocity
}

func (ag *MovableEntity) Direction() utils.Direction {
	return ag.movementManager.CalculateDirection()
}

func (ag *MovableEntity) Speed() float64 {
	return ag.speed
}
