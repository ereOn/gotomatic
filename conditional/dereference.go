package conditional

type unclosableCondition struct {
	Condition
}

func (unclosableCondition) Close() error { return nil }

// Dereference a condition by making any attempt to close a noop.
func Dereference(condition Condition) Condition {
	return unclosableCondition{Condition: condition}
}
