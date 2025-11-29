package pathfinding

import "AI30_-_BlackFriday/pkg/utils"

// GetNextWaypoint returns the next waypoint to reach
func (p *Path) GetNextWaypoint() (utils.Vec2, bool) {
	if len(p.waypoints) == 0 {
		return utils.Vec2{}, false
	}
	return p.waypoints[0], true
}

func (p *Path) IsComplete() bool {
	return len(p.waypoints) == 0
}

func (p *Path) GetWaypoints() []utils.Vec2 {
	return p.waypoints
}

func (p *Path) GetTarget() utils.Vec2 {
	return p.target
}

func (p *Path) RemoveFirstWaypoint() {
	if len(p.waypoints) > 0 {
		p.waypoints = p.waypoints[1:]
	}
}
