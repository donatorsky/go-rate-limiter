package ratelimiter

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/mock"
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

		limiter := RateLimiter{
			rate: rate,
		}

		require.Equal(t, rate, limiter.rate)
		require.Equal(t, rate, limiter.Rate())
	})
}

func TestRateLimiter_Limit(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("The limit may be received", func(t *testing.T) {
		var limit = uint64(fakerInstance.Int64Between(1, 999999999))

		limiter := RateLimiter{
			limit: limit,
		}

		require.Equal(t, limit, limiter.limit)
		require.Equal(t, limit, limiter.Limit())
	})
}

func TestRateLimiter_Begin(t *testing.T) {
	t.Run("Begin does not start the worker when one is already running", func(t *testing.T) {
		limiter := RateLimiter{
			running: true,
			rate:    time.Minute,
		}

		require.True(t, limiter.running)
		require.Nil(t, limiter.rateTicker)

		limiter.Begin()

		require.True(t, limiter.running)
		require.Nil(t, limiter.rateTicker)
	})

	t.Run("Begin starts the worker", func(t *testing.T) {
		limiter := RateLimiter{
			rate: time.Minute,
		}

		require.False(t, limiter.running)
		require.Nil(t, limiter.rateTicker)

		limiter.Begin()

		require.True(t, limiter.running)
		require.NotNil(t, limiter.rateTicker)
	})
}

func TestRateLimiter_Do(t *testing.T) {
	t.Run("Panics when worker is not started", func(t *testing.T) {
		limiter := RateLimiter{
			running: false,
		}

		require.PanicsWithValue(t, "the worker has not been started yet", func() {
			limiter.Do(func() (interface{}, error) {
				return nil, nil
			})
		})
	})
}

func TestRateLimiter_Finish(t *testing.T) {
	t.Run("Already finished Rate Limiter is finished immediately", func(t *testing.T) {
		rateTicker := make(<-chan time.Time, 1)
		worker := make(chan jobDefinition, 1)

		rateLimiter := RateLimiter{
			running:    false,
			rateTicker: rateTicker,
			worker:     worker,
		}

		rateLimiter.Finish()

		require.False(t, rateLimiter.running)
		require.Equal(t, rateTicker, rateLimiter.rateTicker)
		require.Equal(t, worker, rateLimiter.worker)
	})

	t.Run("Not finished Rate Limiter is finished properly", func(t *testing.T) {
		rateTicker := make(<-chan time.Time, 1)
		worker := make(chan jobDefinition, 1)

		rateLimiter := RateLimiter{
			running:    true,
			rateTicker: rateTicker,
			worker:     worker,
		}

		rateLimiter.Finish()

		require.False(t, rateLimiter.running)
		require.Equal(t, rateTicker, rateLimiter.rateTicker)
		_, bar := <-rateLimiter.worker
		require.Equal(t, worker, rateLimiter.worker)
		require.False(t, bar)
	})
}

type doublyLinkedListMock struct {
	mock.Mock
}

func (l *doublyLinkedListMock) PushBack(v interface{}) {
	l.Called(v)
}

func (l *doublyLinkedListMock) PopFront() (interface{}, error) {
	args := l.Called()
	return args.Get(0), args.Error(1)
}

func (l *doublyLinkedListMock) IsEmpty() bool {
	return l.Called().Bool(0)
}

func TestRateLimiter_enqueue(t *testing.T) {
	t.Run("Panics when job cannot be queued", func(t *testing.T) {
		queueMock := doublyLinkedListMock{}
		queueMock.On("PopFront").Return(nil, errors.New("mocked error"))

		rateLimiter := RateLimiter{
			limit: 1,
			queue: &queueMock,
		}

		require.PanicsWithError(t, "unknown error when queuing next job: mocked error", rateLimiter.enqueue)
	})
}

