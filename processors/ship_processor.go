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

func (sp *ShipProcessor) Refine(produce models.TradeSymbol) (p, c []models.RefinedTradeGood, err error) {
	if sp.gameShip().OnCooldown() {
		return nil, nil, fmt.Errorf("ship on cooldown for %ds", sp.gameShip().Cooldown.Expiration())
	}

	resp, err := api.ShipRefine(sp.Symbol, produce)

	if err != nil {
		return nil, nil, err
	}

	sp.Game.AddCooldown(resp.Data.Cooldown)
	sp.Game.Ships[sp.Symbol].Cargo = resp.Data.Cargo

	return resp.Data.Produced, resp.Data.Consumed, nil
}

func (sp *ShipProcessor) Refuel() (*models.MarketTransaction, error) {
	resp, err := api.RefuelShip(sp.Symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Agent = &resp.Data.Agent
	sp.gameShip().Fuel = resp.Data.Fuel

	return &resp.Data.Transaction, nil
}

func (sp *ShipProcessor) SellCargo(item models.TradeSymbol, units int) (*models.MarketTransaction, error) {
	resp, err := api.SellCargo(sp.Symbol, item, units)

	if err != nil {
		return nil, err
	}

	sp.Game.Agent = &resp.Data.Agent
	sp.Game.Ships[sp.Symbol].Cargo = resp.Data.Cargo

	return &resp.Data.Transaction, nil
}

func (sp *ShipProcessor) Survey() ([]models.Survey, error) {
	if sp.gameShip().OnCooldown() {
		return nil, fmt.Errorf("ship on cooldown for %ds", sp.gameShip().Cooldown.Expiration())
	}

	resp, err := api.CreateSurvey(sp.Symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.AddSurveys(resp.Data.Surveys)
	sp.Game.AddCooldown(resp.Data.Cooldown)

	return resp.Data.Surveys, nil
}

func (sp *ShipProcessor) TransferCargo(dest string, item models.TradeSymbol, units int) (*ext.Ship, error) {
	d, ok := sp.Game.Ships[dest]

	if !ok {
		return nil, fmt.Errorf("unknown destination ship %s", dest)
	}

	if d.Nav.WaypointSymbol != sp.gameShip().Nav.WaypointSymbol || d.Nav.Status != sp.gameShip().Nav.Status {
		return nil, fmt.Errorf("ships must be at the same waypoint with the same status")
	}

	resp, err := api.TransferCargo(sp.Symbol, dest, item, units)

	if err != nil {
		return nil, err
	}

	cargo, err := api.GetShipCargo(dest)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Cargo = resp.Data.Cargo
	sp.Game.Ships[dest].Cargo = *cargo
	return sp.gameShip(), nil
}

func (sp *ShipProcessor) Warp(waypoint string) (*ext.Ship, error) {
	wp, err := sp.Game.Waypoints.Waypoint(waypoint)

	if err != nil {
		return nil, err
	}

	resp, err := api.WarpShip(sp.Symbol, waypoint)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Nav = resp.Data.Nav
	sp.Game.Ships[sp.Symbol].Fuel = resp.Data.Fuel
	sp.Game.Ships[sp.Symbol].Waypoint = wp

	return sp.gameShip(), nil
}

func (sp *ShipProcessor) gameShip() *ext.Ship {
	return sp.Game.Ships[sp.Symbol]
}
