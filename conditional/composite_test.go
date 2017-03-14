package conditional

import "testing"

func TestCompositeConditionNoSubconditions(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("instantiation was supposed to panic")
		}
	}()

	NewCompositeCondition(OperatorAnd)
}

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

func TestCompositeConditionOperatorOr(t *testing.T) {
	a := NewManualCondition(false)
	b := NewManualCondition(false)

	condition := NewCompositeCondition(OperatorOr, a, b)
	defer condition.Close()
	assertConditionState(t, condition, false, "initialization to: and, (false, false)")

	a.Set(true)
	assertConditionState(t, condition, true, "a set to true")

	b.Set(true)
	assertConditionState(t, condition, true, "b set to true")

	condition = NewCompositeCondition(OperatorOr, a, b)
	defer condition.Close()
	assertConditionState(t, condition, true, "initialization to: and, (true, true)")

	a.Set(false)
	assertConditionState(t, condition, true, "a set to false")

	b.Set(false)
	assertConditionState(t, condition, false, "b set to false")

	assertConditionChanged(t, condition, false, "a set to true again", func() { a.Set(true) })
}

func TestCompositeConditionOperatorXor(t *testing.T) {
	a := NewManualCondition(false)
	b := NewManualCondition(false)

	condition := NewCompositeCondition(OperatorXor, a, b)
	defer condition.Close()
	assertConditionState(t, condition, false, "initialization to: and, (false, false)")

	a.Set(true)
	assertConditionState(t, condition, true, "a set to true")

	b.Set(true)
	assertConditionState(t, condition, false, "b set to true")

	condition = NewCompositeCondition(OperatorXor, a, b)
	defer condition.Close()
	assertConditionState(t, condition, false, "initialization to: and, (true, true)")

	a.Set(false)
	assertConditionState(t, condition, true, "a set to false")

	b.Set(false)
	assertConditionState(t, condition, false, "b set to false")

	assertConditionChanged(t, condition, false, "a set to true again", func() { a.Set(true) })
}

func TestOperatorAndString(t *testing.T) {
	value := OperatorAnd.String()
	expected := "and"

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestOperatorOrString(t *testing.T) {
	value := OperatorOr.String()
	expected := "or"

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestOperatorXorString(t *testing.T) {
	value := OperatorXor.String()
	expected := "xor"

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}
