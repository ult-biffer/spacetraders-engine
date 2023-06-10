package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ult-biffer/spacetraders_engine/ext"
	"github.com/ult-biffer/spacetraders_sdk/api"
	"github.com/ult-biffer/spacetraders_sdk/models"
	"github.com/ult-biffer/spacetraders_sdk/responses"
)

type Game struct {
	Token     string                      `json:"token"`
	Agent     *models.Agent               `json:"agent"`
	Contracts map[string]*models.Contract `json:"contracts"`
	Ships     map[string]*ext.Ship        `json:"ships"`
	Surveys   map[string]*models.Survey   `json:"surveys"`
	Markets   *ext.MarketCache            `json:"-"`
	Waypoints *ext.WaypointCache          `json:"-"`
}

func (g *Game) LoadFromResponse(resp *responses.RegisterResponse) {
	g.Token = resp.Data.Token
	g.Agent = &resp.Data.Agent
	g.Contracts = map[string]*models.Contract{
		resp.Data.Contract.Id: &resp.Data.Contract,
	}
	g.Ships = map[string]*ext.Ship{
		resp.Data.Ship.Symbol: ext.NewShip(resp.Data.Ship, nil, nil),
	}
	g.Surveys = make(map[string]*models.Survey)

	api.GetClient().SetToken(g.Token)

	g.loadCaches()
	g.initShips()
}

func (g *Game) LoadFromSymbol(symbol string) error {
	body, err := os.ReadFile(saveFilePath(symbol))

	if err != nil {
		return err
	}

	var og *Game

	if err := json.Unmarshal(body, og); err != nil {
		return err
	}

	g.loadFromOther(og)
	return nil
}

func (g *Game) Save() error {
	if err := g.saveGame(); err != nil {
		return err
	}

	if err := g.saveWaypoints(); err != nil {
		return err
	}

	return g.saveMarkets()
}

func (g *Game) saveGame() error {
	path := saveFilePath(g.Agent.Symbol)

	if path == "" {
		return fmt.Errorf("failed to get user home directory")
	}

	return writeToPath(g, path)
}

func (g *Game) saveWaypoints() error {
	path := waypointCachePath()

	if path == "" {
		return fmt.Errorf("failed to get user home directory")
	}

	return writeToPath(g.Waypoints, path)
}

func (g *Game) saveMarkets() error {
	path := marketCachePath()

	if path == "" {
		return fmt.Errorf("failed to get user home directory")
	}

	return writeToPath(g.Markets, path)
}

func writeToPath(a any, path string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()
	body, err := json.Marshal(a)

	if err != nil {
		return err
	}

	_, err = file.Write(body)

	if err != nil {
		return err
	}

	return file.Sync()
}

func saveFilePath(agent string) string {
	home, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/.spacetraders/saves/%s.json", home, agent)
}

func marketCachePath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/.spacetraders/cache/markets.json", home)
}

func waypointCachePath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/.spacetraders/cache/waypoints.json", home)
}

func (g *Game) initShips() {
	for k, v := range g.Ships {
		wp, err := g.Waypoints.Waypoint(v.Nav.WaypointSymbol)

		if err != nil {
			continue
		}

		g.Ships[k].Waypoint = wp
	}
}

func (g *Game) loadFromOther(og *Game) {
	g.Token = og.Token
	g.Agent = og.Agent
	g.Contracts = og.Contracts
	g.Ships = og.Ships
	g.Surveys = og.Surveys

	api.GetClient().SetToken(g.Token)

	g.loadCaches()
}

func (g *Game) loadCaches() error {
	if err := g.loadWaypointCache(); err != nil {
		return err
	}

	return g.loadMarketCache()
}

func (g *Game) loadWaypointCache() error {
	path := waypointCachePath()
	body, err := os.ReadFile(path)

	if errors.Is(err, os.ErrNotExist) {
		body = []byte("{}")
	} else if err != nil {
		return err
	}

	var cache *ext.WaypointCache
	if err := json.Unmarshal(body, cache); err != nil {
		return err
	}

	g.Waypoints = cache
	return nil
}

func (g *Game) loadMarketCache() error {
	path := marketCachePath()
	body, err := os.ReadFile(path)

	if errors.Is(err, os.ErrNotExist) {
		body = []byte("{}")
	} else if err != nil {
		return err
	}

	cache := ext.LoadMarketCache(body, g.Waypoints)
	g.Markets = cache
	return nil
}
