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
}

// TimeCondition represents a condition that is met as long as the current time
// is in the specified range.
type TimeCondition struct {
	Condition
	TimeRange TimeRange
	timeFunc  func() time.Time
	sleepFunc func(time.Duration, <-chan struct{})
	done      chan struct{}
}

func sleep(duration time.Duration, interrupt <-chan struct{}) {
	timer := time.NewTimer(duration)

	select {
	case <-timer.C:
	case <-interrupt:
		timer.Stop()
		break
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
	SleepFunction func(time.Duration, <-chan struct{})
}

func (o SleepFunctionOption) apply(condition *TimeCondition) {
	condition.sleepFunc = o.SleepFunction
}

func (condition *TimeCondition) Close() error {
	if condition.done != nil {
		close(condition.done)
		condition.done = nil
	}

	return condition.Condition.Close()
}

func (condition *TimeCondition) checkTime() error {
	for {
		now := condition.timeFunc()
		var next time.Time

		if condition.TimeRange.Contains(now) {
			next = condition.TimeRange.NextStop(now)
			condition.Condition.(*ManualCondition).Set(true)
		} else {
			next = condition.TimeRange.NextStart(now)
			condition.Condition.(*ManualCondition).Set(false)
		}

		condition.sleepFunc(next.Sub(now), condition.done)
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

// Hour returns the hour of the DayTime.
func (d DayTime) Minute() int {
	hours := time.Duration(d.Hour()) * time.Hour
	return int((time.Duration(d) - hours) / time.Minute)
}

// Hour returns the hour of the DayTime.
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
	} else {
		return dayTime < r.Stop || dayTime >= r.Start
	}
}

func next(t time.Time, ref DayTime) time.Time {
	year, month, day := t.Date()
	dayTime := NewDayTime(t.Clock())

	if dayTime < ref {
		return time.Date(year, month, day, ref.Hour(), ref.Minute(), ref.Second(), 0, t.Location())

	}

	return time.Date(year, month, day+1, ref.Hour(), ref.Minute(), ref.Second(), 0, t.Location())
}

// Next returns the next time that will start the range.
func (r DayTimeRange) NextStart(t time.Time) time.Time {
	return next(t, r.Start)
}

// Next returns the next time that will stop the range.
func (r DayTimeRange) NextStop(t time.Time) time.Time {
	return next(t, r.Stop)
}

// String returns the string representation of the DayTimeRange.
func (d DayTimeRange) String() string {
	return fmt.Sprintf("between %s and %s", d.Start, d.Stop)
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

// Next returns the next time that will start the range.
func (r WeekdaysRange) NextStart(t time.Time) time.Time {
	year, month, day := t.Date()
	currentWeekday := t.Weekday()

	// The default is to return next week.
	var result time.Time = time.Date(year, month, day+7, 0, 0, 0, 0, t.Location())

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

// Next returns the next time that will stop the range.
func (r WeekdaysRange) NextStop(t time.Time) time.Time {
	year, month, day := t.Date()
	currentWeekday := t.Weekday()

	// The default is to return next week.
	var result time.Time = time.Date(year, month, day+7, 0, 0, 0, 0, t.Location())

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
func (d WeekdaysRange) String() string {
	s := make([]string, len(d.Weekdays), len(d.Weekdays))

	for i, wd := range d.Weekdays {
		s[i] = wd.String()
	}

	return strings.Join(s, ", ")
}
