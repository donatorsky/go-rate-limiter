package ratelimiter

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/donatorsky/go-promise"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
)

func TestNewRateLimiter(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("Rate limiter cannot be constructed with limit less than 1", func(t *testing.T) {
		require.PanicsWithValue(t, "limit cannot be less than 1", func() {
			NewRateLimiter(time.Duration(fakerInstance.Int64Between(1, 999999999)), 0)
		})
	})

	t.Run("Rate limiter can be constructed with limit equal 1", func(t *testing.T) {
		var rate = time.Duration(fakerInstance.Int64Between(1, 999999999))

		limiter := NewRateLimiter(rate, 1)

		require.Equal(t, rate, limiter.rate)
		require.Equal(t, uint64(1), limiter.limit)
		require.False(t, limiter.running)
	})
}

func TestRateLimiter_Rate(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("The rate may be received", func(t *testing.T) {
		var rate = time.Duration(fakerInstance.Int64Between(1, 999999999))

		limiter := NewRateLimiter(rate, 1)

		require.Equal(t, rate, limiter.rate)
		require.Equal(t, rate, limiter.Rate())
	})
}

func TestRateLimiter_Limit(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("The limit may be received", func(t *testing.T) {
		var limit = uint64(fakerInstance.Int64Between(1, 999999999))

		limiter := NewRateLimiter(time.Duration(fakerInstance.Int64Between(1, 999999999)), limit)

		require.Equal(t, limit, limiter.limit)
		require.Equal(t, limit, limiter.Limit())
	})
}

func TestRateLimiter_Begin(t *testing.T) {
	type fields struct {
		rate                    time.Duration
		limit                   uint64
		inProgressCounter       uint64
		alreadyProcessedCounter uint64
		running                 bool
		queue                   DoublyLinkedList
		rateTicker              <-chan time.Time
		worker                  chan jobDefinition
		mutex                   sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &RateLimiter{
				rate:                    tt.fields.rate,
				limit:                   tt.fields.limit,
				inProgressCounter:       tt.fields.inProgressCounter,
				alreadyProcessedCounter: tt.fields.alreadyProcessedCounter,
				running:                 tt.fields.running,
				queue:                   tt.fields.queue,
				rateTicker:              tt.fields.rateTicker,
				worker:                  tt.fields.worker,
				mutex:                   tt.fields.mutex,
			}
			sdk.Begin()
		})
	}
}

func TestRateLimiter_Do(t *testing.T) {
	type fields struct {
		rate                    time.Duration
		limit                   uint64
		inProgressCounter       uint64
		alreadyProcessedCounter uint64
		running                 bool
		queue                   DoublyLinkedList
		rateTicker              <-chan time.Time
		worker                  chan jobDefinition
		mutex                   sync.RWMutex
	}
	type args struct {
		handler jobHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   promise.Promiser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &RateLimiter{
				rate:                    tt.fields.rate,
				limit:                   tt.fields.limit,
				inProgressCounter:       tt.fields.inProgressCounter,
				alreadyProcessedCounter: tt.fields.alreadyProcessedCounter,
				running:                 tt.fields.running,
				queue:                   tt.fields.queue,
				rateTicker:              tt.fields.rateTicker,
				worker:                  tt.fields.worker,
				mutex:                   tt.fields.mutex,
			}
			if got := sdk.Do(tt.args.handler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateLimiter_Finish(t *testing.T) {
	type fields struct {
		rate                    time.Duration
		limit                   uint64
		inProgressCounter       uint64
		alreadyProcessedCounter uint64
		running                 bool
		queue                   DoublyLinkedList
		rateTicker              <-chan time.Time
		worker                  chan jobDefinition
		mutex                   sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &RateLimiter{
				rate:                    tt.fields.rate,
				limit:                   tt.fields.limit,
				inProgressCounter:       tt.fields.inProgressCounter,
				alreadyProcessedCounter: tt.fields.alreadyProcessedCounter,
				running:                 tt.fields.running,
				queue:                   tt.fields.queue,
				rateTicker:              tt.fields.rateTicker,
				worker:                  tt.fields.worker,
				mutex:                   tt.fields.mutex,
			}
			sdk.Finish()
		})
	}
}

func TestRateLimiter_enqueue(t *testing.T) {
	type fields struct {
		rate                    time.Duration
		limit                   uint64
		inProgressCounter       uint64
		alreadyProcessedCounter uint64
		running                 bool
		queue                   DoublyLinkedList
		rateTicker              <-chan time.Time
		worker                  chan jobDefinition
		mutex                   sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &RateLimiter{
				rate:                    tt.fields.rate,
				limit:                   tt.fields.limit,
				inProgressCounter:       tt.fields.inProgressCounter,
				alreadyProcessedCounter: tt.fields.alreadyProcessedCounter,
				running:                 tt.fields.running,
				queue:                   tt.fields.queue,
				rateTicker:              tt.fields.rateTicker,
				worker:                  tt.fields.worker,
				mutex:                   tt.fields.mutex,
			}
			sdk.enqueue()
		})
	}
}

