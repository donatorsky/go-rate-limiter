package ratelimiter

import (
	"fmt"
	"sync"
)

func newWaitGroup() waitGroup {
	return waitGroup{
		waitGroups: make(map[string]*sync.WaitGroup),
	}
}

type waitGroup struct {
	waitGroups map[string]*sync.WaitGroup
}

func (wg *waitGroup) Initialize(key string, initialDelta int) *waitGroup {
	if _, exists := wg.waitGroups[key]; exists {
		return wg
	}

	wg.waitGroups[key] = &sync.WaitGroup{}

	wg.waitGroups[key].Add(initialDelta)

	return wg
}

func (wg *waitGroup) Add(key string, delta int) {
	if selectedWaitGroup, exists := wg.waitGroups[key]; exists {
		selectedWaitGroup.Add(delta)
	} else {
		panic(fmt.Sprintf("the wait group %q is not initialized", key))
	}
}

func (wg *waitGroup) Wait(key string) {
	if selectedWaitGroup, exists := wg.waitGroups[key]; exists {
		selectedWaitGroup.Wait()
	} else {
		panic(fmt.Sprintf("the wait group %q is not initialized", key))
	}
}

func (wg *waitGroup) Done(key string) {
	if selectedWaitGroup, exists := wg.waitGroups[key]; exists {
		selectedWaitGroup.Done()
	} else {
		panic(fmt.Sprintf("the wait group %q is not initialized", key))
	}
}
