package ticker

import "time"

type Ticker interface {
	Init() <-chan time.Time
}
