package pathfinding

import "AI30_-_BlackFriday/pkg/utils"

// GetNextWaypoint returns the next waypoint to reach
func (p *Path) GetNextWaypoint() (utils.Coordinate, bool) {
	if len(p.waypoints) == 0 {
		return utils.Coordinate{}, false
	}
	return p.waypoints[0], true
}

func (p *Path) IsComplete() bool {
	return len(p.waypoints) == 0
}

func (p *Path) GetWaypoints() []utils.Coordinate {
	return p.waypoints
}

func (p *Path) GetTarget() utils.Coordinate {
	return p.target
}

func (p *Path) RemoveFirstWaypoint() {
	if len(p.waypoints) > 0 {
		p.waypoints = p.waypoints[1:]
	}
}
