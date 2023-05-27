package game

import (
	"context"
	"spacetraders_engine/api"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type MarketCache struct {
	client    *sdk.APIClient
	context   context.Context
	markets   map[string]sdk.Market
	waypoints *WaypointCache
}

func NewMarketCache(client *sdk.APIClient, ctx context.Context, wp *WaypointCache) *MarketCache {
	return &MarketCache{
		client:    client,
		context:   ctx,
		markets:   make(map[string]sdk.Market),
		waypoints: wp,
	}
}

func (mc *MarketCache) MarketForSymbol(symbol string) (sdk.Market, error) {
	if mkt, ok := mc.markets[symbol]; ok {
		return mkt, nil
	}

	wp, err := mc.waypoints.Waypoint(symbol)

	if err != nil {
		return sdk.Market{}, err
	}

	a, err := api.NewWaypoint(mc.client, wp.Symbol)

	if err != nil {
		return sdk.Market{}, err
	}

	return a.Market(mc.context)
}

func (mc *MarketCache) MarketsInSystem(system string) ([]sdk.Market, error) {
	waypoints := make([]sdk.Waypoint, 0)
	result := make([]sdk.Market, 0)
	wp, err := mc.waypoints.WaypointsInSystem(system)

	if err != nil {
		return []sdk.Market{}, err
	}

	for _, v := range wp {
		g := NewWaypoint(v)
		if g.HasMarket() {
			waypoints = append(waypoints, v)
		}
	}

	for _, v := range waypoints {
		mkt, err := mc.MarketForSymbol(v.Symbol)

		if err != nil {
			return []sdk.Market{}, err
		}

		result = append(result, mkt)
	}

	return result, nil
}
