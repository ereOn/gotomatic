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

//
//// CompositeCondition represents an aggregation of Conditions.
//type CompositeCondition struct {
//	condition     ManualCondition
//	operator      CompositeOperator
//	subconditions []Condition
//	stop          chan struct{}
//}
//
//// CompositeOperator represents an operator to use in CompositeConditions.
//type CompositeOperator interface {
//	Reduce(conditions []Condition) bool
//}
//
//var (
//	// OperatorAnd will cause the associated CompositeCondition to be satisfied
//	// only when all its sub-conditions are satisfied.
//	OperatorAnd CompositeOperator = operatorAnd{}
//
//	// OperatorOr will cause the associated CompositeCondition to be satisfied
//	// when at least one of its sub-conditions is satisfied.
//	OperatorOr CompositeOperator = operatorOr{}
//
//	// OperatorXor will cause the associated CompositeCondition to be satisfied
//	// when exactly one of its sub-conditions is satisfied.
//	OperatorXor CompositeOperator = operatorXor{}
//)
//
//// NewCompositeCondition instantiates a new CompositeCondition that uses the
//// specified operator and has the specified sub-conditions.
//func NewCompositeCondition(operator CompositeOperator, conditions ...Condition) *CompositeCondition {
//	condition := &CompositeCondition{
//		condition:     *NewManualCondition(operator.Reduce(conditions)),
//		operator:      operator,
//		subconditions: conditions,
//		stop:          make(chan struct{}),
//	}
//
//	go condition.watchConditions()
//
//	return condition
//}
//
//func (c *CompositeCondition) watchConditions() {
//	for {
//		cases := make([]reflect.SelectCase, len(c.subconditions)+1)
//		cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.stop)}
//
//		for i, condition := range c.subconditions {
//			cases[i+1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(condition.Changed())}
//		}
//
//		// ok will be true if the channel has not been closed.
//		chosen, _, _ := reflect.Select(cases)
//
//		if chosen == 0 {
//			return
//		}
//
//		c.condition.Set(c.operator.Reduce(c.subconditions))
//	}
//}
//
//// Satisfied blocks until the condition is satisfied.
//func (c *CompositeCondition) Satisfied() <-chan struct{} {
//	return c.condition.Satisfied()
//}
//
//// Unsatisfied blocks until the condition is unsatisfied.
//func (c *CompositeCondition) Unsatisfied() <-chan struct{} {
//	return c.condition.Unsatisfied()
//}
//
//// Changed blocks until the condition is satisfied state changes.
//func (c *CompositeCondition) Changed() <-chan struct{} {
//	return c.condition.Changed()
//}
//
//// Bool returns the condition current satisfied state.
//func (c *CompositeCondition) Bool() bool {
//	return c.condition.Bool()
//}
//
//// Close closes the condition.
////
//// It is not specified what happens if there are pending calls on the condition
//// while it is being closed.
//func (c *CompositeCondition) Close() error {
//	close(c.stop)
//
//	return c.condition.Close()
//}
//
//type operatorAnd struct{}
//
//func (o operatorAnd) Reduce(conditions []Condition) bool {
//	for _, condition := range conditions {
//		if !condition.Bool() {
//			return false
//		}
//	}
//
//	return true
//}
//
//type operatorOr struct{}
//
//func (o operatorOr) Reduce(conditions []Condition) bool {
//	for _, condition := range conditions {
//		if condition.Bool() {
//			return true
//		}
//	}
//
//	return false
//}
//
//type operatorXor struct{}
//
//func (o operatorXor) Reduce(conditions []Condition) bool {
//	result := false
//
//	for _, condition := range conditions {
//		if condition.Bool() {
//			if result {
//				return false
//			} else {
//				result = true
//			}
//		}
//	}
//
//	return result
//}
