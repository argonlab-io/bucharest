package utils

import (
	"fmt"
	"time"
)

type AwaitOptions struct {
	Condition func() bool
	Timeout   time.Duration
	Delay     time.Duration
}

func Await(options *AwaitOptions) error {
	if options.Timeout == 0 {
		options.Timeout = time.Second
	}

	if options.Delay == 0 || options.Delay > options.Timeout {
		options.Delay = options.Timeout / 10
	}

	ticker := time.NewTicker(options.Delay)
	defer ticker.Stop()

	timeoutChan := time.After(options.Timeout)

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("timeout after %v", options.Timeout)
		case <-ticker.C:
			if options.Condition() {
				return nil
			}
		}
	}
}
