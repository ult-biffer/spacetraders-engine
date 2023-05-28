package api

import (
	"context"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Waypoint struct {
	Client *sdk.APIClient
	Symbol string
	System string
}

func NewWaypoint(client *sdk.APIClient, symbol string) (*Waypoint, error) {
	system, err := getSystemFromWaypoint(symbol)

	if err != nil {
		return nil, err
	}

	return &Waypoint{
		Client: client,
		Symbol: symbol,
		System: system,
	}, nil
}

func (w *Waypoint) Get(ctx context.Context) (sdk.Waypoint, error) {
	resp, r, err := w.Client.SystemsApi.GetWaypoint(ctx, w.System, w.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.Waypoint{}, err
	}

	return resp.Data, nil
}

func (w *Waypoint) JumpGate(ctx context.Context) (sdk.JumpGate, error) {
	resp, r, err := w.Client.SystemsApi.GetJumpGate(ctx, w.System, w.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.JumpGate{}, err
	}

	return resp.Data, nil
}

func (w *Waypoint) Market(ctx context.Context) (sdk.Market, error) {
	resp, r, err := w.Client.SystemsApi.GetMarket(ctx, w.System, w.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.Market{}, err
	}

	return resp.Data, nil
}

func (w *Waypoint) PurchaseShip(ctx context.Context, t sdk.ShipType) (sdk.PurchaseShip201ResponseData, error) {
	req := *sdk.NewPurchaseShipRequest(t, w.Symbol)
	resp, r, err := w.Client.FleetApi.PurchaseShip(ctx).PurchaseShipRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.PurchaseShip201ResponseData{}, err
	}

	return resp.Data, nil
}

func (w *Waypoint) Shipyard(ctx context.Context) (sdk.Shipyard, error) {
	resp, r, err := w.Client.SystemsApi.GetShipyard(ctx, w.System, w.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.Shipyard{}, err
	}

	return resp.Data, nil
}
