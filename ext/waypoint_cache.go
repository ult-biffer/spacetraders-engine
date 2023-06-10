package ext

import (
	"fmt"

	"github.com/ult-biffer/spacetraders_sdk/api"
	"github.com/ult-biffer/spacetraders_sdk/util"
)

type WaypointCache struct {
	Waypoints map[string][]*Waypoint `json:"waypoints"`
}

func NewWaypointCache() *WaypointCache {
	return &WaypointCache{
		Waypoints: make(map[string][]*Waypoint),
	}
}

func (c *WaypointCache) WaypointsInSystem(system string) ([]*Waypoint, error) {
	if v, ok := c.Waypoints[system]; ok {
		return v, nil
	}

	v, err := api.WaypointsInSystem(system)

	if err != nil {
		return nil, err
	}

	c.Waypoints[system] = NewWaypointList(v)
	return c.Waypoints[system], nil
}

func (c *WaypointCache) Waypoint(symbol string) (*Waypoint, error) {
	loc := util.NewLocation(symbol)

	if !loc.IsWaypoint() {
		return nil, fmt.Errorf("%s is not a valid waypoint", symbol)
	}

	waypoints, err := c.WaypointsInSystem(loc.System)

	if err != nil {
		return nil, err
	}

	for _, v := range waypoints {
		if v.Symbol == symbol {
			return v, nil
		}
	}

	return nil, fmt.Errorf("could not find waypoint %s", symbol)
}
