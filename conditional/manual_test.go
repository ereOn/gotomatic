package conditional

import (
	"testing"
)

func TestManualCondition(t *testing.T) {
	condition := NewManualCondition(false)
	defer assertCloseCondition(t, condition)
	assertConditionState(t, condition, false, "initialization to false")

	condition = NewManualCondition(true)
	defer assertCloseCondition(t, condition)
	assertConditionState(t, condition, true, "initialization to true")

	ch := make(chan bool, 10)
	defer close(ch)
	observer := NewChannelObserver(ch)
	unregister := condition.Register(observer)
	defer unregister()

	condition.Set(true)
	assertConditionState(t, condition, true, "first set to true")

	condition.Set(false)
	assertConditionState(t, condition, false, "second set to false")

	condition.Set(false)
	assertConditionState(t, condition, false, "third set to false")

	condition.Set(true)
	assertConditionState(t, condition, true, "fourth set to true")
	assertConditionChanged(t, condition, true, "fifth set to false", func() { condition.Set(false) })
}
