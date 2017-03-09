package conditional

import "testing"

func TestCompositeConditionOperatorAnd(t *testing.T) {
	a := NewManualCondition(false)
	b := NewManualCondition(false)

	condition := NewCompositeCondition(OperatorAnd, a, b)
	defer condition.Close()
	assertConditionUnsatisfied(t, condition, "initialization to: and, (false, false)")

	a.Set(true)
	assertConditionUnsatisfied(t, condition, "a set to true")

	b.Set(true)
	assertConditionUnsatisfied(t, condition, "b set to true")

	condition = NewCompositeCondition(OperatorAnd, a, b)
	defer condition.Close()
	assertConditionSatisfied(t, condition, "initialization to: and, (true, true)")

	a.Set(false)
	//TODO: There is a race condition here.
	//<-condition.Wait(false)
	assertConditionUnsatisfied(t, condition, "a set to false")

	assertConditionChanged(t, condition, "a set to true again", false, func() { a.Set(true) })
}
