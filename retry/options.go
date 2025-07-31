package retry

import "time"

type Options struct {
	// Attempts specifies the number of retry attempts.
	attempts int
	// Delay is the initial delay between attempts, in seconds.
	delay float32
	// ExponentialBackoff enables exponential increase of delay between attempts.
	exponentialBackoff bool
}

func DefaultOptions() *Options {
	return &Options{
		attempts:           3,
		delay:              1.0,
		exponentialBackoff: false,
	}
}

func (opts *Options) Execute(fn func() error) (err error) {
	delay := opts.delay
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
	return err
}
