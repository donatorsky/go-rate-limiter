package ticker

import "time"

func NewManualTicker() *ManualTicker {
	return &ManualTicker{
		ticks: make(chan time.Time, 1),
	}
}

type ManualTicker struct {
	ticks chan time.Time
}

func (t *ManualTicker) Init() <-chan time.Time {
	return t.ticks
}

func (t *ManualTicker) Tick() {
	t.ticks <- time.Now()
}
