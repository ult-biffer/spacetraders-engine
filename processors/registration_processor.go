package processors

import (
	"github.com/ult-biffer/spacetraders_engine/game"
	"github.com/ult-biffer/spacetraders_sdk/api"
)

type RegistrationProcessor struct {
	Game    *game.Game
	Symbol  string
	Faction string
	Email   *string
}

func NewRegistrationProcessor(g *game.Game, sym, fac string, eml *string) *RegistrationProcessor {
	return &RegistrationProcessor{
		Game:    g,
		Symbol:  sym,
		Faction: fac,
		Email:   eml,
	}
}

func (rp *RegistrationProcessor) Register() error {
	resp, err := api.Register(rp.Symbol, rp.Faction, rp.Email)

	if err != nil {
		return err
	}

	rp.Game.LoadFromResponse(resp)
	return rp.Game.Save()
}
