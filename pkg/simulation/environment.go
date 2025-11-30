package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/utils"
	"math"
)

type PickRequest struct {
	Agt Agent
	//any est temporaire
	ResponseChannel any
}

type MoveRequest struct {
	Agt Agent
	//temporaire : structure à définir
	ResponseChannel chan bool
}

type Environment struct {
	Map                  *Map.Map
	Clients              []*ClientAgent
	pickChan             chan PickRequest
	moveChan             chan MoveRequest
	deltaTime            float64
	neighborSearchRadius float64
}

func NewEnvironment(mapData *Map.Map, deltaTime float64, searchRadius float64) *Environment {
	pickChan := make(chan PickRequest)
	moveChan := make(chan MoveRequest)
	return &Environment{Map: mapData, Clients: make([]*ClientAgent, 0), pickChan: pickChan, moveChan: moveChan, deltaTime: deltaTime, neighborSearchRadius: searchRadius}
}
func (env *Environment) AddClient(agtId string, syncChan chan int) Agent {
	client := NewClientAgent(agtId, env, env.moveChan, env.pickChan, syncChan)
	env.Clients = append(env.Clients, client)
	return client
}

// demande pour prendre un objet (peut être réfusé si l'objet n'est plus dispo)
func (env *Environment) pickRequest() {
	// for {
	// 	select {
	// 	case pickRequest := <-env.PickChannel:

	// 	}
	// }
}
func (env *Environment) isCollision(agt Agent) bool {
	coords := agt.DryRunMove()
	//check de la collision avec les eléments collisables
	for _, walls := range env.Map.GetCollisables() {
		offsetX := math.Abs(float64(walls[0]) - coords.X)
		offsetY := math.Abs(float64(walls[1]) - coords.Y)

		if offsetX < constants.AgentToEnvironmentHitbox && offsetY < constants.AgentToEnvironmentHitbox {
			return true
		}
	}
	return false
}

func (env *Environment) checkAgentCollisions(agt Agent) []*ClientAgent {
	coords := agt.DryRunMove()
	collidingAgents := make([]*ClientAgent, 0)

	for _, neighbor := range env.Clients {
		if agt.ID() == neighbor.ID() {
			continue
		}
		offsetX := math.Abs(neighbor.coordinate.X - coords.X)
		offsetY := math.Abs(neighbor.coordinate.Y - coords.Y)

		if offsetX < constants.AgentToAgentHitbox && offsetY < constants.AgentToAgentHitbox {
			collidingAgents = append(collidingAgents, neighbor)
		}
	}
	return collidingAgents
}

func (env *Environment) getNearbyAgents(agt Agent, radius float64) []*ClientAgent {
	nearbyAgents := make([]*ClientAgent, 0)
	for _, neighbor := range env.Clients {
		if agt.ID() == neighbor.ID() {
			continue
		}
		if agt.Coordinate().Distance(neighbor.Coordinate()) <= radius {
			nearbyAgents = append(nearbyAgents, neighbor)
		}
	}
	return nearbyAgents
}

func (env *Environment) getNearbyCollisables(agt Agent, radius float64) []utils.Vec2 {
	nearbyCollisables := make([]utils.Vec2, 0)
	for _, collisable := range env.Map.GetCollisables() {
		point := ClosestPointToObstacle(agt.Coordinate(), utils.Vec2{X: collisable[0], Y: collisable[1]})
		if agt.Coordinate().Distance(point) <= radius {
			nearbyCollisables = append(nearbyCollisables, point)
		}
	}
	return nearbyCollisables
}

func ClosestPointToObstacle(agentPos utils.Vec2, obstacle utils.Vec2) (pos utils.Vec2) {
	minX := obstacle.X - 0.5
	maxX := obstacle.X + 0.5
	minY := obstacle.Y - 0.5
	maxY := obstacle.Y + 0.5

	if agentPos.X < minX {
		pos.X = minX
	} else if agentPos.X > maxX {
		pos.X = maxX
	} else {
		pos.X = agentPos.X
	}
	if agentPos.Y < minY {
		pos.Y = minY
	} else if agentPos.Y > maxY {
		pos.Y = maxY
	} else {
		pos.Y = agentPos.Y
	}
	return
}

// demande pour bouger (peut être refuser si une personne où un objet n'est plus dispo)
func (env *Environment) moveRequest() {
	for moveRequest := range env.moveChan {
		clientAgent, ok := moveRequest.Agt.(*ClientAgent)
		if !ok {
			isWallCollision := env.isCollision(moveRequest.Agt)
			if isWallCollision {
				moveRequest.ResponseChannel <- true
				continue
			}
			moveRequest.Agt.Move()
			moveRequest.ResponseChannel <- false
			continue
		}

		neighbors := env.getNearbyAgents(clientAgent, env.neighborSearchRadius)
		collisables := env.getNearbyCollisables(clientAgent, env.neighborSearchRadius*2)

		socialForces := CalculateSocialForces(clientAgent, neighbors)
		obstaclesForces := CalculateObstacleForces(clientAgent, collisables)
		socialForces.X += obstaclesForces.X
		socialForces.Y += obstaclesForces.Y
		ApplySocialForce(clientAgent, socialForces, env.deltaTime)

		clientAgent.Move()
		moveRequest.ResponseChannel <- true
	}
}

func (env *Environment) Start() {
	go env.pickRequest()
	go env.moveRequest()
}
