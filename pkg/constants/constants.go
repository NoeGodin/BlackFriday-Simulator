package constants

import (
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type ElementType string

const (
	NB_AGENTS 		= 50

	// Map element types
	WALL     ElementType = "W"
	SHELF    ElementType = "shelf"
	CHECKOUT ElementType = "C"
	DOOR     ElementType = "D"
	VOID     ElementType = " "

	// Graphics constants
	CELL_SIZE      = 32
	FRAME_DURATION = 10
	FRAME_COUNT    = 4
	DIRECTIONS     = 4
	MARGIN         = 20

	// StuckDistanceThreshold minimum distance to be considered moving
	StuckDistanceThreshold = 0.1
	// StuckCounterThreshold number of frames to consider an agent stuck (~1 second at 30 FPS)
	StuckCounterThreshold = 30
	// WaypointReachedThreshold distance threshold to consider a waypoint reached
	WaypointReachedThreshold = 0.6
	AgentToAgentHitbox       = 0.15
	AgentToEnvironmentHitbox = 0.4

	HUD_POS_X float64 = 10.0
	HUD_POS_Y float64 = 10.0

	DELTA_TIME          = 1.0 / 60.0
	AGENT_SEARCH_RADIUS = 5.0

	//SFC
	SOCIAL_STRENGTH   = 1.0
	WALL_RESISTANCE   = 3.0
	AGT_RANGE         = 0.2
	AGT_STRENGTH      = 10.0
	AGT_RADIUS        = 0.3
	FRICTION_COEF     = 10.0
	RELAXATION_FACTOR = 30.0

	VISION_DISTANCE = 10
	VISION_HEIGHT   = 6
	ANGLE_VISION    = 90.0 // For Raycast FOV

	CENTER_OF_CELL = 0.5

	AGENT_SPAWN_INTERVAL_MS = 200
	SPAWN_OFFSET_FROM_DOOR = 0.5 // Avoid spawning behind a door
)

// loaded from .env using autoload
var (
	NUMBER_OF_AGENTS        = envInt("NUMBER_OF_AGENTS", 75)
	BASE_AGENT_SPEED        = envFloat("BASE_AGENT_SPEED", 0.2)
	AGENT_MAX_SHOPPING_LIST = envInt("AGENT_MAX_SHOPPING_LIST", 4)
	SPEED_MULTIPLIER        = envFloat("SFC_SPEED_MULTIPLIER", 1.1)
	SALES_EXPORT_INTERVAL   = time.Duration(envInt("SALES_EXPORT_INTERVAL_SECONDS", 30)) * time.Second
)

func envInt(key string, def int) int {
	if v, _ := strconv.Atoi(os.Getenv(key)); v > 0 {
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

// MovementDirections all possible movement directions
var MovementDirections = [][2]float64{
	{0, -1}, {1, 0}, {0, 1}, {-1, 0}, // Nord, Est, Sud, Ouest
	{1, -1}, {1, 1}, {-1, 1}, {-1, -1}, // Diagonales
}
