package conditional

import "reflect"

// CompositeCondition represents an aggregation of Conditions.
type CompositeCondition struct {
	condition     ManualCondition
	operator      CompositeOperator
	subconditions []Condition
	stop          chan struct{}
}

// CompositeOperator represents an operator to use in CompositeConditions.
type CompositeOperator interface {
	Reduce(values []bool) bool
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
	if len(conditions) == 0 {
		panic("cannot instantiate a composite condition without at least one sub-condition")
	}

	condition := &CompositeCondition{
		condition:     *NewManualCondition(false),
		operator:      operator,
		subconditions: conditions,
		stop:          make(chan struct{}),
	}

	ready := make(chan struct{})
	go condition.watchConditions(ready)
	<-ready

	return condition
}

// Wait returns a channel that blocks until the condition reaches the
// specified satisfied state.
//
// If the condition already has the satisfied state at the moment of the
// call, a closed channel is returned (which won't block).
func (c *CompositeCondition) Wait(satisfied bool) <-chan struct{} {
	return c.condition.Wait(satisfied)
}

// GetAndWaitChange returns the current satisfied state of the condition as
// well as a channel that will block until the condition state changes.
func (c *CompositeCondition) GetAndWaitChange() (bool, <-chan struct{}) {
	return c.condition.GetAndWaitChange()
}

// Close terminates the condition.
//
// Any pending wait on one of the returned channels via Wait() or
// WaitChange() will be unblocked.
//
// Calling Close() twice or more has no effect.
func (c *CompositeCondition) Close() error {
	select {
	case <-c.stop:
	default:
		close(c.stop)
	}

	for _, condition := range c.subconditions {
		condition.Close()
	}

	return c.condition.Close()
}

func (c *CompositeCondition) watchConditions(ready chan struct{}) {
	for {
		count := len(c.subconditions)
		values := make([]bool, count, count)
		cases := make([]reflect.SelectCase, count+1)

		for i, condition := range c.subconditions {
			value, channel := condition.GetAndWaitChange()
			values[i] = value
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)}
		}

		cases[count] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.stop)}
		c.condition.Set(c.operator.Reduce(values))
		close(ready)

		for {
			chosen, _, _ := reflect.Select(cases)

			if chosen == count {
				return
			}

			condition := c.subconditions[chosen]
			value, channel := condition.GetAndWaitChange()
			values[chosen] = value
			cases[chosen] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)}
			c.condition.Set(c.operator.Reduce(values))
		}
	}
}

type operatorAnd struct{}

func (o operatorAnd) Reduce(values []bool) bool {
	for _, value := range values {
		if !value {
			return false
		}
	}

	return true
}

type operatorOr struct{}

func (o operatorOr) Reduce(values []bool) bool {
	for _, value := range values {
		if value {
			return true
		}
	}

	return false
}

type operatorXor struct{}

func (o operatorXor) Reduce(values []bool) bool {
	result := false

	for _, value := range values {
		if value {
			if result {
				return false
			}

			result = true
		}
	}

	return result
}
