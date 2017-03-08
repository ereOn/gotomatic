package conditional

import (
	"reflect"
	"sync"
)

// Condition represents a condition that can be satisfied.
type Condition interface {
	Satisfied() <-chan struct{}
	Unsatisfied() <-chan struct{}
	Changed() <-chan struct{}
	Bool() bool
	Close() error
}

// ManualCondition is a condition that can be set or unset explicitely.
type ManualCondition struct {
	lock        sync.RWMutex
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
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.satisfied
}

// Unsatisfied blocks until the condition is unsatisfied.
func (c *ManualCondition) Unsatisfied() <-chan struct{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.unsatisfied
}

// Changed blocks until the condition is satisfied state changes.
func (c *ManualCondition) Changed() <-chan struct{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	select {
	case <-c.satisfied:
		return c.unsatisfied
	case <-c.unsatisfied:
		return c.satisfied
	}
}

// Bool returns the condition current satisfied state.
func (c *ManualCondition) Bool() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	select {
	case <-c.satisfied:
		return true
	case <-c.unsatisfied:
		return false
	}
}

// Close closes the condition.
//
// It is not specified what happens if there are pending calls on the condition
// while it is being closed.
func (c *ManualCondition) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	select {
	case <-c.satisfied:
		close(c.unsatisfied)
	case <-c.unsatisfied:
		close(c.satisfied)
	}

	return nil
}

// Set the condition satisfied state explicitely.
func (c *ManualCondition) Set(satisfied bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

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

// CompositeCondition represents an aggregation of Conditions.
type CompositeCondition struct {
	condition     ManualCondition
	operator      CompositeOperator
	subconditions []Condition
	stop          chan struct{}
}

// CompositeOperator represents an operator to use in CompositeConditions.
type CompositeOperator interface {
	Reduce(conditions []Condition) bool
}

var (
	// OperatorAnd will cause the associated CompositeCondition to be satisfied
	// only when all its sub-conditions are satisfied.
	OperatorAnd CompositeOperator = operatorAnd{}

	// OperatorOr will cause the associated CompositeCondition to be satisfied
	// when at least one of its sub-conditions is satisfied.
	OperatorOr CompositeOperator = operatorOr{}

	// OperatorXor will cause the associated CompositeCondition to be satisfied
	// when exactly one of its sub-conditions is satisfied.
	OperatorXor CompositeOperator = operatorXor{}
)

// NewCompositeCondition instantiates a new CompositeCondition that uses the
// specified operator and has the specified sub-conditions.
func NewCompositeCondition(operator CompositeOperator, conditions ...Condition) *CompositeCondition {
	condition := &CompositeCondition{
		condition:     *NewManualCondition(operator.Reduce(conditions)),
		operator:      operator,
		subconditions: conditions,
		stop:          make(chan struct{}),
	}

	go condition.watchConditions()

	return condition
}

func (c *CompositeCondition) watchConditions() {
	for {
		cases := make([]reflect.SelectCase, len(c.subconditions)+1)
		cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.stop)}

		for i, condition := range c.subconditions {
			cases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(condition.Changed())}
		}

		// ok will be true if the channel has not been closed.
		chosen, _, _ := reflect.Select(cases)

		if chosen == 0 {
			return
		}

		c.condition.Set(c.operator.Reduce(c.subconditions))
	}
}

// Satisfied blocks until the condition is satisfied.
func (c *CompositeCondition) Satisfied() <-chan struct{} {
	return c.condition.Satisfied()
}

// Unsatisfied blocks until the condition is unsatisfied.
func (c *CompositeCondition) Unsatisfied() <-chan struct{} {
	return c.condition.Unsatisfied()
}

// Changed blocks until the condition is satisfied state changes.
func (c *CompositeCondition) Changed() <-chan struct{} {
	return c.condition.Changed()
}

// Bool returns the condition current satisfied state.
func (c *CompositeCondition) Bool() bool {
	return c.condition.Bool()
}

// Close closes the condition.
//
// It is not specified what happens if there are pending calls on the condition
// while it is being closed.
func (c *CompositeCondition) Close() error {
	close(c.stop)

	return c.condition.Close()
}

type operatorAnd struct{}

func (o operatorAnd) Reduce(conditions []Condition) bool {
	for _, condition := range conditions {
		if !condition.Bool() {
			return false
		}
	}

	return true
}

type operatorOr struct{}

func (o operatorOr) Reduce(conditions []Condition) bool {
	for _, condition := range conditions {
		if condition.Bool() {
			return true
		}
	}

	return false
}

type operatorXor struct{}

func (o operatorXor) Reduce(conditions []Condition) bool {
	result := false

	for _, condition := range conditions {
		if condition.Bool() {
			if result {
				return false
			} else {
				result = true
			}
		}
	}

	return result
}
