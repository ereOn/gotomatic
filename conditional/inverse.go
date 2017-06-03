package conditional

// Inverse returns a Condition that has the reversed satisfied state as the one
// provided.
func Inverse(condition Condition) Condition {
	return &inversedCondition{Condition: condition}
}

type inversedCondition struct {
	Condition
}

// Wait returns a channel that blocks until the condition reaches the
// specified satisfied state.
//
// If the condition already has the satisfied state at the moment of the
// call, a closed channel is returned (which won't block).
func (c *inversedCondition) Wait(satisfied bool) <-chan error {
	return c.Condition.Wait(!satisfied)
}

// GetAndWaitChange returns the current satisfied state of the condition as
// well as a channel that will block until the condition state changes.
func (c *inversedCondition) GetAndWaitChange() (bool, <-chan error) {
	state, channel := c.Condition.GetAndWaitChange()

	return !state, channel
}
