package trigger

import (
	"context"
	"time"
)

type retryAction struct {
	Action
	max   int
	delay time.Duration
}

// Retry returns an action that, upon failure, retries the specified number of
// times, with the specified delay between attempts.
//
// If max is 0 or less, the action gets never called, and thus never fails.
// This is kind-of useless.
func Retry(action Action, max int, delay time.Duration) Action {
	return retryAction{
		Action: action,
		max:    max,
		delay:  delay,
	}
}

func (t retryAction) run(ctx context.Context) (err error) {
	for i := 0; i < t.max; i++ {
		if err = t.Action.run(ctx); err == nil {
			return
		}

		time.Sleep(t.delay)
	}

	return
}
