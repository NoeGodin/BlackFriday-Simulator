package Graphics

import (
	"AI30_-_BlackFriday/pkg/constants"
	Hud "AI30_-_BlackFriday/pkg/hud"
	Simulation "AI30_-_BlackFriday/pkg/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TODO : on pourrait faire plus petit, genre 10 car j'ai peur que ça fasse paté si la map devient grande, on verra
// Si on change il faudra aussi changer les assets car sinon ils vont être déformée vu que reduit

func (g *Game) mapToDrawCoords(mapX float64, mapY float64, offsetX, offsetY int) (float64, float64) {
	drawX := mapX*float64(constants.CELL_SIZE) + float64(g.CameraX+offsetX)
	drawY := mapY*float64(constants.CELL_SIZE) + float64(g.CameraY+offsetY)
	return drawX, drawY
}

// Draw image at the right position with the right scale
func drawImageAt(screen *ebiten.Image, img *ebiten.Image, x, y float64, colorScale *ebiten.ColorScale) {
	if img == nil {
		return
	}
	options := &ebiten.DrawImageOptions{}
	if colorScale != nil {
		options.ColorScale = *colorScale
	}
	options.GeoM.Scale(float64(constants.CELL_SIZE)/float64(img.Bounds().Dx()),
		float64(constants.CELL_SIZE)/float64(img.Bounds().Dy()))
	options.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, options)
}

