package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/utils"
	"math"
)

// https://pedestriandynamics.org/models/social_force_model/ refer to
func CalculateSocialForces(agt Agent, neighbors []Agent) utils.Vec2 {
	socialForce := utils.Vec2{
		X: 0.0,
		Y: 0.0,
	}
	for _, neighbor := range neighbors {
		if agt.ID() == neighbor.ID() {
			continue
		}
		//relative postions
		agtCoord := agt.Coordinate()
		neighborCoord := neighbor.Coordinate()
		p := utils.Vec2{X: agtCoord.X - neighborCoord.X, Y: agtCoord.Y - neighborCoord.Y}
		distance := math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))
		sumRadius := constants.AGT_RADIUS * 2

		if distance < 0.0001 {
			continue
		}
		//normalized vector
		n := utils.Vec2{X: p.X / distance, Y: p.Y / distance}

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
			// tangent vector
			t := utils.Vec2{X: -n.Y, Y: n.X}
			// speed vector
			agtV := agt.Velocity()
			neighborV := neighbor.Velocity()
			vi := utils.Vec2{X: agtV.X - neighborV.X, Y: agtV.Y - neighborV.Y}
			deltaVT := vi.X*t.X + vi.Y*t.Y
			frictionMag := -constants.FRICTION_COEF * overlap * deltaVT

			socialForce.X += t.X * frictionMag
			socialForce.Y += t.Y * frictionMag
		}
	}

	return socialForce
}

func CalculateObstacleForces(agent Agent, obstacles []utils.Vec2) utils.Vec2 {
	totalForce := utils.Vec2{X: 0.0, Y: 0.0}
	coords := agent.Coordinate()

	for _, o := range obstacles {
		// Relative position vector
		p := utils.Vec2{X: coords.X - o.X, Y: coords.Y - o.Y}
		distance := math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))
		sumRadius := constants.AGT_RADIUS

		if distance < 0.0001 {
			continue
		}

		// Normalized vector
		n := utils.Vec2{X: p.X / distance, Y: p.Y / distance}

		socialForceMag := constants.WALL_RESISTANCE * math.Exp((sumRadius-distance)/constants.AGT_RANGE)
		totalForce.X += n.X * socialForceMag
		totalForce.Y += n.Y * socialForceMag

		if distance < sumRadius {
			overlap := sumRadius - distance

			contactForce := constants.AGT_STRENGTH * overlap
			totalForce.X += n.X * contactForce
			totalForce.Y += n.Y * contactForce

			// tangent vector
			t := utils.Vec2{X: -n.Y, Y: n.X}
			// speed vector
			agtV := agent.Velocity()
			deltaVT := agtV.X*t.X + agtV.Y*t.Y
			frictionMag := -constants.FRICTION_COEF * overlap * deltaVT
			totalForce.X += t.X * frictionMag
			totalForce.Y += t.Y * frictionMag

			correctionFactor := overlap + 0.000
			agent.Coordinate().X += n.X * correctionFactor
			agent.Coordinate().Y += n.Y * correctionFactor
		}
	}

	return totalForce
}

func ApplySocialForce(agt Agent, socialForce utils.Vec2, dt float64) {

	goalForce := utils.Vec2{
		X: (agt.DesiredVelocity().X - agt.Velocity().X) * constants.RELAXATION_FACTOR,
		Y: (agt.DesiredVelocity().Y - agt.Velocity().Y) * constants.RELAXATION_FACTOR,
	}

	totalAccX := goalForce.X + socialForce.X
	totalAccY := goalForce.Y + socialForce.Y

	agt.Velocity().X = agt.Velocity().X + totalAccX*dt
	agt.Velocity().Y = agt.Velocity().Y + totalAccY*dt

	currentSpeed := math.Sqrt(agt.Velocity().X*agt.Velocity().X + agt.Velocity().Y*agt.Velocity().Y)
	maxSpeed := agt.Speed() * constants.SPEED_MULTIPLIER

	if currentSpeed > maxSpeed {
		scale := maxSpeed / currentSpeed
		agt.Velocity().X *= scale
		agt.Velocity().Y *= scale
	}
}
