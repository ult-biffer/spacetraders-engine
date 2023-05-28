package api

import (
	"context"
	"spacetraders_engine/ext"
	"time"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type MarketCache struct {
	client    *sdk.APIClient
	context   context.Context
	markets   map[string]marketCacheEntry
	waypoints *WaypointCache
}

type marketCacheEntry struct {
	Market  sdk.Market
	SavedAt time.Time
}

func NewMarketCache(client *sdk.APIClient, ctx context.Context, wp *WaypointCache) *MarketCache {
	return &MarketCache{
		client:    client,
		context:   ctx,
		markets:   make(map[string]marketCacheEntry),
		waypoints: wp,
	}
}

func (mc *MarketCache) MarketForSymbol(symbol string) (sdk.Market, error) {
	if mkt, ok := mc.markets[symbol]; ok && !mkt.IsOld(nil) {
		return mkt.Market, nil
	}

	wp, err := mc.waypoints.Waypoint(symbol)

	if err != nil {
		return sdk.Market{}, err
	}

	a, err := NewWaypoint(mc.client, wp.Symbol)

	if err != nil {
		return sdk.Market{}, err
	}

	mkt, err := a.Market(mc.context)

	if err != nil {
		return sdk.Market{}, err
	}

	mc.markets[symbol] = marketCacheEntry{
		Market:  mkt,
		SavedAt: time.Now(),
	}

	return mkt, nil
}

func (mc *MarketCache) MarketsInSystem(system string) ([]sdk.Market, error) {
	waypoints := make([]ext.Waypoint, 0)
	result := make([]sdk.Market, 0)
	wp, err := mc.waypoints.WaypointsInSystem(system)

	if err != nil {
		return []sdk.Market{}, err
	}

	for _, v := range wp {
		if v.HasMarket() {
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

func (entry marketCacheEntry) IsOld(oldestAcceptable *time.Time) bool {
	if oldestAcceptable == nil {
		t := time.Now().Add(time.Duration(-1) * time.Hour)
		oldestAcceptable = &t
	}

	return entry.SavedAt.Before(*oldestAcceptable)
}
