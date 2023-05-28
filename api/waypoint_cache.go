package api

import (
	"context"
	"spacetraders_engine/ext"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type WaypointCache struct {
	client    *sdk.APIClient
	context   context.Context
	waypoints map[string][]ext.Waypoint
}

func NewWaypointCache(client *sdk.APIClient, ctx context.Context) *WaypointCache {
	return &WaypointCache{
		client:    client,
		context:   ctx,
		waypoints: make(map[string][]ext.Waypoint),
	}
}

func (wpc *WaypointCache) WaypointsInSystem(system string) ([]ext.Waypoint, error) {
	if v, ok := wpc.waypoints[system]; ok {
		return v, nil
	}

	sys := NewSystem(wpc.client, system)
	wp, err := sys.GetWaypoints(wpc.context)

	if err != nil {
		return []ext.Waypoint{}, err
	}

	wpc.waypoints[system] = ext.NewWaypointList(wp)
	return wpc.waypoints[system], nil
}

func (wpc *WaypointCache) Waypoint(symbol string) (ext.Waypoint, error) {
	loc := ext.NewLocation(symbol)

	if !loc.HasSystem() {
		return ext.Waypoint{}, NewInvalidWaypointError(symbol)
	}

	system := loc.System
	waypoints, err := wpc.WaypointsInSystem(system)

	if err != nil {
		return ext.Waypoint{}, err
	}

	for _, v := range waypoints {
		if v.Symbol == symbol {
			return v, nil
		}
	}

	return ext.Waypoint{}, NewInvalidWaypointError(symbol)
}
