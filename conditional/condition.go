package conditional

// Condition represents a condition that can be satisfied.
//
// All methods on a Condition are thread-safe.
type Condition interface {
	// Wait returns a channel that blocks until the condition reaches the
	// specified satisfied state.
	//
	// If the condition already has the satisfied state at the moment of the
	// call, a closed channel is returned (which won't block).
	Wait(satisfied bool) <-chan struct{}

	// GetAndWaitChange returns the current satisfied state of the condition as
	// well as a channel that will block until the condition state changes.
	GetAndWaitChange() (bool, <-chan struct{})

	// Close terminates the condition.
	//
	// Any pending wait on one of the returned channels via Wait() or
	// WaitChange() will be unblocked.
	//
	// Calling Close() twice or more has no effect.
	Close() error
}
