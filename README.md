# go-rate-limiter
An asynchronous rate limiting package fo Go.

It executes up to `n` jobs per `t` time.

Example usage:
```go
// Execute up to 5 requests per 2 seconds
rateLimiter := NewRateLimiter(time.Second*2, 5)

// Start jobs worker
rateLimiter.Begin()

// Stop jobs worker
defer rateLimiter.Finish()

// Queue and execute jobs
wg := sync.WaitGroup{}

logger.Println("Start requesting")

for i := 0; i < 12; i++ {
    wg.Add(1)

    currentRequestId := i

    rateLimiter.
        Do(func() (interface{}, error) {
            logger.Printf("-> Request #%d has been started\n", currentRequestId)

            // ...
        }).
        OnSuccess(func(v interface{}) {
            logger.Printf("<- Request #%d: OnSuccess: %v\n", currentRequestId, v)
        }).
        OnFailure(func(err error) {
            logger.Printf("<- Request #%d: OnFailure: %v\n", currentRequestId, err)
        }).
        OnComplete(func(v interface{}, err error) {
            //logger.Printf("<- Request #%d: OnComplete: %v\n", currentRequestId, v)

            wg.Done()
        })
}

logger.Println("Stop requesting")

wg.Wait()

logger.Println("Processing finished")
```