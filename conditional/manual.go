package conditional

import "sync"

// ManualCondition is a condition that can be set or unset explicitely.
type ManualCondition struct {
	lock      sync.Mutex
	satisfied bool
	channels  []chan struct{}
}

// NewManualCondition instantiates a new ManualCondition in the specified
// initial state satisfied.
func NewManualCondition(satisfied bool) *ManualCondition {
	return &ManualCondition{
		satisfied: satisfied,
		channels:  make([]chan struct{}, 0, 0),
	}
}

// Wait returns a channel that blocks until the condition reaches the
// specified satisfied state.
//
// If the condition already has the satisfied state at the moment of the
// call, a closed channel is returned (which won't block).
func (c *ManualCondition) Wait(satisfied bool) <-chan struct{} {
	channel := make(chan struct{})

	c.lock.Lock()
	defer c.lock.Unlock()

	if satisfied == c.satisfied {
		close(channel)
	} else {
		c.channels = append(c.channels, channel)
	}

	return channel
}

// GetAndWaitChange returns the current satisfied state of the condition as
// well as a channel that will block until the condition state changes.
func (c *ManualCondition) GetAndWaitChange() (bool, <-chan struct{}) {
	channel := make(chan struct{})

	c.lock.Lock()
	defer c.lock.Unlock()

	c.channels = append(c.channels, channel)

	return c.satisfied, channel
}

// Close terminates the condition.
//
// Any pending wait on one of the returned channels via Wait() or
// WaitChange() will be unblocked.
//
// Calling Close() twice or more has no effect.
func (c *ManualCondition) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.unprotectedClose()

	return nil
}

// Set defines the ManualCondition satisfied state explicitely.
//
// Setting the condition to its current state is a no-op and does not unblock
// any previously returned channel.
func (c *ManualCondition) Set(satisfied bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if satisfied != c.satisfied {
		c.satisfied = satisfied

		// This is an implementation detail but it just happens that closing
		// the condition will have no undesirable side-effect yet does exactly
		// what we want.
		//
		// The condition is *NOT* to be considered closed after this. It's
		// magic.
		c.unprotectedClose()
	}
}

func (c *ManualCondition) unprotectedClose() {
	for _, channel := range c.channels {
		close(channel)
	}

	c.channels = make([]chan struct{}, 0, 0)
}
