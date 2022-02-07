package ratelimiter

import (
	"sync"
	"time"
)

func NewRateLimiter(rate time.Duration, limit uint64) *RateLimiter {
	if limit < 1 {
		panic("limit cannot be less than 1")
	}

	return &RateLimiter{
		rate:   rate,
		limit:  limit,
		queue:  DoublyLinkedList{},
		worker: make(chan jobDefinition, limit),
	}
}

type jobHandler func() (interface{}, error)

type jobDefinition struct {
	handler jobHandler
	promise *Promise
}

type RateLimiter struct {
	rate  time.Duration
	limit uint64

	inProgressCounter       uint64
	alreadyProcessedCounter uint64
	running                 bool
	queue                   DoublyLinkedList
	rateTicker              <-chan time.Time
	worker                  chan jobDefinition
	rw                      sync.RWMutex
}

func (sdk *RateLimiter) Rate() time.Duration {
	return sdk.rate
}

func (sdk *RateLimiter) Limit() uint64 {
	return sdk.limit
}

func (sdk *RateLimiter) Begin() {
	if sdk.running {
		return
	}

	sdk.rateTicker = time.Tick(sdk.rate)

	go func() {
		for {
			select {
			case <-sdk.rateTicker:
				sdk.renew()

			case definition := <-sdk.worker:
				go sdk.process(definition)

			default:
				if sdk.rateTicker == nil {
					return
				}

				sdk.enqueue()
			}
		}
	}()

	sdk.running = true
}

func (sdk *RateLimiter) Finish() {
	if !sdk.running {
		return
	}

	sdk.rateTicker = nil
	close(sdk.worker)

	sdk.running = false
}

func (sdk *RateLimiter) Do(handler jobHandler) *Promise {
	if !sdk.running {
		panic("the worker has not been started yet")
	}

	newPromise := NewPromise().
		OnComplete(func(_ interface{}, _ error) { sdk.release() })

	sdk.rw.Lock()
	sdk.queue.PushBack(jobDefinition{
		handler: handler,
		promise: newPromise,
	})
	sdk.rw.Unlock()

	return newPromise
}

func (sdk *RateLimiter) enqueue() {
	sdk.rw.Lock()
	defer sdk.rw.Unlock()

	if (sdk.alreadyProcessedCounter + sdk.inProgressCounter) >= sdk.limit {
		return
	}

	if sdk.queue.IsEmpty() {
		return
	}

	x, err := sdk.queue.PopFront()
	if err != nil {
		return
	}

	sdk.reserve()

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
	sdk.rw.Lock()

	sdk.alreadyProcessedCounter = sdk.inProgressCounter

	sdk.rw.Unlock()
}

func (sdk *RateLimiter) reserve() {
	sdk.inProgressCounter += 1
}

func (sdk *RateLimiter) release() {
	sdk.rw.Lock()

	sdk.inProgressCounter -= 1
	sdk.alreadyProcessedCounter += 1

	sdk.rw.Unlock()
}