func (g *Game) Update() error {
	g.HandleInput()
	g.Hud.Update()
	g.UI.Update()
	if g.Hud.GetSelectedAgent() != nil {
		g.ensureValidHUDAgent()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.DrawMap(screen)
	if g.Hud.DisplayAgentPaths {
		g.DrawPaths(screen)
	}
	if g.Hud.GetSelectedAgent() != nil {
		g.DrawPath(screen)
	} 
	g.DrawAgents(screen)
	g.DrawTexture(screen)
	g.DrawHUD(screen)
	g.UI.Draw(screen)
}

func (g *Game) DrawMap(screen *ebiten.Image) {
	MARGIN := constants.MARGIN
	offsetX := MARGIN
	offsetY := MARGIN

	// DRAW GRID LINES USEFUL FOR DEBUGGING
	/*for i := range g.Map.Width {
		x := float32(i*constants.CELL_SIZE - g.CameraX + offsetX)
		vector.StrokeLine(
			screen,
			x,
			float32(offsetY),
			x,
			float32(g.Map.Height*constants.CELL_SIZE+offsetY),
			1,
			color.Gray{Y: 240},
			false,
		)
	}

	for i := range g.Map.Height {
		y := float32(i*constants.CELL_SIZE - g.CameraY + offsetY)
		vector.StrokeLine(
			screen,
			float32(offsetX),
			y,
			float32(g.Map.Width*constants.CELL_SIZE+offsetX),
			y,
			1,
			color.Gray{Y: 240},
			false,
		)
	}*/

	//DRAW THE GROUND FIRST
	envMap := g.Simulation.Env.Map
	for y := 0.0; y < float64(envMap.Height); y++ {
		for x := 0.0; x < float64(envMap.Width); x++ {
			drawX, drawY := g.mapToDrawCoords(x, y, offsetX, offsetY)
			drawImageAt(screen, groundImg, drawX, drawY, nil)
		}
	}

	// DRAW EVERY OTHER ELEMENTS, SAME LOGIC IN A WAY BUT GOOD PRACTICE TO KEEP BOTH
	for y := 0.0; y < float64(envMap.Height); y++ {
		for x := 0.0; x < float64(envMap.Width); x++ {
			elementType := envMap.GetElementType(x, y)
			if elementType == constants.VOID {
				continue
			}

			drawX, drawY := g.mapToDrawCoords(x, y, offsetX, offsetY)

			//If image not exist, will not render !
			switch elementType {
			case constants.WALL:
				if envMap.GetElementType(x, y+1) != constants.WALL {
					drawImageAt(screen, wallImg, drawX, drawY, nil)
				} else {
					drawImageAt(screen, wallCeiling, drawX, drawY, nil)
				}
			case constants.SHELF:
				if shelf, ok := envMap.GetShelfData(x, y); ok {
					img := g.ShelfAnimator.AnimationFrame([2]float64{x, y}, &shelf)
					drawImageAt(screen, img, drawX, drawY, nil)
				} else {
					drawImageAt(screen, itemImg, drawX, drawY, nil)
				}
			case constants.DOOR:
				drawImageAt(screen, doorImg, drawX, drawY, nil)
			case constants.CHECKOUT:
				drawImageAt(screen, checkoutImg, drawX, drawY, nil)
			}
		}
	}

}

func (g *Game) DrawAgents(screen *ebiten.Image) {
	MARGIN := constants.MARGIN
	offsetX := MARGIN
	offsetY := MARGIN
	for _, agt := range g.Simulation.Agents() {
		if !agt.HasSpawned() {
			continue
		}
		agtCoords := agt.Coordinate()
		drawX, drawY := g.mapToDrawCoords(agtCoords.X, agtCoords.Y, offsetX, offsetY)
		colorScale := g.AgentAnimator.getColorScale(agt)
		drawImageAt(screen, g.AgentAnimator.AnimationFrame(agt), drawX, drawY, colorScale)
		drawImageAt(screen, g.AgentAnimator.GetEmotion(agt), drawX, drawY, nil)
	}
}

func (g *Game) DrawPaths(screen *ebiten.Image) {
	MARGIN := constants.MARGIN
	offsetX := MARGIN
	offsetY := MARGIN

	for _, agt := range g.Simulation.Agents() {
		g.drawAgentPath(screen, agt, offsetX, offsetY)
	}
}

func (g *Game) DrawPath(screen *ebiten.Image) {
	MARGIN := constants.MARGIN
	offsetX := MARGIN
	offsetY := MARGIN

	g.drawAgentPath(screen, g.Hud.GetSelectedAgent(), offsetX, offsetY)
}

func (g *Game) drawAgentPath(screen *ebiten.Image, agent Simulation.Agent, offsetX, offsetY int) {
	path := agent.GetCurrentPath()
	if path == nil {
		return
	}

	waypoints := path.GetWaypoints()
	if len(waypoints) == 0 {
		return
	}

	agentCoord := agent.Coordinate()
	agentX := agentCoord.X*float64(constants.CELL_SIZE) + float64(g.CameraX+offsetX) + float64(constants.CELL_SIZE)/2
	agentY := agentCoord.Y*float64(constants.CELL_SIZE) + float64(g.CameraY+offsetY) + float64(constants.CELL_SIZE)/2

	for i, waypoint := range waypoints {
		waypointX := float64(waypoint.X)*float64(constants.CELL_SIZE) + float64(g.CameraX+offsetX) + float64(constants.CELL_SIZE)/2
		waypointY := float64(waypoint.Y)*float64(constants.CELL_SIZE) + float64(g.CameraY+offsetY) + float64(constants.CELL_SIZE)/2

		// Waypoints : black circle
		vector.FillCircle(screen, float32(waypointX), float32(waypointY), 3, color.Black, false)

		// Lines between waypoints
		if i > 0 {
			prevWaypoint := waypoints[i-1]
			prevX := float64(prevWaypoint.X)*float64(constants.CELL_SIZE) + float64(g.CameraX+offsetX) + float64(constants.CELL_SIZE)/2
			prevY := float64(prevWaypoint.Y)*float64(constants.CELL_SIZE) + float64(g.CameraY+offsetY) + float64(constants.CELL_SIZE)/2

			vector.StrokeLine(screen, float32(prevX), float32(prevY), float32(waypointX), float32(waypointY), 2, color.RGBA{0, 100, 255, 255}, false)
		}
	}

	// Waypoints : green circle
	target := path.GetTarget()
	targetX := float64(target.X)*float64(constants.CELL_SIZE) + float64(g.CameraX+offsetX) + float64(constants.CELL_SIZE)/2
	targetY := float64(target.Y)*float64(constants.CELL_SIZE) + float64(g.CameraY+offsetY) + float64(constants.CELL_SIZE)/2
	vector.FillCircle(screen, float32(targetX), float32(targetY), 5, color.RGBA{0, 255, 0, 255}, false)

	if len(waypoints) > 0 {
		nextWaypoint := waypoints[0]
		nextX := float64(nextWaypoint.X)*float64(constants.CELL_SIZE) + float64(g.CameraX+offsetX) + float64(constants.CELL_SIZE)/2
		nextY := float64(nextWaypoint.Y)*float64(constants.CELL_SIZE) + float64(g.CameraY+offsetY) + float64(constants.CELL_SIZE)/2

		// Redline between agent and current waypoint
		vector.StrokeLine(screen, float32(agentX), float32(agentY), float32(nextX), float32(nextY), 3, color.RGBA{255, 0, 0, 255}, false)
	}
}

func (g *Game) DrawHUD(screen *ebiten.Image) {
	if g.Hud.Hidden() {
		return
	}

	offsetX := constants.MARGIN
	offsetY := constants.MARGIN

	if g.Hud.HudBg == nil {
		return
	}

	if agt := g.Hud.GetSelectedAgent(); agt != nil {
		// Draw Raycasts
		// ax := agt.Coordinate().X + constants.CENTER_OF_CELL
		// ay := agt.Coordinate().Y + constants.CENTER_OF_CELL
		// axDraw, ayDraw := g.mapToDrawCoords(ax, ay, offsetX, offsetY)

		// for _, end := range agt.VisionManager().RaysEndPoints {
		// 	exDraw, eyDraw := g.mapToDrawCoords(end.X, end.Y, offsetX, offsetY)
		// 	vector.StrokeLine(screen, float32(axDraw), float32(ayDraw),
		// 		float32(exDraw), float32(eyDraw),
		// 		1, color.RGBA{0, 255, 255, 128}, false)
		// }

		// Draw rectangle
		p1 := agt.VisionManager().P1
		p2 := agt.VisionManager().P2
		p3 := agt.VisionManager().P3
		p4 := agt.VisionManager().P4

		x1, y1 := g.mapToDrawCoords(p1.X, p1.Y, offsetX, offsetY)
		x2, y2 := g.mapToDrawCoords(p2.X, p2.Y, offsetX, offsetY)
		x3, y3 := g.mapToDrawCoords(p3.X, p3.Y, offsetX, offsetY)
		x4, y4 := g.mapToDrawCoords(p4.X, p4.Y, offsetX, offsetY)

		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 1, color.RGBA{0, 0, 0, 255}, false)
		vector.StrokeLine(screen, float32(x2), float32(y2), float32(x3), float32(y3), 1, color.RGBA{0, 0, 0, 255}, false)
		vector.StrokeLine(screen, float32(x3), float32(y3), float32(x4), float32(y4), 1, color.RGBA{0, 0, 0, 255}, false)
		vector.StrokeLine(screen, float32(x4), float32(y4), float32(x1), float32(y1), 1, color.RGBA{0, 0, 0, 255}, false)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(constants.HUD_POS_X, constants.HUD_POS_Y)
	screen.DrawImage(g.Hud.HudBg, op)

	y := int(constants.HUD_POS_Y) + g.Hud.PaddingY + Hud.FONT.Metrics().Height.Ceil()
	for _, line := range g.Hud.Lines {
		text.Draw(screen, line, Hud.FONT, int(constants.HUD_POS_X)+g.Hud.PaddingX, y, color.White)
		y += Hud.FONT.Metrics().Height.Ceil()
	}

	targetX, targetY := g.mapToDrawCoords(g.Hud.TargetPositionX, g.Hud.TargetPositionY, offsetX, offsetY)
	drawImageAt(screen, targetImg, targetX, targetY, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	envMap := g.Simulation.Env.Map
	// adapt the window size
	mapWidth := envMap.Width * constants.CELL_SIZE
	mapHeight := envMap.Height * constants.CELL_SIZE

	// better with this :)
	MARGIN := constants.MARGIN
	return mapWidth + MARGIN*2, mapHeight + MARGIN*2
}
