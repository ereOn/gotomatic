package conditional

import (
	"time"
)

// TimeCondition represents a condition that is met as long as the current time
// matches the specified moment.
type TimeCondition struct {
	Condition
	Moment    Moment
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
func NewTimeCondition(moment Moment, options ...TimeConditionOption) *TimeCondition {
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

// Frequency represents a frequency.
type Frequency interface {
	getBase(start time.Time, t time.Time) time.Time
	Previous(start time.Time, t time.Time) time.Time
	Next(start time.Time, t time.Time) time.Time
}

var (
	// FrequencyYear represents a moment that happens every year.
	FrequencyYear = frequencyYear{}
	// FrequencyMonth represents a moment that happens every month.
	FrequencyMonth = frequencyMonth{}
	// FrequencyWeek represents a moment that happens every week.
	FrequencyWeek = frequencyWeek{}
	// FrequencyDay represents a moment that happens every day.
	FrequencyDay = frequencyDay{}
	// FrequencyHour represents a moment that happens every hour.
	FrequencyHour = frequencyHour{}
	// FrequencyMinute represents a moment that happens every minute.
	FrequencyMinute = frequencyMinute{}
	// FrequencySecond represents a moment that happens every second.
	FrequencySecond = frequencySecond{}
)

type frequencyYear struct{}

func (frequencyYear) getBase(start time.Time, t time.Time) time.Time {
	t = t.In(start.Location())

	_, month, day := start.Date()
	hour, minute, second := start.Clock()
	return time.Date(t.Year(), month, day, hour, minute, second, start.Nanosecond(), start.Location())
}

func (f frequencyYear) Previous(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if r.After(t) {
		r = r.AddDate(-1, 0, 0)
	}

	return r
}

func (f frequencyYear) Next(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if !r.After(t) {
		r = r.AddDate(1, 0, 0)
	}

	return r
}

type frequencyMonth struct{}

func (frequencyMonth) getBase(start time.Time, t time.Time) time.Time {
	t = t.In(start.Location())

	_, _, day := start.Date()
	hour, minute, second := start.Clock()

	return time.Date(t.Year(), t.Month(), day, hour, minute, second, start.Nanosecond(), start.Location())
}

func (f frequencyMonth) Previous(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if r.After(t) {
		r = r.AddDate(0, -1, 0)
	}

	return r
}

func (f frequencyMonth) Next(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if !r.After(t) {
		r = r.AddDate(0, 1, 0)
	}

	return r
}

type frequencyWeek struct{}

func (frequencyWeek) getBase(start time.Time, t time.Time) time.Time {
	t = t.In(start.Location())

	weekday := start.Weekday()
	hour, minute, second := start.Clock()

	return time.Date(t.Year(), t.Month(), t.Day()+int(weekday)-int(t.Weekday()), hour, minute, second, start.Nanosecond(), start.Location())
}

func (f frequencyWeek) Previous(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if r.After(t) {
		r = r.AddDate(0, 0, -7)
	}

	return r
}
func (f frequencyWeek) Next(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if !r.After(t) {
		r = r.AddDate(0, 0, 7)
	}

	return r
}

type frequencyDay struct{}

func (frequencyDay) getBase(start time.Time, t time.Time) time.Time {
	t = t.In(start.Location())

	hour, minute, second := start.Clock()

	return time.Date(t.Year(), t.Month(), t.Day(), hour, minute, second, start.Nanosecond(), start.Location())
}

func (f frequencyDay) Previous(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if r.After(t) {
		r = r.AddDate(0, 0, -1)
	}

	return r
}

func (f frequencyDay) Next(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if !r.After(t) {
		r = r.AddDate(0, 0, 1)
	}

	return r
}

type frequencyHour struct{}

func (frequencyHour) getBase(start time.Time, t time.Time) time.Time {
	t = t.In(start.Location())

	_, minute, second := start.Clock()

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, second, start.Nanosecond(), start.Location())
}

func (f frequencyHour) Previous(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if r.After(t) {
		r = r.Add(-time.Hour)
	}

	return r
}

func (f frequencyHour) Next(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if !r.After(t) {
		r = r.Add(time.Hour)
	}

	return r
}

type frequencyMinute struct{}

func (frequencyMinute) getBase(start time.Time, t time.Time) time.Time {
	t = t.In(start.Location())

	_, _, second := start.Clock()

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), second, start.Nanosecond(), start.Location())
}

func (f frequencyMinute) Previous(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if r.After(t) {
		r = r.Add(-time.Minute)
	}

	return r
}

func (f frequencyMinute) Next(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if !r.After(t) {
		r = r.Add(time.Minute)
	}

	return r
}

type frequencySecond struct{}

func (frequencySecond) getBase(start time.Time, t time.Time) time.Time {
	t = t.In(start.Location())

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), start.Nanosecond(), start.Location())
}

func (f frequencySecond) Previous(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if r.After(t) {
		r = r.Add(-time.Second)
	}

	return r
}

func (f frequencySecond) Next(start time.Time, t time.Time) time.Time {
	r := f.getBase(start, t)

	if !r.After(t) {
		r = r.Add(time.Second)
	}

	return r
}

// Moment represents a moment in time.
type Moment interface {
	NextInterval(time.Time) (bool, time.Time)
}

type recurrentMoment struct {
	Start     time.Time
	Stop      time.Time
	Frequency Frequency
}

// NewRecurrentMoment instantiates a new recurrent moment.
func NewRecurrentMoment(start, stop time.Time, frequency Frequency) Moment {
	return recurrentMoment{
		Start:     start,
		Stop:      stop,
		Frequency: frequency,
	}
}

// NextInterval returns a boolean flag that indicates whether the specified
// time is within the interval, and the time of the next interval boundary.
func (r recurrentMoment) NextInterval(t time.Time) (bool, time.Time) {
	previousStart := r.Frequency.Previous(r.Start, t)
	nextStart := r.Frequency.Next(r.Start, t)
	currentStop := r.Frequency.Next(r.Stop, previousStart)

	if !currentStop.After(t) {
		return false, nextStart
	}

	return true, currentStop
}
