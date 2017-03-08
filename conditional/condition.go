package conditional

// Condition represents a condition that can be satisfied.
type Condition interface {
	Satisfied() <-chan struct{}
	Unsatisfied() <-chan struct{}
}

// ManualCondition is a condition that can be set or unset explicitely.
type ManualCondition struct {
	satisfied   chan struct{}
	unsatisfied chan struct{}
}

// NewManualCondition instantiates a new ManualCondition in the specified
// initial state satisfied.
func NewManualCondition(satisfied bool) *ManualCondition {
	condition := &ManualCondition{
		satisfied:   make(chan struct{}),
		unsatisfied: make(chan struct{}),
	}

	if satisfied {
		close(condition.satisfied)
	} else {
		close(condition.unsatisfied)
	}

	return condition
}

// Satisfied blocks until the condition is satisfied.
func (c *ManualCondition) Satisfied() <-chan struct{} {
	return c.satisfied
}

// Unsatisfied blocks until the condition is unsatisfied.
func (c *ManualCondition) Unsatisfied() <-chan struct{} {
	return c.unsatisfied
}

// Set the condition satisfied state explicitely.
func (c *ManualCondition) Set(satisfied bool) {
	if satisfied {
		select {
		case <-c.satisfied:
			return
		case <-c.unsatisfied:
			close(c.satisfied)
			c.unsatisfied = make(chan struct{})
		}
	} else {
		select {
		case <-c.unsatisfied:
			return
		case <-c.satisfied:
			close(c.unsatisfied)
			c.satisfied = make(chan struct{})
		}
	}
}
