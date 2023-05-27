package processors

import (
	"spacetraders_engine/api"
	"spacetraders_engine/game"
)

type RegistrationProcessor struct {
	Game    *game.Game
	Symbol  string
	Faction string
}

func NewRegistrationProcessor(game *game.Game, symbol string, faction string) *RegistrationProcessor {
	return &RegistrationProcessor{
		Game:    game,
		Symbol:  symbol,
		Faction: faction,
	}
}

func (rp *RegistrationProcessor) Register() error {
	apiRegistration := api.NewRegistration(rp.Game.Client, rp.Symbol, rp.Faction)
	resp, err := apiRegistration.Register()

	if err != nil {
		return err
	}

	rp.Game.LoadFrom201(resp)

	if err := rp.Game.Save(); err != nil {
		return err
	}

	return nil
}
