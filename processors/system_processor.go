package processors

import (
	"spacetraders_engine/api"
	"spacetraders_engine/game"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type SystemProcessor struct {
	apiSystem *api.System
	Game      *game.Game
	Symbol    string
}

func NewSystemProcessor(g *game.Game, symbol string) (*SystemProcessor, error) {
	if g.Token == "" {
		return nil, NewNotLoggedInError()
	}

	return &SystemProcessor{
		Game:      g,
		Symbol:    symbol,
		apiSystem: api.NewSystem(g.Client, symbol),
	}, nil
}

func (sp *SystemProcessor) GetWaypoints() ([]sdk.Waypoint, error) {
	return sp.Game.Waypoints.WaypointsInSystem(sp.Symbol)
}

func (sp *SystemProcessor) GetMarkets() ([]sdk.Market, error) {
	return sp.Game.Markets.MarketsInSystem(sp.Symbol)
}
