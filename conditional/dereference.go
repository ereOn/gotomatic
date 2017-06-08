package conditional

type unclosableCondition struct {
	Condition
}

func (unclosableCondition) Close() error { return nil }

type unclosableSettableCondition struct {
	Condition
	Settable
}

// Dereference a condition by making its Close() function a no-op.
func Dereference(condition Condition) Condition {
	if settable, ok := condition.(Settable); ok {
		return unclosableSettableCondition{Condition: unclosableCondition{Condition: condition}, Settable: settable}
	}

	return unclosableCondition{Condition: condition}
}
