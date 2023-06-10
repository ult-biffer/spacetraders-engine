package processors

import (
	"fmt"

	"github.com/ult-biffer/spacetraders_engine/ext"
	"github.com/ult-biffer/spacetraders_engine/game"
	"github.com/ult-biffer/spacetraders_sdk/api"
	"github.com/ult-biffer/spacetraders_sdk/models"
)

type ShipProcessor struct {
	Game   *game.Game
	Symbol string
}

func NewShipProcessor(g *game.Game, symbol string) (*ShipProcessor, error) {
	if g.Token == "" {
		return nil, fmt.Errorf("not logged in, please login or register")
	}

	return &ShipProcessor{
		Game:   g,
		Symbol: symbol,
	}, nil
}

func (sp *ShipProcessor) Dock() (*ext.Ship, error) {
	resp, err := api.DockShip(sp.Symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Nav = *resp
	return sp.Game.Ships[sp.Symbol], nil
}

func (sp *ShipProcessor) Extract(survey string) (*models.Extraction, error) {
	if sp.gameShip().OnCooldown() {
		return nil, fmt.Errorf("ship on cooldown for %ds", sp.gameShip().Cooldown.Expiration())
	}

	s := sp.Game.Surveys[survey]
	resp, err := api.Extract(sp.Symbol, s)

	if err != nil {
		return nil, err
	}

	sp.Game.AddCooldown(resp.Data.Cooldown)
	sp.Game.Ships[sp.Symbol].Cargo = resp.Data.Cargo
	return &resp.Data.Extraction, nil
}

func (sp *ShipProcessor) Navigate(waypoint string) (*ext.Ship, error) {
	wp, err := sp.Game.Waypoints.Waypoint(waypoint)

	if err != nil {
		return nil, err
	}

	resp, err := api.NavigateShip(sp.Symbol, waypoint)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Fuel = resp.Data.Fuel
	sp.Game.Ships[sp.Symbol].Nav = resp.Data.Nav
	sp.Game.Ships[sp.Symbol].Waypoint = wp

	return sp.Game.Ships[sp.Symbol], nil
}

func (sp *ShipProcessor) Orbit() (*ext.Ship, error) {
	resp, err := api.OrbitShip(sp.Symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Nav = *resp
	return sp.Game.Ships[sp.Symbol], nil
}

func (sp *ShipProcessor) gameShip() *ext.Ship {
	return sp.Game.Ships[sp.Symbol]
}
