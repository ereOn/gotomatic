package conditional

import (
	"context"
	"testing"
	"time"
)

func assertConditionSatisfied(t *testing.T, condition Condition, ctx string) {
	select {
	case <-condition.Wait(true):
	case <-condition.Wait(false):
		t.Errorf("condition should be satisfied after %s", ctx)
	default:
		t.Errorf("the condition cannot be both in satisfied and unsatisfied states after %s", ctx)
	}
}

func assertConditionUnsatisfied(t *testing.T, condition Condition, ctx string) {
	select {
	case <-condition.Wait(true):
		t.Errorf("condition should not be satisfied after %s", ctx)
	case <-condition.Wait(false):
	default:
		t.Errorf("the condition cannot be both in satisfied and unsatisfied states after %s", ctx)
	}
}

func assertConditionChanged(t *testing.T, condition Condition, ctx string, satisfied bool, f func()) {
	state, channel := condition.GetAndWaitChange()

	if state != satisfied {
		t.Errorf("condition satisfied state was %v but %v was expected after %s", state, satisfied, ctx)
	}

	timeout, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	result := make(chan bool)

	go func() {
		select {
		case <-channel:
			result <- true
		case <-timeout.Done():
			result <- false
		}
	}()

	f()

	if !<-result {
		t.Errorf("condition should have changed after %s", ctx)
	}
}
