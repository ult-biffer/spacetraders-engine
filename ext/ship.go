package ext

import (
	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Ship struct {
	sdk.Ship
	Cooldown *Cooldown `json:"cooldown"`
	Waypoint *Waypoint `json:"waypoint"`
}

func NewShip(ship sdk.Ship, cd *Cooldown, wp *Waypoint) *Ship {
	return &Ship{
		Ship:     ship,
		Cooldown: cd,
		Waypoint: wp,
	}
}

func (s *Ship) OnCooldown() bool {
	if s.Cooldown == nil {
		return false
	}

	expired, _ := s.Cooldown.Expiration()
	return !expired
}
