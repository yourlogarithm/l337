package retry

import (
	"fmt"
	"time"
)

type Options struct {
	/* Attempts specifies the number of retry attempts. Values:
	 * 0 - Infinite retries
	 * 1 - No retries, just execute once
	 * n - Retry n times */
	attempts int
	// Delay is the initial delay between attempts, in seconds.
	delay float32
	// ExponentialBackoff enables exponential increase of delay between attempts.
	exponentialBackoff bool
}

func NewOptions(attempts int, delay float32, exponentialBackoff bool) (*Options, error) {
	if attempts < 0 {
		return nil, fmt.Errorf("attempts must be 0 or greater, got %d", attempts)
	}
	if delay < 0 {
		return nil, fmt.Errorf("delay must be 0 or greater, got %f", delay)
	}
	return &Options{
		attempts:           attempts,
		delay:              delay,
		exponentialBackoff: exponentialBackoff,
	}, nil
}

func Default() *Options {
	return &Options{
		attempts:           1,
		delay:              0,
		exponentialBackoff: false,
	}
}

func (opts *Options) Execute(fn func() error) (err error) {
	if opts.attempts == 1 {
		return fn()
	}

	delay := opts.delay
	if opts.attempts == 0 {
		for {
			err = fn()
			if err == nil {
				return nil
			}
			time.Sleep(time.Duration(delay * float32(time.Second)))
			if opts.exponentialBackoff {
				delay *= 2
			}
		}
	} else {
		for i := 0; i < opts.attempts; i++ {
			err = fn()
			if err == nil {
				return nil
			}
			time.Sleep(time.Duration(delay * float32(time.Second)))
			if opts.exponentialBackoff {
				delay *= 2
			}
		}
	}

	return err
}
