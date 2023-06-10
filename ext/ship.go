package ext

import (
	"github.com/ult-biffer/spacetraders_sdk/models"
)

type Ship struct {
	models.Ship
	Waypoint *Waypoint `json:"waypoint"`
	Cooldown *Cooldown `json:"cooldown"`
}

func NewShip(ship models.Ship, wp *Waypoint, cd *Cooldown) *Ship {
	return &Ship{
		Ship:     ship,
		Waypoint: wp,
		Cooldown: cd,
	}
}

func (s *Ship) OnCooldown() bool {
	if s.Cooldown == nil {
		return false
	}

	return s.Cooldown.Expiration() > 0
}
