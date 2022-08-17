package ticker

import "time"

func NewTimeTicker(rate time.Duration) *TimeTicker {
	return &TimeTicker{rate: rate}
}

type TimeTicker struct {
	rate time.Duration
}

func (t *TimeTicker) Init() <-chan time.Time {
	return time.Tick(t.rate)
}

func (t *TimeTicker) Rate() time.Duration {
	return t.rate
}
