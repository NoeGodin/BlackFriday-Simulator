package pathfinding

import "AI30_-_BlackFriday/pkg/utils"

// GetNextWaypoint returns the next waypoint to reach
func (p *Path) GetNextWaypoint() (utils.IntCoordinate, bool) {
	if len(p.waypoints) == 0 {
		return utils.IntCoordinate{}, false
	}
	return p.waypoints[0], true
}

func (p *Path) IsComplete() bool {
	return len(p.waypoints) == 0
}

func (p *Path) GetWaypoints() []utils.IntCoordinate {
	return p.waypoints
}

func (p *Path) GetTarget() utils.IntCoordinate {
	return p.target
}

func (p *Path) RemoveFirstWaypoint() {
	if len(p.waypoints) > 0 {
		p.waypoints = p.waypoints[1:]
	}
}
