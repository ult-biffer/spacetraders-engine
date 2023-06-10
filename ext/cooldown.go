package ext

import (
	"math"
	"time"

	"github.com/ult-biffer/spacetraders_sdk/models"
)

type Cooldown struct {
	models.Cooldown
	SavedAt time.Time `json:"savedAt"`
}

func NewCooldown(cd models.Cooldown) *Cooldown {
	return &Cooldown{
		Cooldown: cd,
		SavedAt:  time.Now(),
	}
}

func (cd *Cooldown) Expiration() int {
	seconds := time.Second * time.Duration(cd.RemainingSeconds)
	expiresAt := cd.SavedAt.Add(seconds)

	if expiresAt.After(time.Now()) {
		secondsTilExpiration := float64(time.Until(expiresAt)) / float64(time.Second)
		return int(math.Ceil(secondsTilExpiration))
	} else {
		return 0
	}
}
