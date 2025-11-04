package Graphics

import Map "AI30_-_BlackFriday/pkg/map"

type Game struct {
	ScreenWidth, ScreenHeight int
	CameraX, CameraY          int
	Map                       Map.Map
}