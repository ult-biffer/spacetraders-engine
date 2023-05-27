package processors

import "fmt"

const NOT_LOGGED_IN = "not logged in, login or register to continue"
const SHIP_ON_COOLDOWN = "ship on cooldown"
const SHIP_NOT_FOUND = "could not find ship %s"

type NotLoggedInError struct{}

func (*NotLoggedInError) Error() string {
	return NOT_LOGGED_IN
}

func NewNotLoggedInError() error {
	return &NotLoggedInError{}
}

type ShipOnCooldownError struct{}

func (*ShipOnCooldownError) Error() string {
	return SHIP_ON_COOLDOWN
}

func NewShipOnCooldownError() error {
	return &ShipOnCooldownError{}
}

type ShipNotFoundError struct {
	symbol string
}

func (err *ShipNotFoundError) Error() string {
	return fmt.Sprintf(SHIP_NOT_FOUND, err.symbol)
}

func NewShipNotFoundError(symbol string) error {
	return &ShipNotFoundError{
		symbol: symbol,
	}
}
