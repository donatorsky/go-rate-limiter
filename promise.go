package go_rate_limiter

import (
	"errors"
	"fmt"
)

type SuccessCallback func(v interface{})
type FailureCallback func(err error)
type CompleteCallback func(v interface{}, err error)
type State uint8

const (
	StatePending   = State(0)
	StateSucceeded = State(1)
	StateFailed    = State(2)
)

var (
	ErrPromiseAlreadyResolved = errors.New("the promise is already resolved")
)

func NewPromise() *Promise {
	return &Promise{
		state:       StatePending,
		successes:   make([]SuccessCallback, 0),
		failures:    make([]FailureCallback, 0),
		completions: make([]CompleteCallback, 0),
	}
}

type Promise struct {
	state       State
	successes   []SuccessCallback
	failures    []FailureCallback
	completions []CompleteCallback
}

func (p *Promise) OnSuccess(callback SuccessCallback) *Promise {
	p.successes = append(p.successes, callback)

	return p
}

func (p *Promise) OnFailure(callback FailureCallback) *Promise {
	p.failures = append(p.failures, callback)

	return p
}

func (p *Promise) OnComplete(callback CompleteCallback) *Promise {
	p.completions = append(p.completions, callback)

	return p
}

func (p *Promise) Resolve(v interface{}) (e error) {
	if p.state != StatePending {
		return ErrPromiseAlreadyResolved
	}

	defer func() {
		p.state = StateSucceeded

		if panicError := p.getError(recover()); panicError != nil {
			e = panicError
		}
	}()

	for _, successCallback := range p.successes {
		successCallback(v)
	}

	for _, completeCallback := range p.completions {
		completeCallback(v, nil)
	}

	return nil
}

func (p *Promise) Reject(err error) (e error) {
	if p.state != StatePending {
		return ErrPromiseAlreadyResolved
	}

	defer func() {
		p.state = StateFailed

		if panicError := p.getError(recover()); panicError != nil {
			e = panicError
		}
	}()

	for _, failureCallback := range p.failures {
		failureCallback(err)
	}

	for _, completeCallback := range p.completions {
		completeCallback(nil, err)
	}

	return nil
}

func (*Promise) getError(i interface{}) error {
	if i == nil {
		return nil
	}

	switch v := i.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		if s, ok := i.(fmt.Stringer); ok {
			return errors.New(s.String())
		} else {
			return errors.New(fmt.Sprintf("%v", i))
		}
	}
}
