package proxyguard

import (
	"context"
	"errors"
	"time"
)

var ErrMaxRestarts = errors.New("max restarts exceeded")

// restartUntilErr retries a function `work` in a loop, taking a context and boolean as argument
// This boolean is true if it is the first try
// If the function returns an error the restart loop is immediately stopped,
// thus it begins to restart if the function returns a `nil` error.
// The time to wait between restarts is in the `wt` duration slice, which is cycled through.
// Furthermore, A restart is a 'failure' if the total execution time
// of the function is less than delta `d`. If the function takes more than delta `d`
// the wait time is reset to the first value of `wt`
// Another exit condition is based on the number of max failed restarts, this is len(wt)
func restartUntilErr(ctx context.Context, work func(context.Context, bool) error, wt []time.Duration, d time.Duration) error {
	if len(wt) == 0 {
		return errors.New("no restart wait times available")
	}
	failed := 0
	first := true
	for {
		st := time.Now()
		err := work(ctx, first)
		et := time.Now()
		if err != nil {
			return err
		}
		// if the time it takes for the function is less than delta
		// we consider it as failed
		if et.Sub(st) < d {
			// max tries exceeded now
			if failed == len(wt)-1 {
				return ErrMaxRestarts
			}
			failed = (failed + 1) % len(wt)
		} else {
			failed = 0
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wt[failed]):
		}
		first = false
	}
}
