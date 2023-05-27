package api

import (
	"context"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Agent struct {
	Client *sdk.APIClient
	Object sdk.Agent
}

func NewAgent(client *sdk.APIClient, object sdk.Agent) *Agent {
	return &Agent{
		Client: client,
		Object: object,
	}
}

func (a *Agent) Get(ctx context.Context) (sdk.Agent, error) {
	resp, r, err := a.Client.AgentsApi.GetMyAgent(ctx).Execute()

	if err != nil {
		handleHttpError(r)
		return a.Object, err
	}

	a.Object = resp.Data
	return a.Object, nil
}
