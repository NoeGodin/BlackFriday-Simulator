package Simulation

import (
	"AI30_-_BlackFriday/pkg/constants"
	"AI30_-_BlackFriday/pkg/logger"
	Map "AI30_-_BlackFriday/pkg/map"
	"AI30_-_BlackFriday/pkg/utils"
	"fmt"
	"math"
)

type ClientAgent struct {
	*BaseAgent
	shoppingList []Map.Item
	cart         map[string]*Map.Item

	pickChan chan PickRequest

	pickChanResponse chan PickResponse

	// Behavior
	state                            AgentState
	nextAction                       ActionType
	itemsToPick                      []Map.Item
	interactTargetX, interactTargetY float64

	visitedShelves map[[2]float64]Map.Shelf
}

func NewClientAgent(id string, pos [2]float64, env *Environment, moveChan chan MoveRequest, pickChan chan PickRequest, startChan chan StartRequest, exitChan chan ExitRequest, syncChan chan int, agentIndex int) *ClientAgent {

	agent := &ClientAgent{
		BaseAgent:        NewBaseAgent(id, pos, env, moveChan, syncChan, startChan, exitChan),
		shoppingList:     env.GenerateShoppingListDeterministic(agentIndex),
		cart:             make(map[string]*Map.Item),
		pickChan:         pickChan,
		pickChanResponse: make(chan PickResponse),
		state:            StateWandering,
		visitedShelves:   make(map[[2]float64]Map.Shelf),
	}
	agent.agentBehavior = &ClientAgentBehavior{ag: agent}

	return agent
}

func (ag *ClientAgent) State() AgentState {
	return ag.state
}

func (ag *ClientAgent) NextAction() ActionType {
	return ag.nextAction
}

func (ag *ClientAgent) ShoppingList() []Map.Item {
	return ag.shoppingList
}

func (ag *ClientAgent) findWantedItemLocation() (float64, float64, bool) {
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

func (ag *ClientAgent) chooseItemToPickFromTargetedShelf(env *Environment) {
	ag.itemsToPick = []Map.Item{}
	shelfCoords := [2]float64{ag.interactTargetX, ag.interactTargetY}

	env.Mutex.RLock()
	defer env.Mutex.RUnlock()

	shelf, ok := env.Map.ShelfData[shelfCoords]

	if !ok {
		logger.Warnf("Shelf (%.2f %.2f) does not exist in the current environment", shelfCoords[0], shelfCoords[1])
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

func (ag *ClientAgent) UpdateFOVRays(dx, dy float64, numRays int, env *Environment) {
	ax := ag.coordinate.X + constants.CENTER_OF_CELL
	ay := ag.coordinate.Y + constants.CENTER_OF_CELL
	v := ag.visionManager
	v.RaysEndPoints = make([]utils.Vec2, numRays)

	fovAngle := constants.ANGLE_VISION * math.Pi / 180.0
	halfFOV := fovAngle / 2

	// Angle direction
	baseAngle := math.Atan2(dy, dx)
	for i := 0; i < numRays; i++ {
		angle := baseAngle - halfFOV + (float64(i)/float64(numRays-1))*fovAngle
		rayX, rayY := ax, ay
		step := 0.1

		for d := 0.0; d < v.visionDistance; d += step {
			rayX = math.Ceil(ax + math.Cos(angle)*d)
			rayY = math.Ceil(ay + math.Sin(angle)*d)

			if env.IsObstacleAt(rayX, rayY) {
				break
			}

			if env.IsShelfAt(rayX, rayY) {
				coords := [2]float64{rayX, rayY}
				if shelf, ok := env.Map.ShelfData[coords]; ok {
					ag.visitedShelves[coords] = shelf
				}
			}
		}

		v.RaysEndPoints[i] = utils.Vec2{X: rayX, Y: rayY}
	}
}

func (ag *ClientAgent) DetectShelvesInFOV(env *Environment) {
	// Keep in comment in case we want to know the items and shelves percepted by the agent
	// shelves := []Map.Shelf{}

	for coords := range env.Map.ShelfData {

		cx := float64(coords[0]) + constants.CENTER_OF_CELL
		cy := float64(coords[1]) + constants.CENTER_OF_CELL

		p := utils.Vec2{X: cx, Y: cy}
		v := ag.visionManager
		if v.areCoordinatesIntersectingFOV(p) {
			ag.visitedShelves[coords] = env.Map.ShelfData[coords]
			// shelves = append(shelves, shelf)
		}
	}
	// for _, s := range shelves {
	//     for _, i := range s.Items {
	// 			fmt.Println(i.Name)
	//     }
	// }
}

type ClientAgentBehavior struct {
	ag *ClientAgent
}

func (bh *ClientAgentBehavior) Percept() {
	ag := bh.ag
	// ag.UpdateFOVRays(ag.dx, ag.dy, 10, ag.env)
	ag.visionManager.UpdateFOV(ag.dx, ag.dy)
	ag.DetectShelvesInFOV(ag.env)
}

func (bh *ClientAgentBehavior) Deliberate() {
	ag := bh.ag
	ag.stuckDetector.DetectAndResolve()

	switch ag.state {
	case StateWandering:

		// Agent has finished shopping (either if he has collected all his shopping list, or if he couldnt find more items)
		if (len(ag.visitedShelves) >= len(ag.env.Map.ShelfData)) || (len(ag.GetMissingItems()) == 0) {
			if len(ag.cart) == 0 {
				ag.state = StateMovingToExit
				ag.nextAction = ActionWait
				break
			}

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
				logger.Warnf("No path found to a location near this destination (%.2f %.2f) ", shelfX, shelfY)
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
		if (len(ag.GetMissingItems()) != 0) && (len(ag.visitedShelves) < len(ag.env.Map.ShelfData)) { // si vol d'items ? (� impl�menter plus tard)
			ag.state = StateWandering
			ag.nextAction = ActionWait
			break
		}
		ag.processPathMovement()
		if !ag.hasDestination {
			ag.state = StateCheckingOut
			ag.nextAction = ActionWait
		}
	case StateCheckingOut: // fait pas grand chose, peut-�tre voir pour refacto
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

func (bh *ClientAgentBehavior) Act() {
	ag := bh.ag
	switch ag.nextAction {
	case ActionMove:
		ag.moveChan <- MoveRequest{Agt: ag, ResponseChannel: ag.moveChanResponse}
		<-ag.moveChanResponse

	case ActionPick:
		targetShelf := [2]float64{ag.interactTargetX, ag.interactTargetY}
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
				logger.Warnf("Bad request, PickRequest denied for (%s) at shelf (%.2f, %.2f)", item.Name, targetShelf[0], targetShelf[1])
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
		ag.exitChan <- ExitRequest{
			Agt:             ag,
			ResponseChannel: ag.exitChanResponse,
		}
	}
}
