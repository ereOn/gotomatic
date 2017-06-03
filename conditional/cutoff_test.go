package conditional

import (
	"context"
	"testing"
	"time"
)

func TestCutOffConditionZeroThresholds(t *testing.T) {
	ch := make(chan bool, 1)
	defer close(ch)
	ch <- true

	condition := NewCutOffCondition(0, 0, time.Millisecond, func(context.Context) bool { return <-ch })
	defer condition.Close()

	assertConditionState(t, condition, true, "initialization to true")
	assertConditionChanged(t, condition, true, "callback returns false", func() { ch <- false })
	assertConditionChanged(t, condition, false, "callback returns true", func() { ch <- true })
}

func TestCutOffConditionPositiveThresholds(t *testing.T) {
	ch := make(chan bool, 1)
	defer close(ch)
	ch <- false

	condition := NewCutOffCondition(2, 3, time.Millisecond, func(context.Context) bool { return <-ch })
	defer condition.Close()

	assertConditionState(t, condition, false, "initialization to false")
	ch <- true
	assertConditionState(t, condition, false, "first set to true")
	ch <- true
	assertConditionChanged(t, condition, false, "callback returns true", func() { ch <- true })
	ch <- false
	assertConditionState(t, condition, true, "first set to false")
	ch <- false
	assertConditionState(t, condition, true, "second set to false")
	ch <- false
	assertConditionChanged(t, condition, true, "callback returns false", func() { ch <- false })
}
