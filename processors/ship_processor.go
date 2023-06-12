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
		return nil, NewNotLoggedInError()
	}

	return &ShipProcessor{
		Game:   g,
		Symbol: symbol,
	}, nil
}

func (sp *ShipProcessor) Chart() (*models.Waypoint, error) {
	resp, err := api.CreateChart(sp.Symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Waypoint = ext.NewWaypoint(resp.Data.Waypoint)
	return &resp.Data.Waypoint, nil
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
		return nil, NewShipOnCooldownError(sp.gameShip().Cooldown)
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

func (sp *ShipProcessor) GetCooldown() (*ext.Cooldown, error) {
	resp, err := api.GetShipCooldown(sp.Symbol)

	if err != nil {
		return nil, err
	}

	if (*resp == models.Cooldown{}) {
		return ext.NewCooldown(*resp), nil
	}

	return sp.Game.AddCooldown(*resp), nil
}

func (sp *ShipProcessor) Jettison(item models.TradeSymbol, units int) (*ext.Ship, error) {
	resp, err := api.JettisonCargo(sp.Symbol, item, units)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Cargo = *resp
	return sp.gameShip(), nil
}

func (sp *ShipProcessor) Jump(system string) (*ext.Ship, error) {
	if sp.gameShip().OnCooldown() {
		return nil, NewShipOnCooldownError(sp.gameShip().Cooldown)
	}

	resp, err := api.JumpShip(sp.Symbol, system)

	if err != nil {
		return nil, err
	}

	sp.Game.AddCooldown(resp.Data.Cooldown)
	sp.Game.Ships[sp.Symbol].Nav = resp.Data.Nav
	// it's hard to predict the waypoint, so we will update it on arrival
	return sp.updateWaypoint(resp.Data.Nav.WaypointSymbol)
}

func (sp *ShipProcessor) Market() (*models.Market, error) {
	if sp.gameShip().Waypoint == nil {
		return nil, fmt.Errorf("waypoint is nil")
	}

	mkt, err := sp.Game.Markets.MarketForWaypoint(sp.gameShip().Waypoint)

	if err != nil {
		return nil, err
	}

	return mkt, nil
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

func (sp *ShipProcessor) PurchaseCargo(item models.TradeSymbol, units int) (*models.MarketTransaction, error) {
	resp, err := api.PurchaseCargo(sp.Symbol, item, units)

	if err != nil {
		return nil, err
	}

	sp.Game.Agent = &resp.Data.Agent
	sp.Game.Ships[sp.Symbol].Cargo = resp.Data.Cargo

	return &resp.Data.Transaction, nil
}

func (sp *ShipProcessor) PurchaseShip(t models.ShipType) (*models.ShipyardTransaction, error) {
	resp, err := api.PurchaseShip(t, sp.gameShip().Waypoint.Symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Agent = &resp.Data.Agent
	newSym := resp.Data.Ship.Symbol
	sp.Game.Ships[newSym] = ext.NewShip(resp.Data.Ship, nil, nil)

	// we only ignore error here because the error case returns nil
	newWp, _ := sp.Game.Waypoints.Waypoint(resp.Data.Ship.Nav.WaypointSymbol)
	sp.Game.Ships[newSym].Waypoint = newWp

	return &resp.Data.Transaction, nil
}

func (sp *ShipProcessor) Refine(produce models.TradeSymbol) (p, c []models.RefinedTradeGood, err error) {
	if sp.gameShip().OnCooldown() {
		return nil, nil, NewShipOnCooldownError(sp.gameShip().Cooldown)
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
		return nil, NewShipOnCooldownError(sp.gameShip().Cooldown)
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

func (sp *ShipProcessor) updateWaypoint(symbol string) (*ext.Ship, error) {
	wp, err := sp.Game.Waypoints.Waypoint(symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Waypoint = wp
	return sp.gameShip(), nil
}
