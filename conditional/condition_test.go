package conditional

import (
	"context"
	"testing"
	"time"
)

func waitChannel(channel <-chan error) chan bool {
	timeout, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))

	result := make(chan bool)

	go func() {
		select {
		case <-channel:
			cancel()
			result <- true
		case <-timeout.Done():
			result <- false
		}
	}()

	return result
}

func assertConditionState(t *testing.T, condition Condition, satisfied bool, ctx string) {
	if !<-waitChannel(condition.Wait(satisfied)) {
		if satisfied {
			t.Fatalf("condition should be satisfied after %s", ctx)
		} else {
			t.Fatalf("condition should not be satisfied after %s", ctx)
		}
	}
}

func assertConditionChanged(t *testing.T, condition Condition, satisfied bool, ctx string, f func()) {
	assertConditionState(t, condition, satisfied, ctx)

	state, channel := condition.GetAndWaitChange()

	if state != satisfied {
		t.Errorf("condition satisfied state was %v but %v was expected after %s", state, satisfied, ctx)
	}

	result := waitChannel(channel)

	f()

	if !<-result {
		t.Errorf("condition should have changed after %s", ctx)
	}
}

func assertCloseCondition(t *testing.T, condition Condition) {
	_, channel := condition.GetAndWaitChange()

	condition.Close()

	err := <-channel

	if err != ErrConditionClosed {
		t.Errorf("condition channel returned an unexpected error: %s", err)
	}
}
