package ticker

import "time"

//go:generate moq -out ticker_moq_test.go . Ticker
type Ticker interface {
	Init() <-chan time.Time
}
