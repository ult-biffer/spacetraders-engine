package ext

import (
	"math"

	"github.com/ult-biffer/spacetraders_sdk/models"
	"github.com/ult-biffer/spacetraders_sdk/util"
)

type Waypoint struct {
	models.Waypoint
}

func NewWaypoint(wp models.Waypoint) *Waypoint {
	return &Waypoint{Waypoint: wp}
}

func NewWaypointList(wps []models.Waypoint) []*Waypoint {
	var result []*Waypoint
	for i := range wps {
		result = append(result, NewWaypoint(wps[i]))
	}
	return result
}

func (wp *Waypoint) Location() *util.Location {
	return util.NewLocation(wp.Symbol)
}

func (wp *Waypoint) DistanceTo(other Waypoint) int {
	if wp.Location().System != other.Location().System {
		return -1
	}

	dx := float64(other.X - wp.X)
	dy := float64(other.Y - wp.Y)

	// a**2 + b**2 = c**2
	dd := math.Sqrt(dx*dx + dy*dy)
	return int(math.Ceil(dd))
}

func (wp *Waypoint) FuelCostTo(other Waypoint, ship Ship) int {
	fm := ship.Nav.FlightMode

	switch fm {
	case models.DRIFT:
		return 1
	case models.STEALTH, models.CRUISE:
		return wp.DistanceTo(other)
	case models.BURN:
		return 2 * wp.DistanceTo(other)
	default:
		return -1
	}
}

func (wp *Waypoint) TimeTo(other Waypoint, ship Ship) int {
	fm := ship.Nav.FlightMode
	speed := float64(ship.Engine.Speed)
	distance := float64(wp.DistanceTo(other))
	hz := distance / speed

	switch fm {
	case models.DRIFT:
		return int(15 + math.Ceil(hz*100))
	case models.STEALTH:
		return int(15 + math.Ceil(hz*20))
	case models.CRUISE:
		return int(15 + math.Ceil(hz*10))
	case models.BURN:
		return int(15 + math.Ceil(hz*5))
	default:
		return -1
	}
}

func (wp *Waypoint) HasTrait(trait models.WaypointTraitSymbol) bool {
	for _, v := range wp.Traits {
		if v.Symbol == trait {
			return true
		}
	}

	return false
}

func (wp *Waypoint) IsType(t models.WaypointType) bool {
	return t == wp.Type
}

func (wp *Waypoint) HasMarket() bool {
	return wp.HasTrait(models.WP_TRAIT_MARKETPLACE)
}

func (wp *Waypoint) HasShipyard() bool {
	return wp.HasTrait(models.WP_TRAIT_SHIPYARD)
}

func (wp *Waypoint) IsAsteroidField() bool {
	return wp.IsType(models.WP_ASTEROID_FIELD)
}

func (wp *Waypoint) IsJumpGate() bool {
	return wp.IsType(models.WP_JUMP_GATE)
}