func TestRateLimiter_process(t *testing.T) {
	type fields struct {
		rate                    time.Duration
		limit                   uint64
		inProgressCounter       uint64
		alreadyProcessedCounter uint64
		running                 bool
		queue                   DoublyLinkedList
		rateTicker              <-chan time.Time
		worker                  chan jobDefinition
		mutex                   sync.RWMutex
	}
	type args struct {
		definition jobDefinition
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &RateLimiter{
				rate:                    tt.fields.rate,
				limit:                   tt.fields.limit,
				inProgressCounter:       tt.fields.inProgressCounter,
				alreadyProcessedCounter: tt.fields.alreadyProcessedCounter,
				running:                 tt.fields.running,
				queue:                   tt.fields.queue,
				rateTicker:              tt.fields.rateTicker,
				worker:                  tt.fields.worker,
				mutex:                   tt.fields.mutex,
			}
			sdk.process(tt.args.definition)
		})
	}
}

func TestRateLimiter_release(t *testing.T) {
	type fields struct {
		rate                    time.Duration
		limit                   uint64
		inProgressCounter       uint64
		alreadyProcessedCounter uint64
		running                 bool
		queue                   DoublyLinkedList
		rateTicker              <-chan time.Time
		worker                  chan jobDefinition
		mutex                   sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &RateLimiter{
				rate:                    tt.fields.rate,
				limit:                   tt.fields.limit,
				inProgressCounter:       tt.fields.inProgressCounter,
				alreadyProcessedCounter: tt.fields.alreadyProcessedCounter,
				running:                 tt.fields.running,
				queue:                   tt.fields.queue,
				rateTicker:              tt.fields.rateTicker,
				worker:                  tt.fields.worker,
				mutex:                   tt.fields.mutex,
			}
			sdk.release()
		})
	}
}

func TestRateLimiter_renew(t *testing.T) {
	type fields struct {
		rate                    time.Duration
		limit                   uint64
		inProgressCounter       uint64
		alreadyProcessedCounter uint64
		running                 bool
		queue                   DoublyLinkedList
		rateTicker              <-chan time.Time
		worker                  chan jobDefinition
		mutex                   sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &RateLimiter{
				rate:                    tt.fields.rate,
				limit:                   tt.fields.limit,
				inProgressCounter:       tt.fields.inProgressCounter,
				alreadyProcessedCounter: tt.fields.alreadyProcessedCounter,
				running:                 tt.fields.running,
				queue:                   tt.fields.queue,
				rateTicker:              tt.fields.rateTicker,
				worker:                  tt.fields.worker,
				mutex:                   tt.fields.mutex,
			}
			sdk.renew()
		})
	}
}

func TestRateLimiter_reserve(t *testing.T) {
	type fields struct {
		rate                    time.Duration
		limit                   uint64
		inProgressCounter       uint64
		alreadyProcessedCounter uint64
		running                 bool
		queue                   DoublyLinkedList
		rateTicker              <-chan time.Time
		worker                  chan jobDefinition
		mutex                   sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := &RateLimiter{
				rate:                    tt.fields.rate,
				limit:                   tt.fields.limit,
				inProgressCounter:       tt.fields.inProgressCounter,
				alreadyProcessedCounter: tt.fields.alreadyProcessedCounter,
				running:                 tt.fields.running,
				queue:                   tt.fields.queue,
				rateTicker:              tt.fields.rateTicker,
				worker:                  tt.fields.worker,
				mutex:                   tt.fields.mutex,
			}
			sdk.reserve()
		})
	}
}
