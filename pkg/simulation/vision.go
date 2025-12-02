package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/utils"
	"math"
)

type VisionManager struct {
	agent *ClientAgent

	// Don't directly use the constant for Distance and Height in case we can
	// to add agents obstruction (with walls or other agents)
	visionDistance float64
	visionHeight   float64
	P1, P2, P3, P4 utils.Vec2
	RaysEndPoints  []utils.Vec2
}

func NewVisionManager(ag *ClientAgent) *VisionManager {
	return &VisionManager{
		agent:          ag,
		visionDistance: constants.VISION_DISTANCE,
		visionHeight:   constants.VISION_HEIGHT,
	}
}

// Update for Raycast FOV
func (v *VisionManager) UpdateFOVRays(dx, dy float64, numRays int, env *Environment) {
    ax := v.agent.Coordinate().X + constants.CENTER_OF_CELL
    ay := v.agent.Coordinate().Y + constants.CENTER_OF_CELL
	v.RaysEndPoints = make([]utils.Vec2, numRays)

	fovAngle := constants.ANGLE_VISION * math.Pi / 180.0
	halfFOV := fovAngle / 2

	// Angle direction
	baseAngle := math.Atan2(dy, dx)
    for i := 0; i < numRays; i++ {
        angle := baseAngle - halfFOV + (float64(i) / float64(numRays - 1)) * fovAngle
        rayX, rayY := ax, ay
        step := 0.1

        for d := 0.0; d < v.visionDistance; d += step {
            rayX = math.Ceil(ax + math.Cos(angle) * d)
            rayY = math.Ceil(ay + math.Sin(angle) * d)

            if env.IsObstacleAt(rayX, rayY) {
                break
            }

            if env.IsShelfAt(rayX, rayY) {
                v.agent.visitedShelves[[2]int{int(rayX), int(rayY)}] = false
            }
        }

		v.RaysEndPoints[i] = utils.Vec2{X: rayX, Y: rayY}
    }
}

// Update for rectangle FOV
func (v *VisionManager) UpdateFOV(dx, dy float64) {
								 // Add with the center of the agent's sprite
	ax := v.agent.Coordinate().X + constants.CENTER_OF_CELL
	ay := v.agent.Coordinate().Y + constants.CENTER_OF_CELL

	fx := ax + dx * v.visionDistance
	fy := ay + dy * v.visionDistance

	px := -dy
	py := dx

	length := math.Sqrt(px * px + py * py)
	if length != 0 {
		px /= length
		py /= length
	}

	halfH := v.visionHeight / 2

	v.P1 = utils.Vec2{X: ax + px * halfH, Y: ay + py * halfH}
	v.P2 = utils.Vec2{X: ax - px * halfH, Y: ay - py * halfH}
	v.P3 = utils.Vec2{X: fx - px * halfH, Y: fy - py * halfH}
	v.P4 = utils.Vec2{X: fx + px * halfH, Y: fy + py * halfH}
}

func pointInTriangle(p, a, b, c utils.Vec2) bool {
	v0 := utils.Vec2{X: c.X - a.X, Y: c.Y - a.Y}
	v1 := utils.Vec2{X: b.X - a.X, Y: b.Y - a.Y}
	v2 := utils.Vec2{X: p.X - a.X, Y: p.Y - a.Y}

	dot00 := v0.X * v0.X + v0.Y * v0.Y
	dot01 := v0.X * v1.X + v0.Y * v1.Y
	dot02 := v0.X * v2.X + v0.Y * v2.Y
	dot11 := v1.X * v1.X + v1.Y * v1.Y
	dot12 := v1.X * v2.X + v1.Y * v2.Y

	invDenom := 1 / (dot00 * dot11 - dot01 * dot01)
	u := (dot11 * dot02 - dot01 * dot12) * invDenom
	v := (dot00 * dot12 - dot01 * dot02) * invDenom

	return (u >= 0) && (v >= 0) && (u+v < 1)
}

func (v *VisionManager) areCoordinatesIntersectingFOV(p utils.Vec2) bool {
	return pointInTriangle(p, v.P1, v.P2, v.P3) ||
		pointInTriangle(p, v.P1, v.P3, v.P4)
}

func (v *VisionManager) DetectShelvesInFOV(env *Environment) {
	// Keep in comment in case we want to know the items and shelves percepted by the agent
	// shelves := []Map.Shelf{}

	for coords := range env.Map.ShelfData {

		cx := float64(coords[0]) + constants.CENTER_OF_CELL
		cy := float64(coords[1]) + constants.CENTER_OF_CELL

		p := utils.Vec2{X: cx, Y: cy}

		if v.areCoordinatesIntersectingFOV(p) {
			v.agent.visitedShelves[[2]int{int(p.X), int(p.Y)}] = false
			// shelves = append(shelves, shelf)
		}
	}
	// for _, s := range shelves {
		//     for _, i := range s.Items {
			// 			fmt.Println(i.Name)
	//     }
	// }
}
