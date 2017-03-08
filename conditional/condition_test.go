package conditional

import (
	"context"
	"testing"
	"time"
)

func assertConditionSatisfied(t *testing.T, condition Condition, ctx string) {
	select {
	case <-condition.Satisfied():
	case <-condition.Unsatisfied():
		t.Errorf("condition should be satisfied after %s", ctx)
	default:
		t.Errorf("the condition cannot be both in satisfied and unsatisfied states after %s", ctx)
	}
}

func assertConditionUnsatisfied(t *testing.T, condition Condition, ctx string) {
	select {
	case <-condition.Satisfied():
		t.Errorf("condition should not be satisfied after %s", ctx)
	case <-condition.Unsatisfied():
	default:
		t.Errorf("the condition cannot be both in satisfied and unsatisfied states after %s", ctx)
	}
}

func assertConditionUnchanged(t *testing.T, condition Condition, ctx string) {
	select {
	case <-condition.Changed():
		t.Errorf("condition should not have changed after %s", ctx)
	default:
	}
}

func assertConditionChanged(t *testing.T, condition Condition, ctx string, f func()) {
	timeout, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*3))
	defer cancel()
	changed := condition.Changed()
	result := make(chan bool)

	go func() {
		select {
		case <-changed:
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

func TestManualCondition(t *testing.T) {
	condition := NewManualCondition(false)
	defer condition.Close()
	assertConditionUnsatisfied(t, condition, "initialization to false")
	assertConditionUnchanged(t, condition, "initialization to false")

	condition = NewManualCondition(true)
	defer condition.Close()
	assertConditionSatisfied(t, condition, "initialization to true")
	assertConditionUnchanged(t, condition, "initialization to true")

	condition.Set(true)
	assertConditionSatisfied(t, condition, "first set to true")

	condition.Set(false)
	assertConditionUnsatisfied(t, condition, "second set to false")

	condition.Set(false)
	assertConditionUnsatisfied(t, condition, "third set to false")

	condition.Set(true)
	assertConditionSatisfied(t, condition, "fourth set to true")

	assertConditionChanged(t, condition, "fifth set to false", func() { condition.Set(false) })
	assertConditionChanged(t, condition, "sixth set to true", func() { condition.Set(true) })
}

func TestEmptyCompositeConditionAnd(t *testing.T) {
	condition := NewCompositeCondition(OperatorAnd)
	defer condition.Close()
	assertConditionSatisfied(t, condition, "empty initialization")
}

func TestCompositeConditionAnd(t *testing.T) {
	a := NewManualCondition(false)
	defer a.Close()
	b := NewManualCondition(false)
	defer b.Close()
	condition := NewCompositeCondition(OperatorAnd, a, b)
	defer condition.Close()
	assertConditionUnsatisfied(t, condition, "empty initialization")

	a.Set(true)
	assertConditionUnsatisfied(t, condition, "a was set")

	b.Set(true)
	time.Sleep(time.Second)
	assertConditionSatisfied(t, condition, "b was set")

	a.Set(false)
	assertConditionUnsatisfied(t, condition, "a was unset")
}
