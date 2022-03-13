package ratelimiter

import (
	"fmt"
	"sync"
)

func newWaitGroup() WaitGroup {
	return WaitGroup{
		waitGroups: make(map[string]*sync.WaitGroup),
	}
}

type WaitGroup struct {
	waitGroups map[string]*sync.WaitGroup
}

func (wg *WaitGroup) Initialize(key string, initialDelta int) *WaitGroup {
	if _, exists := wg.waitGroups[key]; exists {
		return wg
	}

	wg.waitGroups[key] = &sync.WaitGroup{}

	wg.waitGroups[key].Add(initialDelta)

	return wg
}

func (wg *WaitGroup) Add(key string, delta int) {
	if selectedWaitGroup, exists := wg.waitGroups[key]; exists {
		selectedWaitGroup.Add(delta)
	} else {
		panic(fmt.Sprintf("the wait group %q is not initialized", key))
	}
}

func (wg *WaitGroup) Wait(key string) {
	if selectedWaitGroup, exists := wg.waitGroups[key]; exists {
		selectedWaitGroup.Wait()
	} else {
		panic(fmt.Sprintf("the wait group %q is not initialized", key))
	}
}

func (wg *WaitGroup) Done(key string) {
	if selectedWaitGroup, exists := wg.waitGroups[key]; exists {
		selectedWaitGroup.Done()
	} else {
		panic(fmt.Sprintf("the wait group %q is not initialized", key))
	}
}
