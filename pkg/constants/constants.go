package constants

import (
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type ElementType string

const (
	NB_AGENTS = 50

	// Map element types
	WALL     ElementType = "W"
	SHELF    ElementType = "shelf"
	CHECKOUT ElementType = "C"
	DOOR     ElementType = "D"
	VOID     ElementType = " "

	// Graphics constants
	CELL_SIZE            = 32
	FRAME_DURATION       = 10
	FRAME_COUNT          = 4
	DIRECTIONS           = 4
	MARGIN               = 20
	AGENT_STATE_DURATION = time.Second * 2

	// StuckDistanceThreshold minimum distance to be considered moving (calculated as agent speed / 2)
	// StuckDistanceThreshold = BASE_AGENT_SPEED / 2
	// StuckCounterThreshold number of frames to consider an agent stuck (~1 second at 30 FPS)
	StuckCounterThreshold = 30
	// WaypointReachedThreshold distance threshold to consider a waypoint reached
	WaypointReachedThreshold = 0.6
	AgentToAgentHitbox       = 0.15
	AgentToEnvironmentHitbox = 0.4

	HUD_POS_X float64 = 10.0
	HUD_POS_Y float64 = 10.0

	DELTA_TIME                        = 1.0 / 60.0
	AGENT_SEARCH_RADIUS               = 5.0
	AGENT_AGGRESSIVENESS_TRESHOLD     = 1.0
	AGENT_AGGRESSIVENESS_INCREASEMENT = 0.33
	BASE_AGENT_AGGRESSIVENESS         = 0.33
	AGENT_STEAL_RANGE                 = 1.0

	//SFC
	SOCIAL_STRENGTH   = 1.0
	WALL_RESISTANCE   = 10.0
	AGT_RANGE         = 0.2
	AGT_STRENGTH      = 10.0
	AGT_RADIUS        = 0.3
	FRICTION_COEF     = 10.0
	RELAXATION_FACTOR = 30.0
	SPEED_MULTIPLIER  = 1.1

	VISION_DISTANCE = 10
	VISION_HEIGHT   = 6
	ANGLE_VISION    = 90.0 // For Raycast FOV

	CENTER_OF_CELL = 0.5

	AGENT_SPAWN_INTERVAL_MS = 200
	SPAWN_OFFSET_FROM_DOOR  = 1 // Avoid spawning behind a door

	MAX_TIC_DURATION   = 150
	MIN_TIC_DURATION   = 10
	MAX_CLIENTS_NUMBER = 2000
	MAX_GUARD_NUMBER   = 100

	// file system
	FONT_PATH = "assets/fonts/Monaco.ttf"
	MAPS_PATH = "maps/store"
)

// loaded from .env using autoload
var (
	NUMBER_OF_CLIENTS          = envInt("NUMBER_OF_CLIENTS", 75)
	NUMBER_OF_GUARDS           = envInt("NUMBER_OF_GUARDS", 2)
	BASE_AGENT_SPEED           = envFloat("BASE_AGENT_SPEED", 0.2)
	AGENT_MAX_SHOPPING_LIST    = envInt("AGENT_MAX_SHOPPING_LIST", 4)
	TIC_DURATION               = envInt("TIC_DURATION", 100)
	SALES_EXPORT_INTERVAL      = time.Duration(envInt("SALES_EXPORT_INTERVAL_SECONDS", 30)) * time.Second
	MAX_QUANTITY_PER_ITEM      = envInt("MAX_QUANTITY_PER_ITEM", 5)
	AGRESSIVE_AGENT_PROPORTION = envFloat("AGRESSIVE_AGENT_PROPORTION", 0)
)

func envInt(key string, def int) int {
	if v, _ := strconv.Atoi(os.Getenv(key)); v >= 0 {
		return v
	}
	return def
}

func envFloat(key string, def float64) float64 {
	if v, _ := strconv.ParseFloat(os.Getenv(key), 64); v > 0 {
		return v
	}
	return def
}

// StuckDistanceThreshold returns the minimum distance to be considered moving (agent speed / 2)
func StuckDistanceThreshold() float64 {
	return BASE_AGENT_SPEED / 2
}

// MovementDirections all possible movement directions
var MovementDirections = [][2]float64{
	{0, -1}, {1, 0}, {0, 1}, {-1, 0}, // Nord, Est, Sud, Ouest
	{1, -1}, {1, 1}, {-1, 1}, {-1, -1}, // Diagonales
}
