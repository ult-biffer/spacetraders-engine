package game

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"spacetraders_engine/api"
	"spacetraders_engine/ext"
	"time"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Game struct {
	Client    *sdk.APIClient           `json:"-"`
	Agent     *sdk.Agent               `json:"agent"`
	Contracts map[string]*sdk.Contract `json:"contracts"`
	Markets   *api.MarketCache         `json:"-"`
	Ships     map[string]*ext.Ship     `json:"ships"`
	Surveys   map[string]*sdk.Survey   `json:"survey"`
	Token     string                   `json:"token"`
	Waypoints *api.WaypointCache       `json:"-"`
}

func NewGame() *Game {
	cfg := sdk.NewConfiguration()
	transport := NewThrottledTransport(time.Second, 2, http.DefaultTransport)
	cfg.HTTPClient = &http.Client{Transport: transport}
	client := sdk.NewAPIClient(cfg)

	return &Game{Client: client}
}

func (g *Game) Save() error {
	path := saveFilePath(g.Agent.Symbol)

	if path == "" {
		return fmt.Errorf("failed to get user config directory")
	}

	if err := g.writeToPath(path); err != nil {
		return err
	}

	return nil
}

func (g *Game) LoadFrom201(data sdk.Register201ResponseData) {
	g.Agent = &data.Agent
	g.Contracts = make(map[string]*sdk.Contract)
	g.Ships = make(map[string]*ext.Ship)
	g.Surveys = make(map[string]*sdk.Survey)

	g.Contracts[data.Contract.Id] = &data.Contract
	g.Ships[data.Ship.Symbol] = ext.NewShip(data.Ship, nil, nil)
	g.Token = data.Token

	g.initCaches()
	g.initShips()
}

func (g *Game) LoadFromSymbol(symbol string) error {
	path := saveFilePath(symbol)
	body, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	var ng *Game

	if err := json.Unmarshal(body, ng); err != nil {
		return err
	}

	g.loadFromOther(ng)
	return nil
}

func (g *Game) AddContracts(contracts []sdk.Contract) {
	for i := range contracts {
		g.Contracts[contracts[i].Id] = &contracts[i]
	}
}

func (g *Game) AddCooldown(cd sdk.Cooldown) {
	g.Ships[cd.ShipSymbol].Cooldown = ext.NewCooldown(cd)
}

func (g *Game) AddSurveys(surveys []sdk.Survey) {
	for i := range surveys {
		g.Surveys[surveys[i].Signature] = &surveys[i]
	}
}

func (g *Game) AuthContext() context.Context {
	return context.WithValue(context.Background(), sdk.ContextAccessToken, g.Token)
}

func (g *Game) ReplaceContracts(contracts []sdk.Contract) {
	result := make(map[string]*sdk.Contract)
	for i := range contracts {
		result[contracts[i].Id] = &contracts[i]
	}

	g.Contracts = result
}

func (g *Game) ReplaceShips(ships []sdk.Ship) {
	result := make(map[string]*ext.Ship)
	for i := range ships {
		if v, ok := g.Ships[ships[i].Symbol]; ok {
			result[ships[i].Symbol] = ext.NewShip(ships[i], v.Cooldown, v.Waypoint)
		} else {
			result[ships[i].Symbol] = ext.NewShip(ships[i], nil, nil)
		}
	}

	g.Ships = result
	g.initShips()
}

func (g *Game) ShipSymbol(input string) string {
	return fmt.Sprintf("%s-%s", g.Agent.Symbol, input)
}

func (g *Game) loadFromOther(ng *Game) {
	g.Agent = ng.Agent
	g.Contracts = ng.Contracts
	g.Ships = ng.Ships
	g.Surveys = ng.Surveys
	g.Token = ng.Token

	g.initCaches()
}

func (g *Game) initCaches() {
	g.Waypoints = api.NewWaypointCache(g.Client, g.AuthContext())
	g.Markets = api.NewMarketCache(g.Client, g.AuthContext(), g.Waypoints)
}

func (g *Game) initShips() {
	for k, v := range g.Ships {
		wp, err := g.Waypoints.Waypoint(v.Nav.WaypointSymbol)

		if err != nil {
			continue
		}

		g.Ships[k].Waypoint = &wp
	}
}

func (g *Game) writeToPath(path string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	body, err := json.Marshal(g)

	if err != nil {
		return err
	}

	_, err = file.Write(body)

	if err != nil {
		return err
	}

	return file.Sync()
}

func saveFilePath(symbol string) string {
	cfg, err := os.UserConfigDir()

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/spacetraders/saves/%s.json", cfg, symbol)
}
