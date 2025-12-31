package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	Hud "AI30_-_BlackFriday/pkg/hud"
	"AI30_-_BlackFriday/pkg/logger"
	Map "AI30_-_BlackFriday/pkg/map"
	Simulation "AI30_-_BlackFriday/pkg/simulation"
	"fmt"
	"math/rand"

	"github.com/ebitenui/ebitenui"
)

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	AgentAnimator             *AgentAnimator
	ShelfAnimator             *ShelfAnimator
	Simulation                *Simulation.Simulation
	Hud                       Hud.HUD
	UI                        *ebitenui.UI
	startMenu                 *Hud.StartMenu
	inMenu                    bool
	guardsFOV                 bool
	ticDuration               int
}

func NewGame(screenWidth, screenHeight int) *Game {
	game := &Game{
		ScreenWidth:   screenWidth,
		ScreenHeight:  screenHeight,
		CameraX:       0,
		CameraY:       0,
		AgentAnimator: NewAgentAnimator(),
		ShelfAnimator: NewShelfAnimator(),
		Hud:           *Hud.NewHud(),
		inMenu:        true,
		guardsFOV:     false,
		ticDuration:   constants.TIC_DURATION,
	}
	game.startMenu = Hud.NewStartMenu(func(p Hud.Properties) {
		game.startSimulation(p)
	})

	game.UI = createUI(game)

	return game
}

func (game *Game) SetTicDuration(duration int) {
	game.ticDuration = duration
	game.Simulation.SetTicDuration(duration)

}
func (game *Game) startSimulation(p Hud.Properties) {
	mapPath := fmt.Sprintf("%s/%s", constants.MAPS_PATH, p.Filename)
	mapData, err := Map.LoadMapFromFile(mapPath)
	if err != nil {
		logger.Errorf("Error loading map: %s", err.Error())
		panic("Error loading map: " + err.Error())
	}

	mapName := extractMapName(mapPath)
	shoppingListsPath := fmt.Sprintf("%s/shopping_lists.json", constants.MAPS_PATH)
	simu := Simulation.NewSimulation(0, game.ticDuration, mapData, constants.DELTA_TIME, constants.AGENT_SEARCH_RADIUS, mapName, shoppingListsPath)
	logger.Info("Adding agents...")
	logger.Infof("Creating %d client from config", p.ClientNumber)
	for i := 1; i <= p.ClientNumber; i++ {
		agentID := fmt.Sprintf("agent%d", i)
		aggressiveness := 0.0
		if rand.Float64() > (1.0 - p.AgentAggressiveness) {
			aggressiveness = 1.0
		}
		simu.AddClient(agentID, aggressiveness)
	}
	for i := 1; i <= p.GuardNumber; i++ {
		agentID := fmt.Sprintf("guard%d", i)
		simu.AddGuard(agentID)
	}

	game.Simulation = simu
	game.Simulation.Run()
	game.inMenu = false

}

func (game *Game) TogglePause() {
	if game.Simulation == nil {
		return
	}
	game.Simulation.TogglePause()
}
