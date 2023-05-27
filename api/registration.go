package api

import (
	"context"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Registration struct {
	Client  *sdk.APIClient
	Symbol  string
	Faction string
}

func NewRegistration(client *sdk.APIClient, symbol string, faction string) *Registration {
	return &Registration{
		Client:  client,
		Symbol:  symbol,
		Faction: faction,
	}
}

func (reg *Registration) Register() (sdk.Register201ResponseData, error) {
	req := *sdk.NewRegisterRequest(reg.Faction, reg.Symbol)
	resp, r, err := reg.Client.DefaultApi.Register(context.Background()).RegisterRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.Register201ResponseData{}, err
	}

	return resp.Data, nil
}
