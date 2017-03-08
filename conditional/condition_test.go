package conditional

import "testing"

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

func TestManualCondition(t *testing.T) {
	condition := NewManualCondition(false)
	assertConditionUnsatisfied(t, condition, "initialization to false")

	condition = NewManualCondition(true)
	assertConditionSatisfied(t, condition, "initialization to true")

	condition.Set(true)
	assertConditionSatisfied(t, condition, "first set to true")

	condition.Set(false)
	assertConditionUnsatisfied(t, condition, "second set to false")

	condition.Set(false)
	assertConditionUnsatisfied(t, condition, "third set to false")

	condition.Set(true)
	assertConditionSatisfied(t, condition, "fourth set to true")
}
