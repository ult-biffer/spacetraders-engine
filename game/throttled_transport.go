package game

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ThrottledTransport struct {
	rt http.RoundTripper
	rl *rate.Limiter
}

func NewThrottledTransport(period time.Duration, requests int, rt http.RoundTripper) http.RoundTripper {
	return &ThrottledTransport{
		rt: rt,
		rl: rate.NewLimiter(rate.Every(period), requests),
	}
}

func (tt *ThrottledTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	err := tt.rl.Wait(r.Context())

	if err != nil {
		return nil, err
	}

	return tt.rt.RoundTrip(r)
}
