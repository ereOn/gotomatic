// Package conditional defines condition primitives and logic to ease writing
// of complex condition compositions.
package conditional

import "errors"

// ErrConditionClosed is the error returned when a wait on a condition is
// interrupted because the channel was closed.
var ErrConditionClosed = errors.New("condition was closed")

// A Condition that can be either satisfied or unsatisfied.
//
// All methods on a Condition are thread-safe.
type Condition interface {
	// Wait returns a channel that blocks until the condition reaches the
	// specified satisfied state.
	//
	// If the condition already has the satisfied state at the moment of the
	// call, a closed channel is returned (which won't block).
	//
	// If the condition is closed or the wait fails for whatever reason,
	// `ErrConditionClosed` is returned on the channel.
	Wait(satisfied bool) <-chan error

	// GetAndWaitChange returns the current satisfied state of the condition as
	// well as a channel that will block until the condition state changes.
	//
	// If the condition is closed or the wait fails for whatever reason,
	// `ErrConditionClosed` is returned on the channel.
	GetAndWaitChange() (bool, <-chan error)

	// Close terminates the condition.
	//
	// Any pending wait on one of the returned channels via Wait() or
	// WaitChange() will be unblocked and `ErrConditionClosed` put in the wait
	// channels.
	//
	// Calling Close() twice or more has no effect.
	Close() error

	// Register an observer for changes.
	//
	// Any change will cause the following observer to be called with the
	// current state until the returned cancel function is called.
	Register(ConditionStateObserver) func()
}
