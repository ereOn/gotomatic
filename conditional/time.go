package conditional

import (
	"fmt"
	"strings"
	"time"
)

// TimeRange represents a time range.
type TimeRange interface {
	// Contains returns true if the specified time is contained in the range,
	// false otherwise.
	Contains(time time.Time) bool

	// Next returns the next time that will start the range.
	NextStart(time time.Time) time.Time

	// Next returns the next time that will stop the range.
	NextStop(time time.Time) time.Time

	// String returns the string representation of the time range.
	String() string
}

// TimeCondition represents a condition that is met as long as the current time
// is in the specified range.
type TimeCondition struct {
	Condition
	TimeRange TimeRange
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
func NewTimeCondition(timeRange TimeRange, options ...TimeConditionOption) *TimeCondition {
	condition := &TimeCondition{
		TimeRange: timeRange,
		timeFunc:  time.Now,
		sleepFunc: sleep,
		done:      make(chan struct{}),
	}

	for _, option := range options {
		option.apply(condition)
	}

	condition.Condition = NewManualCondition(timeRange.Contains(condition.timeFunc()))

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
		var next time.Time

		if condition.TimeRange.Contains(now) {
			next = condition.TimeRange.NextStop(now)
			condition.Condition.(*ManualCondition).Set(true)
		} else {
			next = condition.TimeRange.NextStart(now)
			condition.Condition.(*ManualCondition).Set(false)
		}

		delay = next.Sub(now)
	}
}

// DayTime represent at time of the day.
type DayTime time.Duration

// NewDayTime creates a DayTime instance for its hour, minute and second parts.
func NewDayTime(hour int, minute int, second int) DayTime {
	return DayTime(
		time.Duration(hour)*time.Hour +
			time.Duration(minute)*time.Minute +
			time.Duration(second)*time.Second,
	)
}

// Hour returns the hour of the DayTime.
func (d DayTime) Hour() int {
	return int(time.Duration(d) / time.Hour)
}

// Minute returns the minute of the DayTime.
func (d DayTime) Minute() int {
	hours := time.Duration(d.Hour()) * time.Hour
	return int((time.Duration(d) - hours) / time.Minute)
}

// Second returns the second of the DayTime.
func (d DayTime) Second() int {
	hours := time.Duration(d.Hour()) * time.Hour
	minutes := time.Duration(d.Minute()) * time.Minute
	return int((time.Duration(d) - hours - minutes) / time.Second)
}

// String returns the string representation of the DayTime.
func (d DayTime) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", d.Hour(), d.Minute(), d.Second())
}

// DayTimeRange represents a range of hours within a day.
//
// If Start happens after Stop, the range will represents all hours *NOT*
// between Start and Stop. For instance, if Start is 16:00:00 and Stop is
// 12:00:00, the range represents all hours of the day before 12:00:00 and
// after 16:00:00.
type DayTimeRange struct {
	Start DayTime
	Stop  DayTime
}

// Contains returns true if the specified time is contained in the range,
// false otherwise.
func (r DayTimeRange) Contains(t time.Time) bool {
	dayTime := NewDayTime(t.Clock())

	if r.Start < r.Stop {
		return r.Start <= dayTime && dayTime < r.Stop
	}

	return dayTime < r.Stop || dayTime >= r.Start
}

func next(t time.Time, ref DayTime) time.Time {
	year, month, day := t.Date()
	dayTime := NewDayTime(t.Clock())

	if dayTime < ref {
		return time.Date(year, month, day, ref.Hour(), ref.Minute(), ref.Second(), 0, t.Location())

	}

	return time.Date(year, month, day+1, ref.Hour(), ref.Minute(), ref.Second(), 0, t.Location())
}

// NextStart returns the next time that will start the range.
func (r DayTimeRange) NextStart(t time.Time) time.Time {
	return next(t, r.Start)
}

// NextStop returns the next time that will stop the range.
func (r DayTimeRange) NextStop(t time.Time) time.Time {
	return next(t, r.Stop)
}

// String returns the string representation of the DayTimeRange.
func (r DayTimeRange) String() string {
	return fmt.Sprintf("between %s and %s", r.Start, r.Stop)
}

// WeekdaysRange represents a range of days within a week.
type WeekdaysRange struct {
	Weekdays Weekdays
}

// Weekdays represents a list of week days.
type Weekdays []time.Weekday

// Contains returns true if the specified weekday is in the weekdays list.
func (w Weekdays) Contains(weekday time.Weekday) bool {
	for _, wd := range w {
		if wd == weekday {
			return true
		}
	}

	return false
}