func TestRateLimiter(t *testing.T) {
	t.Parallel()

	fakerInstance := faker.New()

	t.Run("Successful job returns resolved promise", func(t *testing.T) {
		waitGroup := newWaitGroup()
		callsStack := newCallsRegistry(2)

		rateLimiter := NewRateLimiter(time.Second, 1)

		rateLimiter.Begin()
		waitGroup.Initialize("queue", 2)

		resolutionValue := fakerInstance.Int()

		promise := rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register("root")
			waitGroup.Done("queue")

			return resolutionValue, nil
		})

		promise.Then(func(value interface{}) (interface{}, error) {
			require.Equal(t, resolutionValue, value)

			callsStack.Register("success")
			waitGroup.Done("queue")

			return nil, nil
		})

		promise.Catch(func(reason error) {
			require.NoError(t, reason)

			callsStack.Register("failure")
			waitGroup.Done("queue")
		})

		waitGroup.Wait("queue")
		rateLimiter.Finish()

		callsStack.AssertCompletedInOrder(t, []string{"root", "success"})
	})

	t.Run("Erroneous job returns rejected promise", func(t *testing.T) {
		waitGroup := newWaitGroup()
		callsStack := newCallsRegistry(2)

		rateLimiter := NewRateLimiter(time.Second, 1)

		rateLimiter.Begin()
		waitGroup.Initialize("queue", 2)

		failureReason := fakerInstance.Lorem().Sentence(6)

		promise := rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register("root")
			waitGroup.Done("queue")

			return nil, errors.New(failureReason)
		})

		promise.Then(func(value interface{}) (interface{}, error) {
			require.True(t, false, "Promise is expected to be rejected, but is resolved.")

			callsStack.Register("success")
			waitGroup.Done("queue")

			return nil, nil
		})

		promise.Catch(func(reason error) {
			require.EqualError(t, reason, failureReason)

			callsStack.Register("failure")
			waitGroup.Done("queue")
		})

		waitGroup.Wait("queue")
		rateLimiter.Finish()

		callsStack.AssertCompletedInOrder(t, []string{"root", "failure"})
	})

	t.Run("Jobs are executed sequentially when rate is 1 job/1s", func(t *testing.T) {
		waitGroup := newWaitGroup()
		callsStack := newCallsRegistry(3)

		rateLimiter := NewRateLimiter(time.Second, 1)

		rateLimiter.Begin()
		waitGroup.Initialize("queue", 3)

		jobsStarted := time.Now()

		rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register(fmt.Sprintf("Job 1: %t", lastedAtLeast(jobsStarted, 0)))
			waitGroup.Done("queue")

			return nil, nil
		})

		rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register(fmt.Sprintf("Job 1: %t", lastedAtLeast(jobsStarted, time.Second)))
			waitGroup.Done("queue")

			return nil, nil
		})

		rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register(fmt.Sprintf("Job 1: %t", lastedAtLeast(jobsStarted, time.Second*2)))
			waitGroup.Done("queue")

			return nil, nil
		})

		waitGroup.Wait("queue")
		jobsFinished := time.Now()
		rateLimiter.Finish()

		require.GreaterOrEqual(t, jobsFinished.Sub(jobsStarted), time.Second*2)
		callsStack.AssertCompletedInOrder(t, []string{"Job 1: true", "Job 1: true", "Job 1: true"})
	})

	t.Run("Jobs are executed in batches of 2 jobs/1s", func(t *testing.T) {
		waitGroup := newWaitGroup()
		callsStack := newCallsRegistry(3)

		rateLimiter := NewRateLimiter(time.Second, 2)

		waitGroup.
			Initialize("batch-1", 2).
			Initialize("batch-2", 1)

		rateLimiter.Begin()
		jobsStarted := time.Now()

		rateLimiter.Do(func() (interface{}, error) {
			time.Sleep(time.Millisecond * 500)
			callsStack.Register(fmt.Sprintf("Job 1: %t", lastedAtLeast(jobsStarted, time.Millisecond*500)))
			waitGroup.Done("batch-1")

			return nil, nil
		})

		rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register(fmt.Sprintf("Job 2: %t", lastedAtLeast(jobsStarted, 0)))
			waitGroup.Done("batch-1")

			return nil, nil
		})

		rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register(fmt.Sprintf("Job 3: %t", lastedAtLeast(jobsStarted, time.Second)))
			waitGroup.Done("batch-2")

			return nil, nil
		})

		waitGroup.Wait("batch-1")
		callsStack.AssertCurrentCallsStackIs(t, []string{"Job 1: true", "Job 2: true"})
		waitGroup.Wait("batch-2")
		jobsFinished := time.Now()
		rateLimiter.Finish()

		require.GreaterOrEqual(t, jobsFinished.Sub(jobsStarted), time.Second*1)
		callsStack.AssertCompleted(t, []string{"Job 1: true", "Job 2: true", "Job 3: true"})
	})

	t.Run("When job does not finish in one cycle, it is continued in the next cycle", func(t *testing.T) {
		waitGroup := newWaitGroup()
		callsStack := newCallsRegistry(3)

		rateLimiter := NewRateLimiter(time.Second, 2)

		waitGroup.
			Initialize("batch-1", 2).
			Initialize("batch-2", 1)

		rateLimiter.Begin()
		jobsStarted := time.Now()

		rateLimiter.Do(func() (interface{}, error) {
			time.Sleep(time.Millisecond * 1500)
			callsStack.Register(fmt.Sprintf("Job 1: %t", lastedAtLeast(jobsStarted, time.Millisecond*1500)))
			waitGroup.Done("batch-1")

			return nil, nil
		})

		rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register(fmt.Sprintf("Job 2: %t", lastedAtLeast(jobsStarted, 0)))
			waitGroup.Done("batch-1")

			return nil, nil
		})

		rateLimiter.Do(func() (interface{}, error) {
			callsStack.Register(fmt.Sprintf("Job 3: %t", lastedAtLeast(jobsStarted, time.Second)))
			waitGroup.Done("batch-2")

			return nil, nil
		})

		waitGroup.Wait("batch-2")
		callsStack.AssertCurrentCallsStackInOrderIs(t, []string{"Job 2: true", "Job 3: true"})
		waitGroup.Wait("batch-1")
		jobsFinished := time.Now()
		rateLimiter.Finish()

		require.GreaterOrEqual(t, jobsFinished.Sub(jobsStarted), time.Second*1)
		callsStack.AssertCompletedInOrder(t, []string{"Job 2: true", "Job 3: true", "Job 1: true"})
	})
}

//TODO: Temporary measure, get rid when refactored to custom reset channel
func lastedAtLeast(t time.Time, d time.Duration) bool {
	return time.Now().Sub(t) >= d
}
