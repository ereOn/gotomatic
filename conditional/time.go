package conditional

import (
	"time"

	gtime "github.com/intelux/gotomatic/time"
)

// TimeCondition represents a condition that is met as long as the current time
// matches the specified moment.
type TimeCondition struct {
	Condition
	Moment    gtime.Moment
	timeFunc  func() time.Time
	sleepFunc func(time.Duration, <-chan struct{}) bool
	done      chan struct{}
}

func sleep(duration time.Duration, interrupt <-chan struct{}) bool {
	timer := time.NewTimer(duration)

	select {
	case <-timer.C:
		return true
	case <-interrupt:
		timer.Stop()
		return false
	}
}

// NewTimeCondition instantiates a new TimeCondition.
func NewTimeCondition(moment gtime.Moment, options ...TimeConditionOption) *TimeCondition {
	condition := &TimeCondition{
		Moment:    moment,
		timeFunc:  time.Now,
		sleepFunc: sleep,
		done:      make(chan struct{}),
	}

	for _, option := range options {
		option.apply(condition)
	}

	inMoment, _ := moment.NextInterval(condition.timeFunc())
	condition.Condition = NewManualCondition(inMoment)

	go condition.checkTime()

	return condition
}

// TimeConditionOption represents an option for a TimeCondition.
type TimeConditionOption interface {
	apply(condition *TimeCondition)
}

// TimeFunctionOption defines the time function used by a TimeCondition.
type TimeFunctionOption struct {
	TimeFunction func() time.Time
}

func (o TimeFunctionOption) apply(condition *TimeCondition) {
	condition.timeFunc = o.TimeFunction
}

// SleepFunctionOption defines the sleep function used by a TimeCondition.
type SleepFunctionOption struct {
	SleepFunction func(time.Duration, <-chan struct{}) bool
}

func (o SleepFunctionOption) apply(condition *TimeCondition) {
	condition.sleepFunc = o.SleepFunction
}

// Close terminates the condition.
//
// Any pending wait on one of the returned channels via Wait() or
// WaitChange() will be unblocked.
//
// Calling Close() twice or more has no effect.
func (condition *TimeCondition) Close() error {
	if condition.done != nil {
		close(condition.done)
		condition.done = nil
	}

	return condition.Condition.Close()
}

func (condition *TimeCondition) checkTime() {
	var delay time.Duration

	for condition.sleepFunc(delay, condition.done) {
		now := condition.timeFunc()

		inMoment, nextChange := condition.Moment.NextInterval(now)
		condition.Condition.(*ManualCondition).Set(inMoment)

		delay = nextChange.Sub(now)
	}
}
