package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/logger"
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/pathfinding"
	"AI30_-_BlackFriday/pkg/utils"
	"fmt"
	"math/rand"
)

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

type ClientAgent struct {
	id           AgentID
	Speed        float64
	env          *Environment
	coordinate   utils.Vec2
	dx, dy       float64
	shoppingList []Map.Item
	cart         map[string]*Map.Item
	pickChan     chan PickRequest
	moveChan     chan MoveRequest

	syncChan chan int
	//temporaire
	moveChanResponse chan bool
	//rajouter un type action ?
	pickChanResponse chan PickResponse

	// Pathfinding
	currentPath              *pathfinding.Path
	moveTargetX, moveTargetY float64
	hasDestination           bool

	// Anti-blocage
	stuckCounter int
	lastPosition utils.Vec2

	velocity        utils.Vec2
	desiredVelocity utils.Vec2

	// Gestionnaires
	movementManager *MovementManager
	stuckDetector   *StuckDetector

	// Behavior
	state      AgentState
	nextAction ActionType
	itemsToPick                      []Map.Item
	interactTargetX, interactTargetY int

	// Manage vision and agent memory
	visionManager  *VisionManager
	visitedShelves map[[2]int]Map.Shelf
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
		coordinate:       utils.Vec2{X: float64(startX), Y: float64(startY)},
		dx:               0,
		dy:               0,
		shoppingList:     []Map.Item{}, //generateShoppingList(env),
		cart:             make(map[string]*Map.Item),
		pickChan:         pickChan,
		moveChan:         moveChan,
		syncChan:         syncChan,
		moveChanResponse: make(chan bool),
		pickChanResponse: make(chan PickResponse),
		hasDestination:   false,
		stuckCounter:     0,
		lastPosition:     utils.Vec2{X: float64(startX), Y: float64(startY)},
		state:            StateWandering,
		visitedShelves:   make(map[[2]int]Map.Shelf),
	}
	agent.movementManager = NewMovementManager(agent)
	agent.stuckDetector = NewStuckDetector(agent)
	agent.visionManager = NewVisionManager(agent)

	return agent
}

