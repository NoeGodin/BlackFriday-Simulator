package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/shopping"
	"AI30_-_BlackFriday/pkg/utils"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type PickResponse struct {
	Status       bool
	PickedAmount int
}

type PickRequest struct { // une pick request par item
	Agt             Agent
	ItemName        string
	ShelfCoords     [2]float64
	WantedAmount    int
	ResponseChannel chan PickResponse
}

type MoveRequest struct {
	Agt Agent
	//temporaire : structure à définir
	ResponseChannel chan bool
}

type StartRequest struct {
	Agt             Agent
	ResponseChannel chan bool
}

type ExitRequest struct {
	Agt Agent
	//temporaire : structure à définir
	ResponseChannel chan bool
}

type Environment struct {
	Map                  *Map.Map
	Clients              []*ClientAgent
	Guards               []*GuardAgent
	Profit               float64
	pickChan             chan PickRequest
	moveChan             chan MoveRequest
	startChans           map[[2]float64]chan StartRequest
	exitChan             chan ExitRequest
	deltaTime            float64
	neighborSearchRadius float64
	Mutex                sync.RWMutex
	SalesTracker         *SalesTracker
	ShoppingListLoader   *shopping.ShoppingListLoader
	AgentCounter         int
}

func NewEnvironment(mapData *Map.Map, deltaTime float64, searchRadius float64, mapName string, shoppingListsPath string) *Environment {
	pickChan := make(chan PickRequest)
	moveChan := make(chan MoveRequest)
	startChan := make(map[[2]float64]chan StartRequest)
	for _, co := range mapData.Doors {
		startChan[co] = make(chan StartRequest)
	}

	exitChan := make(chan ExitRequest)

	salesTracker := NewSalesTracker(mapName)

	// Load shoppingList
	shoppingListLoader, err := shopping.NewShoppingListLoader(shoppingListsPath)
	if err != nil {
		// If file doesn't exist, loading old methods
		shoppingListLoader = nil
	}

	return &Environment{
		Map:                  mapData,
		Clients:              make([]*ClientAgent, 0),
		pickChan:             pickChan,
		startChans:           startChan,
		moveChan:             moveChan,
		exitChan:             exitChan,
		deltaTime:            deltaTime,
		neighborSearchRadius: searchRadius,
		SalesTracker:         salesTracker,
		ShoppingListLoader:   shoppingListLoader,
		AgentCounter:         0,
	}
}
func (env *Environment) AddClient(agtId string, syncChan chan int) Agent {
	doorCo := env.GetRandomDoor()

	client := NewClientAgent(agtId, env.getSpawnablePos(doorCo), env, env.moveChan, env.pickChan, env.startChans[doorCo], env.exitChan, syncChan, env.AgentCounter)
	env.Clients = append(env.Clients, client)
	env.AgentCounter++
	return client
}
func (env *Environment) AddGuard(agtId string, syncChan chan int) (Agent, error) {
	x, y, freePostion := env.Map.GetRandomFreeCoordinate()
	if !freePostion {
		return nil, fmt.Errorf("Error generating agent %s, no spawnable position", agtId)
	}
	pos := [2]float64{x, y}
	guard := NewGuardAgent(agtId, pos, env, env.moveChan, nil, env.exitChan, syncChan, env.AgentCounter)
	env.Guards = append(env.Guards, guard)
	env.AgentCounter++
	return guard, nil
}

