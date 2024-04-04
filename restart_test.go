package proxyguard

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRestartUntilErr(t *testing.T) {
	// test boolean 'first' argument
	gf := false

	_ = restartUntilErr(context.Background(), func(_ context.Context, first bool) error {
		gf = first
		return nil
	}, []time.Duration{0 * time.Second}, time.Duration(1*time.Hour))

	if !gf {
		t.Fatalf("first argument true is not passed")
	}

	tErr := errors.New("test")
	// test if a function exits when an error occurs
	gErr := restartUntilErr(context.Background(), func(_ context.Context, _ bool) error {
		return tErr
	}, []time.Duration{0 * time.Second}, time.Duration(1*time.Hour))

	if !errors.Is(gErr, tErr) {
		t.Fatalf("got error: %v, not equal to want error: %v", gErr, tErr)
	}

	// test if a function retries n times
	wt := []time.Duration{
		1 * time.Microsecond,
		1 * time.Microsecond,
		1 * time.Microsecond,
		1 * time.Microsecond,
		1 * time.Microsecond,
	}
	d := 1 * time.Hour
	n := len(wt)
	restarted := 0

	gErr = restartUntilErr(context.Background(), func(_ context.Context, _ bool) error {
		restarted++
		return nil
	}, wt, d)

	if restarted != n {
		t.Fatalf("restart count: %v, not equal to want count: %v", restarted, n)
	}

	if !errors.Is(gErr, ErrMaxRestarts) {
		t.Fatalf("restart error is not max restarts: %v", gErr)
	}

	// test again and set a limit to 10 calls
	restarted = 0
	limit := 10

	// now the function that is restarted is always slower
	// so it should restart until the limit
	d = 1 * time.Microsecond
	errLimit := errors.New("limit exceeded")
	gErr = restartUntilErr(context.Background(), func(_ context.Context, _ bool) error {
		restarted++

		// add some wait so it's always slower
		time.Sleep(5 * time.Microsecond)
		if restarted == limit {
			return errLimit
		}
		return nil
	}, wt, d)

	if restarted == n {
		t.Fatalf("restarted count: %v, equal to count: %v", restarted, n)
	}

	if !errors.Is(gErr, errLimit) {
		t.Fatalf("restart error is not limit exceeded: %v", gErr)
	}

	restarted = 0
	// add a higher delta such that only the first 2 calls are slower
	d = 1 * time.Second
	// again let's test but now only the first 2 calls take 2 seconds, whereas the last ones are instant (no wait).
	// And let's set the fail duration to 1 second, this should make it so that it only starts failing after the second attempt, totalling 7 calls
	gErr = restartUntilErr(context.Background(), func(_ context.Context, _ bool) error {
		restarted++

		if restarted <= 2 {
			time.Sleep(2 * time.Second)
		}
		return nil
	}, wt, d)

	if restarted != n+2 {
		t.Fatalf("restarted count: %v, not equal to count plus 2: %v", restarted, n+2)
	}

	if !errors.Is(gErr, ErrMaxRestarts) {
		t.Fatalf("restart error is not max restarts: %v", gErr)
	}
}
