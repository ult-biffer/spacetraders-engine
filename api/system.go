package api

import (
	"context"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type System struct {
	Client *sdk.APIClient
	Symbol string
}

func GetSystems(client *sdk.APIClient, ctx context.Context) ([]sdk.System, error) {
	req := client.SystemsApi.GetSystems(ctx).Limit(MAX_LIMIT).Page(1)
	resp, r, err := req.Execute()

	if err != nil {
		handleHttpError(r)
		return []sdk.System{}, err
	}

	pages := getPagesFromMeta(resp.Meta)
	result := resp.Data

	if pages > 1 {
		for i := int32(2); i <= pages; i++ {
			resp, r, err = req.Page(i).Execute()

			if err != nil {
				handleHttpError(r)
				return []sdk.System{}, err
			}

			result = append(result, resp.Data...)
		}
	}

	return result, nil
}

func GetSystemsPage(client *sdk.APIClient, ctx context.Context, page int32) ([]sdk.System, error) {
	req := client.SystemsApi.GetSystems(ctx).Limit(MAX_LIMIT).Page(page)
	resp, r, err := req.Execute()

	if err != nil {
		handleHttpError(r)
		return []sdk.System{}, err
	}

	return resp.Data, nil
}

func NewSystem(client *sdk.APIClient, symbol string) *System {
	return &System{
		Client: client,
		Symbol: symbol,
	}
}

func (s *System) Get(ctx context.Context) (sdk.System, error) {
	resp, r, err := s.Client.SystemsApi.GetSystem(ctx, s.Symbol).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.System{}, err
	}

	return resp.Data, nil
}

func (s *System) GetWaypoints(ctx context.Context) ([]sdk.Waypoint, error) {
	req := s.Client.SystemsApi.GetSystemWaypoints(ctx, s.Symbol).Limit(MAX_LIMIT).Page(1)
	resp, r, err := req.Execute()

	if err != nil {
		handleHttpError(r)
		return []sdk.Waypoint{}, err
	}

	pages := getPagesFromMeta(resp.Meta)
	result := resp.Data

	if pages > 1 {
		for i := int32(2); i <= pages; i++ {
			resp, r, err = req.Page(i).Execute()

			if err != nil {
				handleHttpError(r)
				return []sdk.Waypoint{}, err
			}

			result = append(result, resp.Data...)
		}
	}

	return result, nil
}