func (env *Environment) getSpawnablePos(co [2]float64) [2]float64 {
	if co[0] == 0 {
		co[0] += constants.SPAWN_OFFSET_FROM_DOOR
	} else if co[1] == 0 {
		co[1] += constants.SPAWN_OFFSET_FROM_DOOR
	} else if co[1] == float64(env.Map.Height-1) {
		co[1] -= constants.SPAWN_OFFSET_FROM_DOOR
	} else if co[0] == float64(env.Map.Width-1) {
		co[0] -= constants.SPAWN_OFFSET_FROM_DOOR
	}

	return co
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

func (env *Environment) getNearbyAgents(agt Agent, radius float64) []Agent {
	nearbyAgents := make([]Agent, 0)
	for _, neighbor := range env.Clients {
		if agt.ID() == neighbor.ID() {
			continue
		}
		if agt.Coordinate().Distance(*neighbor.Coordinate()) <= radius {
			nearbyAgents = append(nearbyAgents, neighbor)
		}
	}
	return nearbyAgents
}

func (env *Environment) getNearbyCollisables(agt Agent, radius float64) []utils.Vec2 {
	nearbyCollisables := make([]utils.Vec2, 0)
	for _, collisable := range env.Map.GetCollisables() {
		point := ClosestPointToObstacle(*agt.Coordinate(), utils.Vec2{X: float64(collisable[0]), Y: float64(collisable[1])})
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
		clientAgent := moveRequest.Agt

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

func (env *Environment) pickRequest() {
	for pickRequest := range env.pickChan {
		if pickRequest.ItemName == "" {
			pickRequest.ResponseChannel <- PickResponse{false, 0}
			continue
		}

		if pickRequest.WantedAmount <= 0 {
			pickRequest.ResponseChannel <- PickResponse{false, 0}
			continue
		}

		env.Mutex.Lock() //TODO : potentiellement à refacto

		shelf, ok := env.Map.ShelfData[pickRequest.ShelfCoords]
		if !ok {
			env.Mutex.Unlock()
			pickRequest.ResponseChannel <- PickResponse{false, 0}
			continue
		}

		targetedShelf := shelf.Items
		var itemToPick *Map.Item

		for itemIndex := range targetedShelf {
			if targetedShelf[itemIndex].Name == pickRequest.ItemName {
				itemToPick = &targetedShelf[itemIndex]
				break
			}
		}

		if itemToPick == nil {
			env.Mutex.Unlock()
			pickRequest.ResponseChannel <- PickResponse{false, 0}
			continue
		}
		if itemToPick.Quantity >= pickRequest.WantedAmount {
			pickRequest.ResponseChannel <- PickResponse{true, pickRequest.WantedAmount}
			itemToPick.Quantity -= pickRequest.WantedAmount

		} else {
			pickRequest.ResponseChannel <- PickResponse{true, itemToPick.Quantity}
			itemToPick.Quantity = 0
		}
		env.Mutex.Unlock()
	}
}

// Launch all start requests from all doors
func (env *Environment) startRequests() {
	for _, startRequest := range env.startChans {
		go func(startRequest chan StartRequest) {
			for s := range startRequest {
				s.ResponseChannel <- true
				time.Sleep(constants.AGENT_SPAWN_INTERVAL_MS * time.Millisecond)
			}
		}(startRequest)
	}
}

func remove(name Agent, nations []*ClientAgent) []*ClientAgent {
	i := 0
	for idx, item := range nations {
		if item.ID() != name.ID() {
			nations[i] = nations[idx]
			i++
		}
	}
	return nations[:i]
}

func removeAgentFromClients(agentID AgentID, clients []*ClientAgent) []*ClientAgent {
	i := 0
	for _, c := range clients {
		if c.ID() != agentID {
			clients[i] = c
			i++
		}
	}
	return clients[:i]
}

func (env *Environment) Start() {
	go env.pickRequest()
	go env.moveRequest()
	env.startRequests()
}

func (env *Environment) IsObstacleAt(x, y float64) bool {
	for _, wall := range env.Map.GetCollisables() {
		if math.Abs(float64(wall[0])-x) < constants.CENTER_OF_CELL && math.Abs(float64(wall[1])-y) < constants.CENTER_OF_CELL {
			return true
		}
	}
	return false
}

func (env *Environment) IsShelfAt(x, y float64) bool {
	for coords := range env.Map.ShelfData {
		if x == coords[0] && y == coords[1] {
			return true
		}
	}
	return false
}

func (env *Environment) ProcessPayment(amout float64) {
	env.Mutex.Lock()
	defer env.Mutex.Unlock()
	env.Profit += amout
	env.SalesTracker.RecordSale(amout, env.Profit)
}

func (env *Environment) ExportSalesData() error {
	return env.SalesTracker.ExportToCSV()
}

func (env *Environment) GetRandomDoor() [2]float64 {
	idx := rand.Intn(len(env.Map.Doors))
	return env.Map.Doors[idx]
}

func (env *Environment) GenerateShoppingListDeterministic(agentIndex int) []Map.Item {
	// If list, loading it
	if env.ShoppingListLoader != nil {
		predefList := env.ShoppingListLoader.GetShoppingList(agentIndex)
		if len(predefList) > 0 {
			return predefList
		}
	}

	// Random option if not loaded
	return env.generateShoppingList()
}

func (env *Environment) generateShoppingList() []Map.Item {
	totalAttractiveness := 0.0
	shopList := []Map.Item{}
	for _, item := range env.Map.Items {
		totalAttractiveness += item.Attractiveness
	}

	for range rand.Intn(constants.AGENT_MAX_SHOPPING_LIST) + 1 {
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
