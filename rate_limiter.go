package ratelimiter

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/donatorsky/go-promise"
	"github.com/donatorsky/go-rate-limiter/internal"
	tickers "github.com/donatorsky/go-rate-limiter/ticker"
)

func NewRateLimiter(ticker tickers.Ticker, limit uint64) *RateLimiter {
	if limit < 1 {
		panic("limit cannot be less than 1")
	}

	return &RateLimiter{
		limit:  limit,
		ticker: ticker,
		queue:  &internal.DoublyLinkedList{},
		worker: make(chan jobDefinition, limit),
	}
}

type jobHandler func() (interface{}, error)

type jobDefinition struct {
	handler        jobHandler
	promise        promise.Promiser
	processRequest chan bool
}

//go:generate moq -out ticker_moq_test.go . ticker
type ticker = tickers.Ticker

type RateLimiter struct {
	limit uint64

	inProgressCounter       uint64
	alreadyProcessedCounter uint64
	running                 bool
	queue                   internal.DoublyLinkedListInterface
	ticks                   <-chan time.Time
	ticker                  ticker
	worker                  chan jobDefinition
	mutex                   sync.RWMutex
}

func (sdk *RateLimiter) Limit() uint64 {
	return sdk.limit
}

func (sdk *RateLimiter) Begin() {
	sdk.mutex.Lock()
	defer sdk.mutex.Unlock()

	if sdk.running {
		return
	}

	sdk.ticks = sdk.ticker.Init()

	go func() {
		for {
			select {
			case <-sdk.ticks:
				sdk.renew()

			case definition, ok := <-sdk.worker:
				if ok {
					go sdk.process(definition)
				}

			default:
				sdk.mutex.Lock()

				if !sdk.running {
					sdk.ticks = nil
					sdk.mutex.Unlock()

					return
				}

				sdk.mutex.Unlock()

				sdk.enqueue()
			}
		}
	}()

	sdk.running = true
}

func (sdk *RateLimiter) Finish() {
	sdk.mutex.Lock()
	defer sdk.mutex.Unlock()

	if !sdk.running {
		return
	}

	close(sdk.worker)

	sdk.running = false
}

func (sdk *RateLimiter) Do(handler jobHandler) promise.Promiser {
	sdk.mutex.RLock()

	if !sdk.running {
		sdk.mutex.RUnlock()

		panic("the worker has not been started yet")
	}

	sdk.mutex.RUnlock()

	newPromise := promise.Pending()

	newPromise.Finally(sdk.release)

	sdk.mutex.Lock()
	sdk.queue.PushBack(jobDefinition{
		handler: handler,
		promise: newPromise,
	})
	sdk.mutex.Unlock()

	return newPromise
}

func (sdk *RateLimiter) enqueue() {
	sdk.mutex.Lock()

	if (sdk.alreadyProcessedCounter + sdk.inProgressCounter) >= sdk.limit {
		sdk.mutex.Unlock()

		return
	}

	x, err := sdk.queue.PopFront()
	if err != nil {
		defer sdk.mutex.Unlock()

		if errors.Is(err, internal.ErrListIsEmpty) {
			return
		}

		// Currently, the ErrListIsEmpty error is the only error possible, but panic just in case it changed
		panic(fmt.Errorf("unknown error when queuing next job: %w", err))
	}

	sdk.reserve()

	sdk.mutex.Unlock()

	sdk.worker <- x.(jobDefinition)
}

func (sdk *RateLimiter) process(definition jobDefinition) {
	if result, err := definition.handler(); err == nil {
		_ = definition.promise.Resolve(result)
	} else {
		_ = definition.promise.Reject(err)
	}
}

func (sdk *RateLimiter) renew() {
	sdk.mutex.Lock()

	sdk.alreadyProcessedCounter = 0

	sdk.mutex.Unlock()
}

func (sdk *RateLimiter) reserve() {
	sdk.inProgressCounter += 1
}

func (sdk *RateLimiter) release() {
	sdk.mutex.Lock()

	sdk.inProgressCounter -= 1
	sdk.alreadyProcessedCounter += 1

	sdk.mutex.Unlock()
}
