package conditional

type unclosableCondition struct {
	Condition
}

func (unclosableCondition) Close() error { return nil }

// Dereference a condition by making its Close() function a no-op.
func Dereference(condition Condition) Condition {
	return unclosableCondition{Condition: condition}
}
