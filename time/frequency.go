package time

import "time"

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
