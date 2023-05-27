package game

import (
	"context"
	"fmt"
	"spacetraders_engine/api"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type WaypointCache struct {
	client    *sdk.APIClient
	context   context.Context
	waypoints map[string][]sdk.Waypoint
}

func NewWaypointCache(client *sdk.APIClient, ctx context.Context) *WaypointCache {
	return &WaypointCache{
		client:    client,
		context:   ctx,
		waypoints: make(map[string][]sdk.Waypoint),
	}
}

func (wpc *WaypointCache) WaypointsInSystem(system string) ([]sdk.Waypoint, error) {
	if v, ok := wpc.waypoints[system]; ok {
		return v, nil
	}

	sys := api.NewSystem(wpc.client, system)
	wp, err := sys.GetWaypoints(wpc.context)

	if err != nil {
		return []sdk.Waypoint{}, err
	}

	wpc.waypoints[system] = wp
	return wp, nil
}

func (wpc *WaypointCache) Waypoint(symbol string) (sdk.Waypoint, error) {
	system, err := api.GetSystemFromWaypoint(symbol)

	if err != nil {
		return sdk.Waypoint{}, err
	}

	waypoints, err := wpc.WaypointsInSystem(system)

	if err != nil {
		return sdk.Waypoint{}, err
	}

	for _, v := range waypoints {
		if v.Symbol == symbol {
			return v, nil
		}
	}

	return sdk.Waypoint{}, fmt.Errorf("invalid waypoint %s", symbol)
}
