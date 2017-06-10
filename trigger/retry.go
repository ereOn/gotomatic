package trigger

import (
	"io"
	"time"
)

type retryTrigger struct {
	Trigger
	max   int
	delay time.Duration
}

// Retry returns a trigger that, upon failure, retries the specified number of
// times, with the specified delay between attempts.
//
// If max is 0 or less, the trigger gets never called, and thus never fails.
// This is kind-of useless.
func Retry(trigger Trigger, max int, delay time.Duration) Trigger {
	return retryTrigger{
		Trigger: trigger,
		max:     max,
		delay:   delay,
	}
}

func (t retryTrigger) run(w io.Writer, name string, state bool) (err error) {
	for i := 0; i < t.max; i++ {
		if err = t.Trigger.run(w, name, state); err == nil {
			return
		}

		time.Sleep(t.delay)
	}

	return
}
