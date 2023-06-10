package ext

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ult-biffer/spacetraders_sdk/api"
	"github.com/ult-biffer/spacetraders_sdk/models"
)

const DEFAULT_MARKET_OLD = (time.Duration(-1) * time.Hour)

type MarketCache struct {
	Markets   map[string]marketCacheEntry `json:"markets"`
	waypoints *WaypointCache
}

type marketCacheEntry struct {
	models.Market
	SavedAt time.Time `json:"savedAt"`
}

func NewMarketCache(wp *WaypointCache) *MarketCache {
	return &MarketCache{
		Markets:   make(map[string]marketCacheEntry),
		waypoints: wp,
	}
}

func LoadMarketCache(body []byte, wp *WaypointCache) *MarketCache {
	var cache *MarketCache

	if err := json.Unmarshal(body, cache); err != nil {
		return NewMarketCache(wp)
	}

	cache.waypoints = wp
	return cache
}

func newMarketCacheEntry(mkt *models.Market) marketCacheEntry {
	return marketCacheEntry{
		Market:  *mkt,
		SavedAt: time.Now(),
	}
}

func (c *MarketCache) MarketForSymbol(waypoint string) (*models.Market, error) {
	if mkt, ok := c.Markets[waypoint]; ok && !mkt.IsOld() {
		return &mkt.Market, nil
	}

	wp, err := c.waypoints.Waypoint(waypoint)

	if err != nil {
		return nil, err
	}

	if !wp.HasMarket() {
		return nil, fmt.Errorf("waypoint %s has no market", waypoint)
	}

	mkt, err := api.GetMarket(waypoint)

	if err != nil {
		return nil, err
	}

	c.Markets[waypoint] = newMarketCacheEntry(mkt)
	return mkt, nil
}

func (c *MarketCache) MarketsInSystem(system string) ([]models.Market, error) {
	waypoints := make([]Waypoint, 0)
	result := make([]models.Market, 0)
	wps, err := c.waypoints.WaypointsInSystem(system)

	if err != nil {
		return result, err
	}

	for i := range wps {
		if wps[i].HasMarket() {
			waypoints = append(waypoints, *wps[i])
		}
	}

	for _, v := range waypoints {
		mkt, err := c.MarketForSymbol(v.Symbol)

		if err != nil {
			return nil, err
		}

		result = append(result, *mkt)
	}

	return result, nil
}

func (e marketCacheEntry) IsOld() bool {
	t := time.Now().Add(DEFAULT_MARKET_OLD)
	return e.SavedAt.Before(t)
}