func generateShoppingList(env *Environment) []Map.Item {
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

func (ag *ClientAgent) State() AgentState {
	return ag.state
}

func (ag *ClientAgent) NextAction() ActionType {
	return ag.nextAction
}

func (ag *ClientAgent) ID() AgentID {
	return ag.id
}

func (ag *ClientAgent) Move() {
	ag.coordinate.X += ag.velocity.X
	ag.coordinate.Y += ag.velocity.Y
}
func (ag *ClientAgent) DryRunMove() utils.Vec2 {
	coordinate := ag.coordinate
	coordinate.X += ag.velocity.X
	coordinate.Y += ag.velocity.Y
	return coordinate
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

func (ag *ClientAgent) Coordinate() utils.Vec2 {
	return ag.coordinate
}

func (ag *ClientAgent) DesiredVelocity() *utils.Vec2 {
	return &ag.desiredVelocity
}

func (ag *ClientAgent) Velocity() *utils.Vec2 {
	return &ag.velocity
}

func (ag *ClientAgent) Direction() utils.Direction {
	return ag.movementManager.CalculateDirection()
}
func (ag *ClientAgent) Percept() {
	// ag.visionManager.UpdateFOVRays(ag.dx, ag.dy, 10, ag.env)
	ag.visionManager.UpdateFOV(ag.dx, ag.dy)
	ag.visionManager.DetectShelvesInFOV(ag.env)
}

func (ag *ClientAgent) Deliberate() {
	ag.stuckDetector.DetectAndResolve()

	switch ag.state {
	case StateWandering:

		// Agent has finished shopping
		if len(ag.GetMissingItems()) == 0 {
			destX, destY, res := FindWalkablePositionNearbyElement(ag.env, ag, "C")
			if !res {
				logger.Warnf("No walkable position nearby element")
				ag.nextAction = ActionWait
			} else {
				ag.movementManager.SetDestination(destX, destY)
				ag.state = StateMovingToCheckout
				ag.nextAction = ActionMove
			}
			break
		}

		// Agent tries to find a wanted item from a visited shelf
		if shelfX, shelfY, exists := ag.findWantedItemLocation(); exists {
			moveTargetX, moveTargetY, found := FindNearestFreePosition(ag.env, shelfX, shelfY)
			if !found {
				logger.Warnf("No path found to a location near this destination (%d %d) ", shelfX, shelfY)
				ag.nextAction = ActionWait
			} else {
				ag.interactTargetX, ag.interactTargetY = shelfX, shelfY
				ag.movementManager.SetDestination(moveTargetX, moveTargetY)
				ag.state = StateMovingToShelf
				ag.nextAction = ActionMove
			}
			break
		}

		// Agent wanders (visits unvisited nearby shelves)
		if !ag.hasDestination || ag.currentPath == nil || ag.currentPath.IsComplete() {
			destX, destY, res := FindWalkablePositionNearbyElement(ag.env, ag, "shelf")
			if res != true {
				logger.Warnf("No walkable position nearby element")
				ag.nextAction = ActionWait
			} else {
				ag.movementManager.SetDestination(destX, destY)
				ag.nextAction = ActionMove
			}
		}

		ag.processPathMovement()

	case StateMovingToShelf:
		ag.processPathMovement()
		if !ag.hasDestination {
			ag.state = StatePickingItem
			ag.nextAction = ActionWait
		}

	case StatePickingItem:
		ag.chooseItemToPickFromTargetedShelf(ag.env)

		if len(ag.itemsToPick) == 0 {
			logger.Warnf("No items to pick")
			ag.state = StateWandering
			break
		}
		ag.state = StateWandering
		ag.nextAction = ActionPick

	case StateMovingToCheckout:
		if len(ag.GetMissingItems()) != 0 { // si vol d'items ? (à implémenter plus tard)
			ag.state = StateWandering
			ag.nextAction = ActionWait
			break
		}
		ag.processPathMovement()
		if !ag.hasDestination {
			ag.state = StateCheckingOut
			ag.nextAction = ActionWait
		}
	case StateCheckingOut: // fait pas grand chose, peut-être voir pour refacto
		ag.state = StateMovingToExit
		ag.nextAction = ActionCheckout

	case StateMovingToExit:
		destX, destY, found := FindWalkablePositionNearbyElement(ag.env, ag, "D")
		if !found {
			logger.Warnf("Could not find walkable position nearby door")
			ag.nextAction = ActionWait
			break
		}
		ag.movementManager.SetDestination(destX, destY)
		ag.state = StateLeaving
		ag.nextAction = ActionMove

	case StateLeaving:
		ag.processPathMovement()
		if !ag.hasDestination {
			ag.nextAction = ActionExit
		}
	}
}

func (ag *ClientAgent) Act() {
	switch ag.nextAction {
	case ActionMove:
		ag.moveChan <- MoveRequest{Agt: ag, ResponseChannel: ag.moveChanResponse}
		<-ag.moveChanResponse

	case ActionPick:
		targetShelf := [2]int{ag.interactTargetX, ag.interactTargetY}
		for _, item := range ag.itemsToPick {
			ag.pickChan <- PickRequest{
				Agt:             ag,
				ItemName:        item.Name,
				ShelfCoords:     targetShelf,
				WantedAmount:    item.Quantity,
				ResponseChannel: ag.pickChanResponse,
			}
			amount := <-ag.pickChanResponse

			if !amount.Status {
				logger.Warnf("Bad request, PickRequest denied for (%s) at shelf (%d, %d)", item.Name, targetShelf[0], targetShelf[1])
				continue
			}

			if existingItem, ok := ag.cart[item.Name]; ok {
				existingItem.Quantity += amount.PickedAmount
			} else {
				newItem := Map.Item{
					Name:           item.Name,
					Price:          item.Price,
					Reduction:      item.Reduction,
					Attractiveness: item.Attractiveness,
					Quantity:       amount.PickedAmount,
				}
				ag.cart[item.Name] = &newItem
			}
		}

	case ActionCheckout:
		cartValue := ag.CalculateCartValue()
		if cartValue > 0 {
			ag.env.ProcessPayment(cartValue)
		}
		ag.cart = make(map[string]*Map.Item)

	case ActionWait:
		ag.dx = 0
		ag.dy = 0

	case ActionExit:
		fmt.Println("profit du magasin : ", ag.env.Profit)
		fmt.Println("exiting....")
	}
}

func (ag *ClientAgent) GetCurrentPath() *pathfinding.Path {
	return ag.currentPath
}

func (ag *ClientAgent) ShoppingList() []Map.Item {
	return ag.shoppingList
}

func (ag *ClientAgent) VisionManager() VisionManager {
	return *ag.visionManager
}

func (ag *ClientAgent) findWantedItemLocation() (int, int, bool) {
	missingItems := ag.GetMissingItems()

	for k, shelf := range ag.visitedShelves {
		for _, item := range shelf.Items {
			for _, wantedItem := range missingItems {
				if (item.Name == wantedItem.Name) && (item.Quantity > 0) {
					return k[0], k[1], true
				}
			}
		}
	}
	return 0, 0, false
}

func (ag *ClientAgent) processPathMovement() {
	// If path, follow it
	if ag.currentPath != nil && !ag.currentPath.IsComplete() {
		ag.movementManager.FollowPath()
	} else {
		// Stop movement
		ag.dx = 0
		ag.dy = 0
	}
}

func (ag *ClientAgent) chooseItemToPickFromTargetedShelf(env *Environment) {
	ag.itemsToPick = []Map.Item{}
	shelfCoords := [2]int{ag.interactTargetX, ag.interactTargetY}

	env.Mutex.RLock()
	defer env.Mutex.RUnlock()

	shelf, ok := env.Map.ShelfData[shelfCoords]

	if !ok {
		logger.Warnf("Shelf (%d %d) does not exist in the current environment", shelfCoords[0], shelfCoords[1])
		return
	}

	missingItems := ag.GetMissingItems()

	for _, shelfItem := range shelf.Items {
		for _, agentNeed := range missingItems {
			if (shelfItem.Name == agentNeed.Name) && (shelfItem.Quantity > 0) {
				ag.itemsToPick = append(ag.itemsToPick, agentNeed)
			}
		}
	}
}

func (ag *ClientAgent) GetMissingItems() []Map.Item {
	missing := []Map.Item{}

	for _, targetItem := range ag.shoppingList {
		have := 0
		if cartItem, ok := ag.cart[targetItem.Name]; ok {
			have = cartItem.Quantity
		}

		need := targetItem.Quantity - have

		if need > 0 {
			missing = append(missing, Map.Item{
				Name:           targetItem.Name,
				Price:          targetItem.Price,
				Reduction:      targetItem.Reduction,
				Attractiveness: targetItem.Attractiveness,
				Quantity:       need,
			})
		}
	}
	return missing
}

func (ag *ClientAgent) CalculateCartValue() float64 {
	var amount float64
	for _, item := range ag.cart {
		amount = amount + ((item.Price - item.Price*item.Reduction) * float64(item.Quantity))
	}
	return amount
}
