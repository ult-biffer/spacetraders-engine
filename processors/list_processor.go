package processors

import (
	"spacetraders_engine/api"
	"spacetraders_engine/game"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type ListProcessor struct {
	Game *game.Game
}

func NewListProcessor(g *game.Game) (*ListProcessor, error) {
	if g.Token == "" {
		return nil, NewNotLoggedInError()
	}

	return &ListProcessor{
		Game: g,
	}, nil
}

func (lp *ListProcessor) GetContracts() error {
	resp, err := api.GetContracts(lp.Game.Client, lp.Game.AuthContext())

	if err != nil {
		return err
	}

	lp.Game.ReplaceContracts(resp)
	return nil
}

func (lp *ListProcessor) GetSystems() ([]sdk.System, error) {
	resp, err := api.GetSystems(lp.Game.Client, lp.Game.AuthContext())

	if err != nil {
		return []sdk.System{}, err
	}

	return resp, nil
}

func (lp *ListProcessor) GetShips() error {
	resp, err := api.GetShips(lp.Game.Client, lp.Game.AuthContext())

	if err != nil {
		return err
	}

	lp.Game.ReplaceShips(resp)
	return nil
}
