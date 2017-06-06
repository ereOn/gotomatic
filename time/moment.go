package time

import "time"

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
