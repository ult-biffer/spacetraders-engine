package processors

import (
	"spacetraders_engine/api"
	"spacetraders_engine/ext"
	"spacetraders_engine/game"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type ShipProcessor struct {
	apiShip *api.Ship
	Game    *game.Game
	Symbol  string
}

func NewShipProcessor(g *game.Game, symbol string) (*ShipProcessor, error) {
	if g.Token == "" {
		return nil, NewNotLoggedInError()
	}

	return &ShipProcessor{
		apiShip: api.NewShip(g.Client, symbol),
		Game:    g,
		Symbol:  symbol,
	}, nil
}

func (sp *ShipProcessor) Get() (*ext.Ship, error) {
	resp, err := sp.apiShip.Get(sp.Game.AuthContext())

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[resp.Symbol].Ship = resp
	return sp.updateWaypoint(resp.Nav.WaypointSymbol)
}

func (sp *ShipProcessor) Chart() (*sdk.Waypoint, error) {
	resp, err := sp.apiShip.Chart(sp.Game.AuthContext())

	if err != nil {
		return nil, err
	}

	return &resp.Waypoint, nil
}

func (sp *ShipProcessor) Cooldown() (*sdk.Cooldown, error) {
	resp, err := sp.apiShip.Cooldown(sp.Game.AuthContext())

	if err != nil {
		return nil, err
	}

	sp.Game.AddCooldown(resp)
	return &resp, nil
}

func (sp *ShipProcessor) Dock() (*ext.Ship, error) {
	resp, err := sp.apiShip.Dock(sp.Game.AuthContext())

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Nav = resp
	return sp.Game.Ships[sp.Symbol], nil
}

func (sp *ShipProcessor) Extract(survey string) (*sdk.Extraction, error) {
	if sp.gameShip().OnCooldown() {
		return nil, NewShipOnCooldownError()
	}

	s := sp.Game.Surveys[survey]
	resp, err := sp.apiShip.Extract(sp.Game.AuthContext(), s)

	if err != nil {
		return nil, err
	}

	sp.Game.AddCooldown(resp.Cooldown)
	sp.Game.Ships[sp.Symbol].Cargo = resp.Cargo
	return &resp.Extraction, nil
}

func (sp *ShipProcessor) Jettison(symbol string, units int32) (*ext.Ship, error) {
	resp, err := sp.apiShip.Jettison(sp.Game.AuthContext(), symbol, units)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Cargo = resp
	return sp.Game.Ships[sp.Symbol], nil
}

func (sp *ShipProcessor) Jump(symbol string) (*ext.Ship, error) {
	if sp.gameShip().OnCooldown() {
		return nil, NewShipOnCooldownError()
	}

	resp, err := sp.apiShip.Jump(sp.Game.AuthContext(), symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.AddCooldown(resp.Cooldown)

	if resp.Nav != nil {
		sp.Game.Ships[sp.Symbol].Nav = *resp.Nav
		_, err := sp.updateWaypoint(resp.Nav.WaypointSymbol)

		if err != nil {
			return nil, err
		}
	}

	return sp.Game.Ships[sp.Symbol], nil
}

func (sp *ShipProcessor) Market() (*sdk.Market, error) {
	if sp.gameShip().Waypoint == nil {
		return nil, api.NewInvalidWaypointError("")
	}

	if sp.gameShip().Waypoint.HasMarket() {
		mkt, err := sp.Game.Markets.MarketForSymbol(sp.gameShip().Waypoint.Symbol)

		if err != nil {
			return nil, err
		}

		return &mkt, nil
	} else {
		return nil, api.NewInvalidWaypointError(sp.gameShip().Waypoint.Symbol)
	}
}

func (sp *ShipProcessor) Navigate(symbol string) (*ext.Ship, error) {
	resp, err := sp.apiShip.Navigate(sp.Game.AuthContext(), symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Fuel = resp.Fuel
	sp.Game.Ships[sp.Symbol].Nav = resp.Nav

	return sp.updateWaypoint(resp.Nav.WaypointSymbol)
}

func (sp *ShipProcessor) Orbit() (*ext.Ship, error) {
	resp, err := sp.apiShip.Orbit(sp.Game.AuthContext())

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Nav = resp
	return sp.Game.Ships[sp.Symbol], nil
}

func (sp *ShipProcessor) PurchaseCargo(symbol string, units int32) (*sdk.MarketTransaction, error) {
	resp, err := sp.apiShip.PurchaseCargo(sp.Game.AuthContext(), symbol, units)

	if err != nil {
		return nil, err
	}

	sp.Game.Agent = &resp.Agent
	sp.Game.Ships[sp.Symbol].Cargo = resp.Cargo

	return &resp.Transaction, nil
}

func (sp *ShipProcessor) PurchaseShip(t sdk.ShipType) (*sdk.ShipyardTransaction, error) {
	wp, err := sp.apiWaypoint()

	if err != nil {
		return nil, err
	}

	resp, err := wp.PurchaseShip(sp.Game.AuthContext(), t)

	if err != nil {
		return nil, err
	}

	sp.Game.Agent = &resp.Agent
	sp.Game.Ships[resp.Ship.Symbol] = ext.NewShip(resp.Ship, nil, nil)

	newWp, err := sp.Game.Waypoints.Waypoint(resp.Ship.Nav.WaypointSymbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[resp.Ship.Symbol].Waypoint = &newWp
	return &resp.Transaction, nil
}

func (sp *ShipProcessor) Refine(produce string) (p, c []sdk.ShipRefine200ResponseDataProducedInner, err error) {
	empty := make([]sdk.ShipRefine200ResponseDataProducedInner, 0)

	if sp.gameShip().OnCooldown() {
		return empty, empty, NewShipOnCooldownError()
	}

	resp, err := sp.apiShip.Refine(sp.Game.AuthContext(), produce)

	if err != nil {
		return empty, empty, err
	}

	sp.Game.AddCooldown(resp.Cooldown)
	sp.Game.Ships[sp.Symbol].Cargo = resp.Cargo

	return resp.Produced, resp.Consumed, nil
}

func (sp *ShipProcessor) Refuel() (*sdk.MarketTransaction, error) {
	resp, err := sp.apiShip.Refuel(sp.Game.AuthContext())

	if err != nil {
		return nil, err
	}

	sp.Game.Agent = &resp.Agent
	sp.Game.Ships[sp.Symbol].Fuel = resp.Fuel

	return &resp.Transaction, nil
}

// TODO: Scans

func (sp *ShipProcessor) SellCargo(symbol string, units int32) (*sdk.MarketTransaction, error) {
	resp, err := sp.apiShip.SellCargo(sp.Game.AuthContext(), symbol, units)

	if err != nil {
		return nil, err
	}

	sp.Game.Agent = &resp.Agent
	sp.Game.Ships[sp.Symbol].Cargo = resp.Cargo

	return &resp.Transaction, nil
}

func (sp *ShipProcessor) Shipyard() (*sdk.Shipyard, error) {
	wp, err := sp.apiWaypoint()

	if err != nil {
		return nil, err
	}

	resp, err := wp.Shipyard(sp.Game.AuthContext())

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (sp *ShipProcessor) Survey() ([]sdk.Survey, error) {
	if sp.gameShip().OnCooldown() {
		return nil, NewShipOnCooldownError()
	}

	resp, err := sp.apiShip.Survey(sp.Game.AuthContext())

	if err != nil {
		return []sdk.Survey{}, err
	}

	sp.Game.AddSurveys(resp.Surveys)
	sp.Game.AddCooldown(resp.Cooldown)

	return resp.Surveys, nil
}

// TODO: Transfer Cargo

func (sp *ShipProcessor) Warp(symbol string) (*ext.Ship, error) {
	resp, err := sp.apiShip.Warp(sp.Game.AuthContext(), symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Fuel = resp.Fuel
	sp.Game.Ships[sp.Symbol].Nav = resp.Nav

	return sp.updateWaypoint(resp.Nav.WaypointSymbol)
}

func (sp *ShipProcessor) apiWaypoint() (*api.Waypoint, error) {
	ship, ok := sp.Game.Ships[sp.Symbol]

	if !ok {
		return nil, NewShipNotFoundError(sp.Symbol)
	}

	wp, err := api.NewWaypoint(sp.Game.Client, ship.Nav.WaypointSymbol)

	if err != nil {
		return nil, err
	}

	return wp, nil
}

func (sp *ShipProcessor) gameShip() *ext.Ship {
	return sp.Game.Ships[sp.Symbol]
}

func (sp *ShipProcessor) updateWaypoint(symbol string) (*ext.Ship, error) {
	wp, err := sp.Game.Waypoints.Waypoint(symbol)

	if err != nil {
		return nil, err
	}

	sp.Game.Ships[sp.Symbol].Waypoint = &wp
	return sp.Game.Ships[sp.Symbol], nil
}
