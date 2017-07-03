// Package trigger implements triggers for conditions.
package trigger

import (
	"context"
	"fmt"

	"github.com/intelux/gotomatic/conditional"
)

// A Trigger than watches a condition and runs when it changes.
type Trigger struct {
	// Up is called whenever the watched condition becomes true.
	Up Action
	// Down is called whenever the watched condition becomes false.
	Down Action
}

// Watch a condition a drive a trigger with its states changes.
//
// The watch exits when the condition is closed, the trigger fails or the
// context expires. The two first cases, return an error. The third one
// doesn't.
func Watch(ctx context.Context, condition conditional.Condition, trigger Trigger) (err error) {
	stateCh := make(chan bool, 1)
	defer close(stateCh)

	unregister := condition.Register(conditional.NewChannelObserver(stateCh))
	defer unregister()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case state := <-stateCh:
			var action Action
			var errPrefix string

			if state {
				action = trigger.Up
				errPrefix = "trigger up"
			} else {
				action = trigger.Down
				errPrefix = "trigger down"
			}

			if action != nil {
				ctx := WithConditionState(ctx, state)

				if err = action.run(ctx); err != nil {
					err = fmt.Errorf("%s: %s", errPrefix, err)
					cancel()
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
