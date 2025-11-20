package pathfinding

import "AI30_-_BlackFriday/pkg/utils"

// GetNextWaypoint returns the next waypoint to reach
func (p *Path) GetNextWaypoint() (utils.IntCoordinate, bool) {
	if len(p.Waypoints) == 0 {
		return utils.IntCoordinate{}, false
	}
	return p.Waypoints[0], true
}

func (p *Path) IsComplete() bool {
	return len(p.Waypoints) == 0
}

func (p *Path) GetWaypoints() []utils.IntCoordinate {
	return p.Waypoints
}

func (p *Path) GetTarget() utils.IntCoordinate {
	return p.Target
}

func (p *Path) RemoveFirstWaypoint() {
	if len(p.Waypoints) > 0 {
		p.Waypoints = p.Waypoints[1:]
	}
}
