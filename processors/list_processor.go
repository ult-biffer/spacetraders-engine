package processors

import (
	"fmt"

	"github.com/ult-biffer/spacetraders_engine/game"
	"github.com/ult-biffer/spacetraders_sdk/api"
)

type ListProcessor struct {
	Game *game.Game
}

func NewListProcessor(g *game.Game) (*ListProcessor, error) {
	if g.Token == "" {
		return nil, fmt.Errorf("not logged in, please login or register")
	}

	return &ListProcessor{Game: g}, nil
}

func (lp *ListProcessor) GetShips() error {
	resp, err := api.ListAllShips()

	if err != nil {
		return err
	}

	lp.Game.ReplaceShips(resp)
	return nil
}
