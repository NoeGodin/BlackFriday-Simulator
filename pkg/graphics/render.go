package Graphics

import (
	Map "AI30_-_BlackFriday/pkg/map"
	Simulation "AI30_-_BlackFriday/pkg/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// TODO : on pourrait faire plus petit, genre 10 car j'ai peur que ça fasse paté si la map devient grande, on verra
// Si on change il faudra aussi changer les assets car sinon ils vont être déformée vu que reduit

func (g *Game) mapToDrawCoords(mapX float64, mapY float64, offsetX, offsetY int) (float64, float64) {
	drawX := mapX*float64(CELL_SIZE) + float64(g.CameraX+offsetX)
	drawY := mapY*float64(CELL_SIZE) + float64(g.CameraY+offsetY)
	return drawX, drawY
}

// Draw image at the right position with the right scale
func drawImageAt(screen *ebiten.Image, img *ebiten.Image, x, y float64) {
	if img == nil {
		return
	}
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(float64(CELL_SIZE)/float64(img.Bounds().Dx()),
		float64(CELL_SIZE)/float64(img.Bounds().Dy()))
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, options)
}

func (g *Game) Update() error {
	g.HandleInput()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.DrawMap(screen)
	g.DrawPaths(screen)
	g.DrawAgents(screen)
	g.DrawHUD(screen)
}

func (g *Game) DrawMap(screen *ebiten.Image) {
	MARGIN := 20
	offsetX := MARGIN
	offsetY := MARGIN

	// DRAW GRID LINES USEFUL FOR DEBUGGING
	/*for i := range g.Map.Width {
		x := float32(i*CELL_SIZE - g.CameraX + offsetX)
		vector.StrokeLine(
			screen,
			x,
			float32(offsetY),
			x,
			float32(g.Map.Height*CELL_SIZE+offsetY),
			1,
			color.Gray{Y: 240},
			false,
		)
	}

	for i := range g.Map.Height {
		y := float32(i*CELL_SIZE - g.CameraY + offsetY)
		vector.StrokeLine(
			screen,
			float32(offsetX),
			y,
			float32(g.Map.Width*CELL_SIZE+offsetX),
			y,
			1,
			color.Gray{Y: 240},
			false,
		)
	}*/

	//DRAW THE GROUND FIRST
	envMap := g.Simulation.Env.Map
	for y := range envMap.Height {
		for x := range envMap.Width {
			drawX, drawY := g.mapToDrawCoords(float64(x), float64(y), offsetX, offsetY)
			drawImageAt(screen, groundImg, drawX, drawY)
		}
	}

	// DRAW EVERY OTHER ELEMENTS, SAME LOGIC IN A WAY BUT GOOD PRACTICE TO KEEP BOTH
	for y := range envMap.Height {
		for x := range envMap.Width {
			elementType := envMap.GetElementType(x, y)
			if elementType == Map.VOID {
				continue
			}

			drawX, drawY := g.mapToDrawCoords(float64(x), float64(y), offsetX, offsetY)

			//If image not exist, will not render !
			switch elementType {
			case Map.WALL:
				drawImageAt(screen, wallImg, drawX, drawY)
			case Map.SHELF:
				drawImageAt(screen, itemImg, drawX, drawY)
			case Map.DOOR:
				drawImageAt(screen, doorImg, drawX, drawY)
			case Map.CHECKOUT:
				drawImageAt(screen, checkoutImg, drawX, drawY)
			}
		}
	}
}

func (g *Game) DrawAgents(screen *ebiten.Image) {
	MARGIN := 20
	offsetX := MARGIN
	offsetY := MARGIN
	for _, agt := range g.Simulation.Agents() {
		agtCoords := agt.Coordinate()
		drawX, drawY := g.mapToDrawCoords(agtCoords.X, agtCoords.Y, offsetX, offsetY)
		drawImageAt(screen, g.AgentAnimator.AnimationFrame(agt), drawX, drawY)
	}
}

func (g *Game) DrawPaths(screen *ebiten.Image) {
	MARGIN := 20
	offsetX := MARGIN
	offsetY := MARGIN

	for _, agt := range g.Simulation.Agents() {
		if clientAgent, ok := agt.(*Simulation.ClientAgent); ok {
			g.drawAgentPath(screen, clientAgent, offsetX, offsetY)
		}
	}
}

func (g *Game) drawAgentPath(screen *ebiten.Image, agent *Simulation.ClientAgent, offsetX, offsetY int) {
	path := agent.GetCurrentPath()
	if path == nil {
		return
	}

	waypoints := path.GetWaypoints()
	if len(waypoints) == 0 {
		return
	}

	agentCoord := agent.Coordinate()
	agentX := agentCoord.X*float64(CELL_SIZE) + float64(g.CameraX+offsetX) + float64(CELL_SIZE)/2
	agentY := agentCoord.Y*float64(CELL_SIZE) + float64(g.CameraY+offsetY) + float64(CELL_SIZE)/2

	for i, waypoint := range waypoints {
		waypointX := float64(waypoint.X)*float64(CELL_SIZE) + float64(g.CameraX+offsetX) + float64(CELL_SIZE)/2
		waypointY := float64(waypoint.Y)*float64(CELL_SIZE) + float64(g.CameraY+offsetY) + float64(CELL_SIZE)/2

		// Waypoints : black circle
		vector.FillCircle(screen, float32(waypointX), float32(waypointY), 3, color.Black, false)

		// Lines between waypoints
		if i > 0 {
			prevWaypoint := waypoints[i-1]
			prevX := float64(prevWaypoint.X)*float64(CELL_SIZE) + float64(g.CameraX+offsetX) + float64(CELL_SIZE)/2
			prevY := float64(prevWaypoint.Y)*float64(CELL_SIZE) + float64(g.CameraY+offsetY) + float64(CELL_SIZE)/2

			vector.StrokeLine(screen, float32(prevX), float32(prevY), float32(waypointX), float32(waypointY), 2, color.RGBA{0, 100, 255, 255}, false)
		}
	}

	// Waypoints : green circle
	target := path.GetTarget()
	targetX := float64(target.X)*float64(CELL_SIZE) + float64(g.CameraX+offsetX) + float64(CELL_SIZE)/2
	targetY := float64(target.Y)*float64(CELL_SIZE) + float64(g.CameraY+offsetY) + float64(CELL_SIZE)/2
	vector.FillCircle(screen, float32(targetX), float32(targetY), 5, color.RGBA{0, 255, 0, 255}, false)

	if len(waypoints) > 0 {
		nextWaypoint := waypoints[0]
		nextX := float64(nextWaypoint.X)*float64(CELL_SIZE) + float64(g.CameraX+offsetX) + float64(CELL_SIZE)/2
		nextY := float64(nextWaypoint.Y)*float64(CELL_SIZE) + float64(g.CameraY+offsetY) + float64(CELL_SIZE)/2

		// Redline between agent and current waypoint
		vector.StrokeLine(screen, float32(agentX), float32(agentY), float32(nextX), float32(nextY), 3, color.RGBA{255, 0, 0, 255}, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	envMap := g.Simulation.Env.Map
	// adapt the window size
	mapWidth := envMap.Width * CELL_SIZE
	mapHeight := envMap.Height * CELL_SIZE

	// better with this :)
	MARGIN := 20
	return mapWidth + MARGIN*2, mapHeight + MARGIN*2
}
