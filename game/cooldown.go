package game

import (
	"time"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Cooldown struct {
	Object  *sdk.Cooldown `json:"cooldown"`
	SavedAt time.Time     `json:"savedAt"`
}

func (cd Cooldown) Expiration() (bool, int64) {
	seconds := time.Second * time.Duration(cd.Object.RemainingSeconds)
	expiresAt := cd.SavedAt.Add(seconds)
	secondsTilExpiration := time.Until(expiresAt) / time.Second

	if expiresAt.After(time.Now()) {
		return false, int64(secondsTilExpiration)
	} else {
		return true, 0
	}
}
