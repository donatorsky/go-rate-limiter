package ratelimiter

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newCallsRegistry(expectedCalls uint) CallsRegistry {
	return CallsRegistry{
		expectedCalls: expectedCalls,
	}
}

type CallsRegistry struct {
	mutex sync.RWMutex

	registry      []string
	expectedCalls uint
}

func (r *CallsRegistry) Register(place string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if uint(len(r.registry)) > r.expectedCalls {
		panic(fmt.Sprintf(
			"trying to register an unexpected call: %s; already registered all calls: %v",
			place,
			r.registry,
		))
	}

	r.registry = append(r.registry, place)
}

func (r *CallsRegistry) Summarize() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return strings.Join(r.registry, "|")
}

func (r *CallsRegistry) AssertCompletedBefore(t *testing.T, expectedRegistry []string, timeLimit time.Duration) {
	sort.Strings(expectedRegistry)

	r.assertCallsStacksAreSameBefore(t, func() ([]string, []string) {
		r.mutex.Lock()
		currentRegistry := make([]string, len(r.registry))
		copy(currentRegistry, r.registry)
		r.mutex.Unlock()

		sort.Strings(currentRegistry)

		return expectedRegistry, currentRegistry
	}, timeLimit)
}

func (r *CallsRegistry) AssertCompletedInOrder(t *testing.T, expectedRegistry []string) {
	r.assertCallsStacksAreSame(
		t,
		func() ([]string, []string) { return expectedRegistry, r.registry },
	)
}

func (r *CallsRegistry) AssertCompleted(t *testing.T, expectedRegistry []string) {
	r.assertCallsStacksAreSame(
		t,
		func() ([]string, []string) {
			r.mutex.Lock()
			currentRegistry := make([]string, len(r.registry))
			copy(currentRegistry, r.registry)
			r.mutex.Unlock()

			sort.Strings(currentRegistry)

			return expectedRegistry, currentRegistry
		},
	)
}

func (r *CallsRegistry) AssertCompletedInOrderBefore(t *testing.T, expectedRegistry []string, timeLimit time.Duration) {
	r.assertCallsStacksAreSameBefore(
		t,
		func() ([]string, []string) { return expectedRegistry, r.registry },
		timeLimit,
	)
}

func (r *CallsRegistry) AssertCompletedCallsStackIsEmpty(t *testing.T) {
	require.Empty(t, r.registry)
	r.AssertCurrentCallsStackIsEmpty(t)
}

func (r *CallsRegistry) AssertCurrentCallsStackIs(t *testing.T, expectedRegistry []string) {
	if nil == expectedRegistry {
		require.Empty(t, r.registry)

		return
	}

	r.mutex.Lock()
	currentRegistry := make([]string, len(r.registry))
	copy(currentRegistry, r.registry)
	r.mutex.Unlock()

	sort.Strings(currentRegistry)
	sort.Strings(expectedRegistry)

	require.Equal(t, expectedRegistry, currentRegistry)
}

func (r *CallsRegistry) AssertCurrentCallsStackInOrderIs(t *testing.T, expectedRegistry []string) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	require.Equal(t, expectedRegistry, r.registry)
}

func (r *CallsRegistry) AssertCurrentCallsStackIsEmpty(t *testing.T) {
	r.AssertCurrentCallsStackIs(t, nil)
}

func (r *CallsRegistry) AssertThereAreNCallsLeft(t *testing.T, numberOfCallsLeft uint) {
	numberOfCurrentCalls := uint(len(r.registry))

	require.LessOrEqual(t, numberOfCurrentCalls, r.expectedCalls)
	require.Equal(t, numberOfCallsLeft, r.expectedCalls-numberOfCurrentCalls)
}

func (r *CallsRegistry) assertCallsStacksAreSame(t *testing.T, h func() ([]string, []string)) {
	expectedRegistry, currentRegistry := h()

	require.Equal(t, expectedRegistry, currentRegistry)
}

func (r *CallsRegistry) assertCallsStacksAreSameBefore(t *testing.T, h func() ([]string, []string), timeLimit time.Duration) {
	timeLimiter := time.After(timeLimit)

	for {
		r.mutex.RLock()
		expectedRegistry, currentRegistry := h()
		r.mutex.RUnlock()

		select {
		case <-timeLimiter:
			require.FailNowf(
				t,
				"Calls registry assertion timeout",
				"There are still %d expected call(s) left. Calls registered (%d): %v.",
				r.expectedCalls-uint(len(currentRegistry)),
				len(currentRegistry),
				currentRegistry,
			)
			return

		default:
			if 0 == r.expectedCalls {
				time.Sleep(timeLimit)
			}

			if uint(len(currentRegistry)) < r.expectedCalls {
				continue
			}

			require.Equal(t, expectedRegistry, currentRegistry)
			return
		}
	}
}
