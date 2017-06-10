// Package trigger implements triggers for conditions.
package trigger

import (
	"io"

	"github.com/intelux/gotomatic/conditional"
)

// A Trigger than watches a condition and runs when it changes.
type Trigger interface {
	run(w io.Writer, name string, state bool) error
}

// Watch a condition a runs a trigger when its state matches the specified one.
//
// The watch exits when the condition is closed or the trigger fails.
func Watch(condition conditional.Condition, trigger Trigger, w io.Writer, name string, state bool) (err error) {
	for err == nil {
		currentState, ch := condition.GetAndWaitChange()

		if currentState == state {
			err = trigger.run(w, name, state)

			if err != nil {
				break
			}
		}

		err = <-ch
	}

	return err
}
