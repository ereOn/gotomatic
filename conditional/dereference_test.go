package conditional

import "testing"

func TestDereference(t *testing.T) {
	condition := Inverse(NewManualCondition(true))
	weakCondition := Dereference(condition)

	weakCondition.Close()

	assertConditionState(t, condition, false, "after close")
}

func TestDereferenceSettable(t *testing.T) {
	condition := NewManualCondition(true)
	weakCondition := Dereference(condition)

	weakCondition.Close()

	assertConditionState(t, condition, true, "after close")
}
