package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/logger"
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/pathfinding"
	"AI30_-_BlackFriday/pkg/utils"
	"math/rand"
)

type ClientAgent struct {
	id          AgentID
	Speed       float64
	env         *Environment
	coordinate  utils.Coordinate
	dx, dy      float64
	shoppingList []Map.Item
	pickChan   chan PickRequest
	moveChan   chan MoveRequest

	syncChan chan int
	//temporaire
	moveChanResponse chan bool
	//rajouter un type action ?

	// Pathfinding
	currentPath      *pathfinding.Path
	targetX, targetY int
	hasDestination   bool

	// Anti-blocage
	stuckCounter int
	lastPosition utils.Coordinate

	// Gestionnaires
	movementManager *MovementManager
	stuckDetector   *StuckDetector
}

func NewClientAgent(id string, env *Environment, moveChan chan MoveRequest, pickChan chan PickRequest, syncChan chan int) *ClientAgent {
	startX, startY, found := env.Map.GetRandomFreeCoordinate()
	if !found {
		startX, startY = 5, 5 // no free coordinate
	}

	agent := &ClientAgent{
		id:               AgentID(id),
		Speed:            constants.BASE_AGENT_SPEED,
		env:              env,
		coordinate:       utils.Coordinate{X: float64(startX), Y: float64(startY)},
		dx:               0,
		dy:               0,
		shoppingList:	  generateShoppingList(env),
		pickChan:         pickChan,
		moveChan:         moveChan,
		syncChan:         syncChan,
		moveChanResponse: make(chan bool),
		hasDestination:   false,
		stuckCounter:     0,
		lastPosition:     utils.Coordinate{X: float64(startX), Y: float64(startY)},
	}
	agent.movementManager = NewMovementManager(agent)
	agent.stuckDetector = NewStuckDetector(agent)

	return agent
}

func generateShoppingList(env *Environment) ([]Map.Item) {
	totalAttractiveness := 0.0
	shopList := []Map.Item{}
	for _, item := range env.Map.Items {
		totalAttractiveness += item.Attractiveness
	}

	for range rand.Intn(4) + 1 {
		wantedItem := rand.Float64() * totalAttractiveness
		cumulative := 0.0
		
		for _, item := range env.Map.Items {
			cumulative += item.Attractiveness
			if wantedItem <= cumulative {
				shopList = append(shopList, item)
				break
			}
		}
	}

	return shopList
}

func (ag *ClientAgent) ID() AgentID {
	return ag.id
}

func (ag *ClientAgent) Move() {
	ag.coordinate.X += ag.dx * ag.Speed
	ag.coordinate.Y += ag.dy * ag.Speed
}
func (ag *ClientAgent) DryRunMove() utils.Coordinate {
	coordinate := ag.coordinate
	coordinate.X += ag.dx * ag.Speed
	coordinate.Y += ag.dy * ag.Speed
	return ag.coordinate
}
func (ag *ClientAgent) Start() {
	logger.Infof("Agent %s starting at position (%.1f, %.1f)", ag.id, ag.coordinate.X, ag.coordinate.Y)

	go func() {
		var step int
		for {
			step = <-ag.syncChan
			// perception := <-ag.viewChan

			ag.Percept()
			ag.Deliberate()
			ag.Act()
			step++
			ag.syncChan <- step
		}
	}()
}

func (ag *ClientAgent) Coordinate() utils.Coordinate {
	return ag.coordinate
}

func (ag *ClientAgent) Direction() utils.Direction {
	return ag.movementManager.CalculateDirection()
}
func (ag *ClientAgent) Percept() {

}

func (ag *ClientAgent) Deliberate() {
	ag.stuckDetector.DetectAndResolve()

	// If no destination, generate a new one
	if !ag.hasDestination || ag.currentPath == nil || ag.currentPath.IsComplete() {
		ag.movementManager.GenerateNewDestination()
	}

	// If path, follow it
	if ag.currentPath != nil && !ag.currentPath.IsComplete() {
		ag.movementManager.FollowPath()
	} else {
		// Stop movement
		ag.dx = 0
		ag.dy = 0
	}
}

func (ag *ClientAgent) Act() {
	ag.moveChan <- MoveRequest{Agt: ag, ResponseChannel: ag.moveChanResponse}
	<-ag.moveChanResponse
}

func (ag *ClientAgent) GetCurrentPath() *pathfinding.Path {
	return ag.currentPath
}

func (ag *ClientAgent) Preferences() []Map.Item {
	return ag.shoppingList
}
