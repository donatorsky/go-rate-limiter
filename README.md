# Go Rate Limiter

A Go's rate limiting package for asynchronous jobs.

[![GitHub license](https://img.shields.io/github/license/donatorsky/go-rate-limiter)](https://github.com/donatorsky/go-rate-limiter/blob/main/LICENSE)

## What it does?

It allows for the execution of up to `n` jobs in `t` time. If not all jobs were completed in current `t`, during the next `t` period it will try to execute `n - unfinished jobs` jobs.

You can add as many jobs as You want. All jobs that exceed `n` limit will be queued and will run automatically, as soon as there is space for them.

## Installation

Install it using a package manager:

```shell
go get github.com/donatorsky/go-rate-limiter
```

And then use it in Your code:

```go
import "github.com/donatorsky/go-rate-limiter"
```

## Example usage

```go
// Execute up to 5 requests per 2 seconds
worker := ratelimiter.NewRateLimiter(time.Second*2, 5)

// Start jobs worker
worker.Begin()

// Stop jobs worker
defer worker.Finish()

// Queue and execute jobs
wg := sync.WaitGroup{}

log.Println("Start requesting")

for i := 0; i < 12; i++ {
    wg.Add(1)

    currentRequestId := i

    worker.
        Do(func() (interface{}, error) {
            log.Printf("-> Request #%d has been started\n", currentRequestId)

            // ...Your job execution definition...
            
            // Return job result or error
            return currentRequestId, nil
        }).
        OnSuccess(func(v interface{}) {
            // In case Do() returned nil error
            log.Printf("<- Request #%d: OnSuccess: %v\n", currentRequestId, v)
        }).
        OnFailure(func(err error) {
            // In case Do() returned non-nil error
            log.Printf("<- Request #%d: OnFailure: %v\n", currentRequestId, err)
        }).
        OnComplete(func(v interface{}, err error) {
            // Always
            //log.Printf("<- Request #%d: OnComplete: %v\n", currentRequestId, v)

            wg.Done()
        })
}

log.Println("Stop requesting")

wg.Wait()

log.Println("Processing finished")
```
