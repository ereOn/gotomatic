package conditional

import "time"

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
// is in all the associated ranges.
type TimeCondition struct {
	Condition
	Ranges []TimeRange
	done   chan struct{}
}

// NewTimeCondition instantiates a new TimeCondition.
func NewTimeCondition(ranges ...TimeRange) *TimeCondition {
	condition := &TimeCondition{
		Condition: NewManualCondition(false),
		Ranges:    ranges,
		done:      make(chan struct{}),
	}

	go condition.checkTime()

	return condition
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
		now := time.Now()
		later := now
		satisfied := true

		for _, range_ := range condition.Ranges {
			var next time.Time

			if range_.Contains(now) {
				next = range_.NextStop(now)
			} else {
				satisfied = false
				next = range_.NextStart(now)
			}

			if next.Before(later) {
				later = next
			}
		}

		// If later == now, we are in the ranges so we set to true.
		condition.Condition.(*ManualCondition).Set(satisfied)
		timer := time.NewTimer(later.Sub(now))

		select {
		case <-timer.C:
		case <-condition.done:
			timer.Stop()
			break
		}
	}
}

// TimeOfDay represents a time in a day.
type TimeOfDay time.Duration

// ToTimeOfDay converts a time.Time to a TimeOfDay.
func ToTimeOfDay(t time.Time) TimeOfDay {
	hour := time.Duration(t.Hour()) * time.Hour
	minute := time.Duration(t.Minute()) * time.Minute
	second := time.Duration(t.Second()) * time.Second

	return TimeOfDay(hour + minute + second)
}

// Before returns true if the time of day is before the specified time.
func (t TimeOfDay) Before(td TimeOfDay) bool {
	return time.Duration(t).Seconds() < time.Duration(td).Seconds()
}

// After returns true if the time of day is after or eauql to the specified time.
func (t TimeOfDay) After(td TimeOfDay) bool {
	return time.Duration(t).Seconds() >= time.Duration(td).Seconds()
}

// DayRange represents a range of hours in a day.
type DayRange struct {
	Start TimeOfDay
	Stop  TimeOfDay
}

// Contains returns true if the specified time is contained in the range,
// false otherwise.
func (r DayRange) Contains(t time.Time) bool {
	td := ToTimeOfDay(t)

	if r.Start.Before(r.Stop) {
		return td.After(r.Start) && td.Before(r.Stop)
	} else {
		return td.Before(r.Start) || td.After(r.Stop)
	}
}

// Next returns the next time that will start the range.
func (r DayRange) NextStart(t time.Time) time.Time {
	td := ToTimeOfDay(t)

	if td.Before(r.Start) {
		return t.Add(time.Duration(r.Start) - time.Duration(td))
	} else {
		return t.Add(time.Duration(td) - time.Duration(r.Start) + time.Hour*24)
	}
}

// Next returns the next time that will stop the range.
func (r DayRange) NextStop(t time.Time) time.Time {
	td := ToTimeOfDay(t)

	if td.Before(r.Stop) {
		return t.Add(time.Duration(r.Stop) - time.Duration(td))
	} else {
		return t.Add(time.Duration(td) - time.Duration(r.Stop) + time.Hour*24)
	}
}
