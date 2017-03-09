package conditional

import "testing"

func TestCompositeConditionOperatorAnd(t *testing.T) {
	a := NewManualCondition(false)
	b := NewManualCondition(false)

	condition := NewCompositeCondition(OperatorAnd, a, b)
	defer condition.Close()
	assertConditionState(t, condition, false, "initialization to: and, (false, false)")

	a.Set(true)
	assertConditionState(t, condition, false, "a set to true")

	b.Set(true)
	assertConditionState(t, condition, true, "b set to true")

	condition = NewCompositeCondition(OperatorAnd, a, b)
	defer condition.Close()
	assertConditionState(t, condition, true, "initialization to: and, (true, true)")

	a.Set(false)
	assertConditionState(t, condition, false, "a set to false")

	assertConditionChanged(t, condition, false, "a set to true again", func() { a.Set(true) })
}
