package conditional

import "testing"

func TestInverse(t *testing.T) {
	m := NewManualCondition(false)
	condition := Inverse(m)
	defer condition.Close()
	assertConditionState(t, condition, true, "initialization to false")

	m.Set(true)
	assertConditionState(t, condition, false, "first set to true")

	m.Set(false)
	assertConditionState(t, condition, true, "second set to false")

	m.Set(false)
	assertConditionState(t, condition, true, "third set to false")

	m.Set(true)
	assertConditionState(t, condition, false, "fourth set to true")
	assertConditionChanged(t, condition, false, "fifth set to false", func() { m.Set(false) })
}
