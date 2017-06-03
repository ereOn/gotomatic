package conditional

import (
	"context"
	"time"
)

type cutOffCondition struct {
	Condition
	upThreshold   uint
	downThreshold uint
	period        time.Duration
	executor      Executor
	done          chan struct{}
	counter       uint
	lastState     bool
	locked        bool
}

// NewCutOffCondition creates a new cut-off condition.
//
// The condition is set when the specified executor returns true for
// `upThreshold` periods of time.
//
// The condition is unset when the specified executor returns false for
// `downThreshold` periods of time.
//
// Any change of the executor return value resets both counters.
//
// The inital status of the condition is the return value of the executor.
func NewCutOffCondition(upThreshold uint, downThreshold uint, period time.Duration, executor Executor) Condition {
	ctx, cancel := context.WithTimeout(context.Background(), period)
	state := executor(ctx)
	cancel()

	condition := &cutOffCondition{
		Condition:     NewManualCondition(state),
		upThreshold:   upThreshold,
		downThreshold: downThreshold,
		period:        period,
		executor:      executor,
		done:          make(chan struct{}),
		lastState:     state,
		locked:        true,
	}

	go condition.run(condition.done)

	return condition
}

func (c *cutOffCondition) run(done <-chan struct{}) {
	ticker := time.NewTicker(c.period)

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), c.period)
			state := c.executor(ctx)
			cancel()
			c.tick(state)
		}
	}
}

func (c *cutOffCondition) tick(state bool) {
	if state == c.lastState {
		c.increment()
	} else {
		c.reset(state)
	}

	c.check(state)
}

func (c *cutOffCondition) increment() {
	if !c.locked {
		c.counter++
	}
}

func (c *cutOffCondition) check(state bool) {
	if state {
		if c.counter >= c.upThreshold {
			c.set(state)
		}
	} else {
		if c.counter >= c.downThreshold {
			c.set(state)
		}
	}
}

func (c *cutOffCondition) set(state bool) {
	if !c.locked {
		c.locked = true
		c.Condition.(*ManualCondition).Set(state)
	}
}

func (c *cutOffCondition) reset(state bool) {
	c.lastState = state
	c.locked = false
	c.counter = 0
}

func (c *cutOffCondition) Close() error {
	close(c.done)

	return c.Condition.Close()
}
