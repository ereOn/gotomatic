package conditional

import "testing"

func TestDereference(t *testing.T) {
	condition := NewManualCondition(true)
	weakCondition := Dereference(condition)

	weakCondition.Close()

	assertConditionState(t, condition, true, "after close")
}
