package common

import "time"

// type retryOptions struct {
// 	MaxRetries int
// 	SleepTime  time.Duration
// }

// Builder Mode, default values: MaxRetries=3, SleepTime=1s
// Retry(f).MaxRetries(3).SleepTime(time.Second).Execute()

// RetryConfig holds the configuration for the retry mechanism.
type RetryConfig struct {
	MaxRetries int
	SleepTime  time.Duration
}

// Retry is the main struct that will use the Builder pattern.
type Retry struct {
	config RetryConfig
}

// NewRetry creates a new Retry instance with default values.
func NewRetry() *Retry {
	return &Retry{
		config: RetryConfig{
			MaxRetries: 3,               // default MaxRetries
			SleepTime:  1 * time.Second, // default SleepTime
		},
	}
}

// MaxRetries sets the maximum number of retries.
func (r *Retry) MaxRetries(retries int) *Retry {
	r.config.MaxRetries = retries
	return r
}

// SleepTime sets the sleep duration between retries.
func (r *Retry) SleepTime(duration time.Duration) *Retry {
	r.config.SleepTime = duration
	return r
}

// Execute runs the provided function with the configured retry mechanism.
func (r *Retry) Execute(f func() error) error {
	var err error
	for i := 0; i < r.config.MaxRetries; i++ {
		err = f()
		if err == nil {
			return nil // successful execution
		}
		time.Sleep(r.config.SleepTime) // wait before retrying
	}
	return err // return the last error after all retries
}