// Contains returns true if the specified time is contained in the range,
// false otherwise.
func (r WeekdaysRange) Contains(t time.Time) bool {
	return r.Weekdays.Contains(t.Weekday())
}

// NextStart returns the next time that will start the range.
func (r WeekdaysRange) NextStart(t time.Time) time.Time {
	year, month, day := t.Date()
	currentWeekday := t.Weekday()

	// The default is to return next week.
	result := time.Date(year, month, day+7, 0, 0, 0, 0, t.Location())

	if len(r.Weekdays) < 7 {
		if r.Contains(t) {
			return r.NextStart(r.NextStop(t))
		}

		for i := 1; i <= 7; i = i + 1 {
			weekday := currentWeekday + time.Weekday(i%7)

			if r.Weekdays.Contains(weekday) {
				result = time.Date(year, month, day+i, 0, 0, 0, 0, t.Location())
				break
			}
		}
	}

	return result
}

// NextStop returns the next time that will stop the range.
func (r WeekdaysRange) NextStop(t time.Time) time.Time {
	year, month, day := t.Date()
	currentWeekday := t.Weekday()

	// The default is to return next week.
	result := time.Date(year, month, day+7, 0, 0, 0, 0, t.Location())

	if len(r.Weekdays) < 7 {
		if !r.Contains(t) {
			return r.NextStop(r.NextStart(t))
		}

		for i := 1; i <= 7; i = i + 1 {
			weekday := currentWeekday + time.Weekday(i%7)

			if !r.Weekdays.Contains(weekday) {
				result = time.Date(year, month, day+int(weekday-currentWeekday), 0, 0, 0, 0, t.Location())
				break
			}
		}
	}

	return result
}

// String returns the string representation of the WeekdaysRange.
func (r WeekdaysRange) String() string {
	s := make([]string, len(r.Weekdays), len(r.Weekdays))

	for i, wd := range r.Weekdays {
		s[i] = wd.String()
	}

	return strings.Join(s, ", ")
}

type compositeRange struct {
	ranges   []TimeRange
	operator CompositeOperator
}

// NewCompositeTimeRange creates a new TimeRange which is the composition of
// several other TimeRanges.
func NewCompositeTimeRange(operator CompositeOperator, ranges ...TimeRange) TimeRange {
	if len(ranges) == 0 {
		panic("cannot instantiate a composite range without at least one sub-range")
	}

	return &compositeRange{
		ranges:   ranges,
		operator: operator,
	}
}

// Contains returns true if the specified time is contained in the range,
// false otherwise.
func (r compositeRange) Contains(t time.Time) bool {
	values := make([]bool, len(r.ranges), len(r.ranges))

	for i, r := range r.ranges {
		values[i] = r.Contains(t)
	}

	return r.operator.Reduce(values)
}

// Next returns the next time that will start the range.
func (r compositeRange) NextStart(t time.Time) time.Time {
	var result time.Time

	for _, r := range r.ranges {
		// Because we don't know what the operator is going to be, we stop at
		// the first possible change. This way we may wake up too soon, but
		// never too late.
		x := r.NextStart(t)
		y := r.NextStop(t)

		if x.Before(result) || (result == time.Time{}) {
			result = x
		}

		if y.Before(result) {
			result = y
		}
	}

	return result
}

// Next returns the next time that will stop the range.
func (r compositeRange) NextStop(t time.Time) time.Time {
	return r.NextStart(t)
}

// String returns the string representation of the composite range.
func (r compositeRange) String() string {
	s := make([]string, len(r.ranges), len(r.ranges))

	for i, r := range r.ranges {
		s[i] = fmt.Sprintf("(%s)", r)
	}

	return strings.Join(s, fmt.Sprintf(" %s ", r.operator))
}

// Frequency represents a frequency.
type Frequency interface {
	getBase(start time.Time, t time.Time) time.Time
	Previous(start time.Time, t time.Time) time.Time
	Next(start time.Time, t time.Time) time.Time
}

var (
	FrequencyYear   = frequencyYear{}
	FrequencyMonth  = frequencyMonth{}
	FrequencyWeek   = frequencyWeek{}
	FrequencyDay    = frequencyDay{}
	FrequencyHour   = frequencyHour{}
	FrequencyMinute = frequencyMinute{}
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

	if !r.Before(t) {
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

	if !r.Before(t) {
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

	if !r.Before(t) {
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

	if !r.Before(t) {
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

	if !r.Before(t) {
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

	if !r.Before(t) {
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

	if !r.Before(t) {
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
type Moment struct {
}
