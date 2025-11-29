package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/utils"
	"math"
)

type Vec2 struct {
	X float64
	Y float64
}

// https://pedestriandynamics.org/models/social_force_model/ refer to
func CalculateSocialForces(agt *ClientAgent, neighbors []*ClientAgent) utils.Vec2 {
	socialForce := utils.Vec2{
		X: 0.0,
		Y: 0.0,
	}
	for _, neighbor := range neighbors {
		if agt.ID() == neighbor.ID() {
			continue
		}
		//positions relatives
		agtCoord := agt.Coordinate()
		neighborCoord := neighbor.Coordinate()
		p := Vec2{X: agtCoord.X - neighborCoord.X, Y: agtCoord.Y - neighborCoord.Y}
		distance := math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))
		sumRadius := constants.AGT_RADIUS * 2

		if distance < 0.0001 {
			continue
		}
		//vecteur normalisé
		n := Vec2{X: p.X / distance, Y: p.Y / distance}

		socialForceMag := constants.SOCIAL_STRENGTH * math.Exp((sumRadius-distance)/constants.AGT_RANGE)
		//pushing force 1
		socialForce.X += n.X * socialForceMag
		socialForce.Y += n.Y * socialForceMag

		if distance < sumRadius {
			overlap := sumRadius - distance

			//pushing force 2
			contactForce := constants.AGT_STRENGTH * overlap
			socialForce.X += n.X * contactForce
			socialForce.Y += n.Y * contactForce

			// sliding force
			// vecteur tengentiel
			t := Vec2{X: -n.Y, Y: n.X}
			//vecteur de vitesse
			agtV := agt.velocity
			neighborV := neighbor.velocity
			vi := Vec2{X: agtV.X - neighborV.X, Y: agtV.Y - neighborV.Y}
			deltaVT := vi.X*t.X + vi.Y*t.Y
			frictionMag := -constants.FRICTION_COEF * overlap * deltaVT

			socialForce.X += t.X * frictionMag
			socialForce.Y += t.Y * frictionMag
		}
	}

	return socialForce
}

func ApplySocialForce(agt *ClientAgent, socialForce utils.Vec2, dt float64) {

	goalForce := utils.Vec2{
		X: (agt.desiredVelocity.X - agt.velocity.X) * constants.RELAXATION_FACTOR,
		Y: (agt.desiredVelocity.Y - agt.velocity.Y) * constants.RELAXATION_FACTOR,
	}

	totalAccX := goalForce.X + socialForce.X
	totalAccY := goalForce.Y + socialForce.Y

	agt.velocity.X = agt.velocity.X + totalAccX*dt
	agt.velocity.Y = agt.velocity.Y + totalAccY*dt

	currentSpeed := math.Sqrt(agt.velocity.X*agt.velocity.X + agt.velocity.Y*agt.velocity.Y)
	maxSpeed := agt.Speed * constants.SPEED_MULTIPLIER

	if currentSpeed > maxSpeed {
		scale := maxSpeed / currentSpeed
		agt.velocity.X *= scale
		agt.velocity.Y *= scale
	}
}
