package api

import (
	"context"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Ship struct {
	Client *sdk.APIClient
	Symbol string
}

func GetShips(client *sdk.APIClient, ctx context.Context) ([]sdk.Ship, error) {
	req := client.FleetApi.GetMyShips(ctx).Limit(MAX_LIMIT).Page(1)
	resp, r, err := req.Execute()

	if err != nil {
		handleHttpError(r)
		return []sdk.Ship{}, err
	}

	pages := getPagesFromMeta(resp.Meta)
	result := resp.Data

	if pages > 1 {
		for i := int32(2); i <= pages; i++ {
			resp, r, err = req.Page(i).Execute()

			if err != nil {
				handleHttpError(r)
				return []sdk.Ship{}, err
			}

			result = append(result, resp.Data...)
		}
	}

	return result, nil
}

func NewShip(c *sdk.APIClient, s string) *Ship {
	return &Ship{
		Client: c,
		Symbol: s,
	}
}

func (s *Ship) Get(ctx context.Context) (sdk.Ship, error) {
	resp, r, err := s.Client.FleetApi.GetMyShip(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.Ship{}, err
	}

	return resp.Data, nil
}

func (s *Ship) Chart(ctx context.Context) (sdk.CreateChart201ResponseData, error) {
	resp, r, err := s.Client.FleetApi.CreateChart(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.CreateChart201ResponseData{}, nil
	}

	return resp.Data, nil
}

func (s *Ship) Cooldown(ctx context.Context) (sdk.Cooldown, error) {
	resp, r, err := s.Client.FleetApi.GetShipCooldown(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.Cooldown{}, err
	}

	return resp.Data, nil
}

func (s *Ship) Dock(ctx context.Context) (sdk.ShipNav, error) {
	resp, r, err := s.Client.FleetApi.DockShip(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.ShipNav{}, err
	}

	return resp.Data.Nav, nil
}

func (s *Ship) Extract(ctx context.Context, survey *sdk.Survey) (sdk.ExtractResources201ResponseData, error) {
	req := *sdk.NewExtractResourcesRequest()

	if survey != nil {
		req.SetSurvey(*survey)
	}

	resp, r, err := s.Client.FleetApi.ExtractResources(ctx, s.Symbol).ExtractResourcesRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.ExtractResources201ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) Jettison(ctx context.Context, symbol string, units int32) (sdk.ShipCargo, error) {
	req := *sdk.NewJettisonRequest(symbol, units)
	resp, r, err := s.Client.FleetApi.Jettison(ctx, s.Symbol).JettisonRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.ShipCargo{}, err
	}

	return resp.Data.Cargo, nil
}

func (s *Ship) Jump(ctx context.Context, systemSymbol string) (sdk.JumpShip200ResponseData, error) {
	req := *sdk.NewJumpShipRequest(systemSymbol)
	resp, r, err := s.Client.FleetApi.JumpShip(ctx, s.Symbol).JumpShipRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.JumpShip200ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) Navigate(ctx context.Context, waypointSymbol string) (sdk.NavigateShip200ResponseData, error) {
	req := *sdk.NewNavigateShipRequest(waypointSymbol)
	resp, r, err := s.Client.FleetApi.NavigateShip(ctx, s.Symbol).NavigateShipRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.NavigateShip200ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) Orbit(ctx context.Context) (sdk.ShipNav, error) {
	resp, r, err := s.Client.FleetApi.OrbitShip(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.ShipNav{}, err
	}

	return resp.Data.Nav, nil
}

func (s *Ship) PurchaseCargo(ctx context.Context, symbol string, units int32) (sdk.SellCargo201ResponseData, error) {
	req := *sdk.NewPurchaseCargoRequest(symbol, units)
	resp, r, err := s.Client.FleetApi.PurchaseCargo(ctx, s.Symbol).PurchaseCargoRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.SellCargo201ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) Refine(ctx context.Context, produce string) (sdk.ShipRefine200ResponseData, error) {
	req := *sdk.NewShipRefineRequest(produce)
	resp, r, err := s.Client.FleetApi.ShipRefine(ctx, s.Symbol).ShipRefineRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.ShipRefine200ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) Refuel(ctx context.Context) (sdk.RefuelShip200ResponseData, error) {
	resp, r, err := s.Client.FleetApi.RefuelShip(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.RefuelShip200ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) ScanShips(ctx context.Context) (sdk.CreateShipShipScan201ResponseData, error) {
	resp, r, err := s.Client.FleetApi.CreateShipShipScan(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.CreateShipShipScan201ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) ScanSystems(ctx context.Context) (sdk.CreateShipSystemScan201ResponseData, error) {
	resp, r, err := s.Client.FleetApi.CreateShipSystemScan(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.CreateShipSystemScan201ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) ScanWaypoints(ctx context.Context) (sdk.CreateShipWaypointScan201ResponseData, error) {
	resp, r, err := s.Client.FleetApi.CreateShipWaypointScan(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.CreateShipWaypointScan201ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) SellCargo(ctx context.Context, symbol string, units int32) (sdk.SellCargo201ResponseData, error) {
	req := *sdk.NewSellCargoRequest(symbol, units)
	resp, r, err := s.Client.FleetApi.SellCargo(ctx, s.Symbol).SellCargoRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.SellCargo201ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) Survey(ctx context.Context) (sdk.CreateSurvey201ResponseData, error) {
	resp, r, err := s.Client.FleetApi.CreateSurvey(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.CreateSurvey201ResponseData{}, err
	}

	return resp.Data, nil
}

func (s *Ship) TransferCargo(ctx context.Context, trade string, units int32, ship string) (sdk.ShipCargo, error) {
	req := *sdk.NewTransferCargoRequest(trade, units, ship)
	resp, r, err := s.Client.FleetApi.TransferCargo(ctx, s.Symbol).TransferCargoRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.ShipCargo{}, err
	}

	return resp.Data.Cargo, nil
}

func (s *Ship) Warp(ctx context.Context, waypointSymbol string) (sdk.NavigateShip200ResponseData, error) {
	req := *sdk.NewNavigateShipRequest(waypointSymbol)
	resp, r, err := s.Client.FleetApi.WarpShip(ctx, s.Symbol).NavigateShipRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.NavigateShip200ResponseData{}, err
	}

	return resp.Data, nil
}
