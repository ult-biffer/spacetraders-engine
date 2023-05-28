package ext

import (
	"time"

	sdk "github.com/ult-biffer/spacetraders-api-go"
)

type Cooldown struct {
	sdk.Cooldown
	SavedAt time.Time `json:"savedAt"`
}

func NewCooldown(cd sdk.Cooldown) *Cooldown {
	return &Cooldown{
		Cooldown: cd,
		SavedAt:  time.Now(),
	}
}

func (cd Cooldown) Expiration() (bool, int64) {
	seconds := time.Second * time.Duration(cd.RemainingSeconds)
	expiresAt := cd.SavedAt.Add(seconds)
	secondsTilExpiration := time.Until(expiresAt) / time.Second

	if expiresAt.After(time.Now()) {
		return false, int64(secondsTilExpiration)
	} else {
		return true, 0
	}
}
