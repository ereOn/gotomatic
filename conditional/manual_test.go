package conditional

import "testing"

func TestManualCondition(t *testing.T) {
	condition := NewManualCondition(false)
	defer condition.Close()
	assertConditionUnsatisfied(t, condition, "initialization to false")

	condition = NewManualCondition(true)
	defer condition.Close()
	assertConditionSatisfied(t, condition, "initialization to true")

	condition.Set(true)
	assertConditionSatisfied(t, condition, "first set to true")

	condition.Set(false)
	assertConditionUnsatisfied(t, condition, "second set to false")

	condition.Set(false)
	assertConditionUnsatisfied(t, condition, "third set to false")

	condition.Set(true)
	assertConditionSatisfied(t, condition, "fourth set to true")
	assertConditionChanged(t, condition, "fifth set to false", true, func() { condition.Set(false) })
}
