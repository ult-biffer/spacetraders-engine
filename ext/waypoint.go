package ext

import (
	"fmt"
	"math"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Waypoint struct {
	sdk.Waypoint
}

func NewWaypoint(wp sdk.Waypoint) *Waypoint {
	return &Waypoint{
		Waypoint: wp,
	}
}

func NewWaypointList(wps []sdk.Waypoint) []Waypoint {
	var result []Waypoint
	for i := range wps {
		result = append(result, *NewWaypoint(wps[i]))
	}
	return result
}

func (wp *Waypoint) DistanceTo(other sdk.Waypoint) int32 {
	// Absolute value doesn't matter because we are squaring
	dx := float64(other.X - wp.X)
	dy := float64(other.Y - wp.Y)

	// a**2 + b**2 = c**2
	dd := math.Sqrt(dx*dx + dy*dy)
	return int32(math.Ceil(dd))
}

func (wp *Waypoint) FuelCostTo(other sdk.Waypoint, ship sdk.Ship) (int32, error) {
	flightMode := ship.Nav.FlightMode

	switch flightMode {
	case sdk.SHIPNAVFLIGHTMODE_DRIFT:
		return 1, nil
	case sdk.SHIPNAVFLIGHTMODE_STEALTH, sdk.SHIPNAVFLIGHTMODE_CRUISE:
		return wp.DistanceTo(other), nil
	case sdk.SHIPNAVFLIGHTMODE_BURN:
		return (2 * wp.DistanceTo(other)), nil
	default:
		return -1, fmt.Errorf("invalid flight mode %s", flightMode)
	}
}

func (wp *Waypoint) TimeTo(other sdk.Waypoint, ship sdk.Ship) (int32, error) {
	flightMode := ship.Nav.FlightMode
	shipSpeed := ship.Engine.Speed
	distance := float32(wp.DistanceTo(other))
	hertz := distance / shipSpeed

	switch flightMode {
	case sdk.SHIPNAVFLIGHTMODE_DRIFT:
		return int32(15 + (100 * hertz)), nil
	case sdk.SHIPNAVFLIGHTMODE_STEALTH:
		return int32(15 + (20 * hertz)), nil
	case sdk.SHIPNAVFLIGHTMODE_CRUISE:
		return int32(15 + (10 * hertz)), nil
	case sdk.SHIPNAVFLIGHTMODE_BURN:
		return int32(15 + (5 * hertz)), nil
	default:
		return -1, fmt.Errorf("invalid flight mode %s", flightMode)
	}
}

func (wp *Waypoint) HasTrait(trait string) bool {
	for _, v := range wp.Traits {
		if v.Symbol == trait {
			return true
		}
	}

	return false
}

func (wp *Waypoint) IsType(t sdk.WaypointType) bool {
	return wp.Type == t
}

func (wp *Waypoint) HasMarket() bool {
	return wp.HasTrait("MARKETPLACE")
}

func (wp *Waypoint) HasShipyard() bool {
	return wp.HasTrait("SHIPYARD")
}

func (wp *Waypoint) IsAsteroidField() bool {
	return wp.IsType(sdk.WAYPOINTTYPE_ASTEROID_FIELD)
}

func (wp *Waypoint) IsJumpGate() bool {
	return wp.IsType(sdk.WAYPOINTTYPE_JUMP_GATE)
}
