package api

import (
	"fmt"
	"math"
	"net/http"
	"spacetraders_engine/ext"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

const MAX_LIMIT int32 = 20
const INVALID_WAYPOINT = "invalid waypoint %s"

type InvalidWaypointError struct {
	symbol string
}

func (err *InvalidWaypointError) Error() string {
	return fmt.Sprintf(INVALID_WAYPOINT, err.symbol)
}

func NewInvalidWaypointError(symbol string) error {
	return &InvalidWaypointError{
		symbol: symbol,
	}
}

func getSystemFromWaypoint(waypoint string) (string, error) {
	loc := ext.NewLocation(waypoint)

	if !loc.HasSystem() {
		return "", NewInvalidWaypointError(waypoint)
	}

	return loc.System, nil
}

func handleHttpError(res *http.Response) {
	fmt.Println("Encountered HTTP error: ", res.Body)
}

func getPagesFromMeta(meta sdk.Meta) int32 {
	unrounded := float64(meta.Total) / float64(MAX_LIMIT)
	rounded := math.Ceil(unrounded)

	return int32(rounded)
}
