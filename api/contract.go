package api

import (
	"context"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Contract struct {
	Client *sdk.APIClient
	Id     string
}

func GetContracts(client *sdk.APIClient, ctx context.Context) ([]sdk.Contract, error) {
	req := client.ContractsApi.GetContracts(ctx).Limit(MAX_LIMIT).Page(1)
	resp, r, err := req.Execute()

	if err != nil {
		handleHttpError(r)
		return []sdk.Contract{}, err
	}

	pages := getPagesFromMeta(resp.Meta)
	result := resp.Data

	if pages > 1 {
		for i := int32(2); i <= pages; i++ {
			resp, r, err = req.Page(i).Execute()

			if err != nil {
				handleHttpError(r)
				return []sdk.Contract{}, err
			}

			result = append(result, resp.Data...)
		}
	}

	return result, nil
}

func NewContract(client *sdk.APIClient, id string) *Contract {
	return &Contract{
		Client: client,
		Id:     id,
	}
}

func (c *Contract) Get(ctx context.Context) (sdk.Contract, error) {
	resp, r, err := c.Client.ContractsApi.GetContract(ctx, c.Id).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.Contract{}, err
	}

	return resp.Data, nil
}

func (c *Contract) Accept(ctx context.Context) (sdk.AcceptContract200ResponseData, error) {
	resp, r, err := c.Client.ContractsApi.AcceptContract(ctx, c.Id).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.AcceptContract200ResponseData{}, err
	}

	return resp.Data, nil
}

func (c *Contract) Deliver(ctx context.Context, ship string, trade string, units int32) (sdk.DeliverContract200ResponseData, error) {
	req := *sdk.NewDeliverContractRequest(ship, trade, units)
	resp, r, err := c.Client.ContractsApi.DeliverContract(ctx, c.Id).DeliverContractRequest(req).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.DeliverContract200ResponseData{}, err
	}

	return resp.Data, nil
}

func (c *Contract) Fulfill(ctx context.Context) (sdk.AcceptContract200ResponseData, error) {
	resp, r, err := c.Client.ContractsApi.FulfillContract(ctx, c.Id).Execute()

	if err != nil {
		handleHttpError(r)
		return sdk.AcceptContract200ResponseData{}, err
	}

	return resp.Data, nil
}
