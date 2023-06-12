package processors

import (
	"fmt"

	"github.com/ult-biffer/spacetraders_engine/ext"
)

const (
	NOT_LOGGED_IN = "not logged in, please login or register"
	SHIP_ON_CD    = "ship on cooldown for %ds"
)

type NotLoggedInError struct{}

func NewNotLoggedInError() error {
	return &NotLoggedInError{}
}

func (*NotLoggedInError) Error() string {
	return NOT_LOGGED_IN
}

type ShipOnCooldownError struct {
	Cooldown *ext.Cooldown
}

func NewShipOnCooldownError(cd *ext.Cooldown) error {
	return &ShipOnCooldownError{Cooldown: cd}
}

func (err *ShipOnCooldownError) Error() string {
	return fmt.Sprintf(SHIP_ON_CD, err.Cooldown.Expiration())
}
