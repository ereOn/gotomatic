// +build !race

package conditional

import (
	"testing"
	"time"
)

func TestSleep(t *testing.T) {
	interrupt := make(chan struct{})
	sleep(0, interrupt)
	close(interrupt)
	sleep(time.Second, interrupt)
}

func TestTimeCondition(t *testing.T) {
	now := time.Date(2017, 3, 9, 11, 59, 59, 0, time.Local)
	dayTimeRange := DayTimeRange{
		Start: NewDayTime(12, 0, 0),
		Stop:  NewDayTime(12, 0, 1),
	}
	condition := NewTimeCondition(
		dayTimeRange,
		TimeFunctionOption{
			TimeFunction: func() time.Time { return now },
		},
		SleepFunctionOption{
			SleepFunction: func(time.Duration, <-chan struct{}) {},
		},
	)
	defer condition.Close()

	assertConditionState(t, condition, false, "initialization before date")

	assertConditionChanged(t, condition, false, "a second passed", func() {
		now = now.Add(time.Second)
	})

	assertConditionChanged(t, condition, true, "another second passed", func() {
		now = now.Add(time.Second)
	})
}

func TestDayTimeHour(t *testing.T) {
	dayTime := NewDayTime(1, 2, 3)
	expected := 1
	value := dayTime.Hour()

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestDayTimeMinute(t *testing.T) {
	dayTime := NewDayTime(1, 2, 3)
	expected := 2
	value := dayTime.Minute()

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestDayTimeSecond(t *testing.T) {
	dayTime := NewDayTime(1, 2, 3)
	expected := 3
	value := dayTime.Second()

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestDayTimeString(t *testing.T) {
	dayTime := NewDayTime(1, 2, 3)
	expected := "01:02:03"
	value := dayTime.String()

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestDayTimeRangeContains(t *testing.T) {
	dayTimeRange := DayTimeRange{
		Start: NewDayTime(17, 0, 0),
		Stop:  NewDayTime(18, 0, 0),
	}
	reversedDayTimeRange := DayTimeRange{
		Start: dayTimeRange.Stop,
		Stop:  dayTimeRange.Start,
	}
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	a := time.Date(2017, 3, 12, 17, 45, 0, 0, edt)
	b := time.Date(2017, 3, 12, 18, 15, 0, 0, edt)

	if !dayTimeRange.Contains(a) {
		t.Errorf("%s does not contain %s", dayTimeRange, a)
	}

	if reversedDayTimeRange.Contains(a) {
		t.Errorf("%s contains %s", reversedDayTimeRange, a)
	}

	if dayTimeRange.Contains(b) {
		t.Errorf("%s contains %s", dayTimeRange, b)
	}

	if !reversedDayTimeRange.Contains(b) {
		t.Errorf("%s does not contain %s", reversedDayTimeRange, b)
	}
}

func TestDayTimeRangeNextStart(t *testing.T) {
	dayTimeRange := DayTimeRange{
		Start: NewDayTime(17, 0, 0),
		Stop:  NewDayTime(18, 0, 0),
	}
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	// Note: these dates happen to overlap a DST change ! Do not change them
	// arbitrarily.
	a := time.Date(2017, 3, 12, 16, 45, 0, 0, edt)
	b := time.Date(2017, 3, 12, 17, 15, 0, 0, edt)
	today := time.Date(2017, 3, 12, 17, 0, 0, 0, edt)
	tomorrow := time.Date(2017, 3, 13, 17, 0, 0, 0, edt)
	value := dayTimeRange.NextStart(a)

	if value != today {
		t.Errorf("expected: %s, got: %s", today, value)
	}

	value = dayTimeRange.NextStart(b)

	if value != tomorrow {
		t.Errorf("expected: %s, got: %s", tomorrow, value)
	}
}

func TestDayTimeRangeNextStop(t *testing.T) {
	dayTimeRange := DayTimeRange{
		Start: NewDayTime(16, 0, 0),
		Stop:  NewDayTime(17, 0, 0),
	}
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	// Note: these dates happen to overlap a DST change ! Do not change them
	// arbitrarily.
	a := time.Date(2017, 3, 12, 16, 45, 0, 0, edt)
	b := time.Date(2017, 3, 12, 17, 15, 0, 0, edt)
	today := time.Date(2017, 3, 12, 17, 0, 0, 0, edt)
	tomorrow := time.Date(2017, 3, 13, 17, 0, 0, 0, edt)
	value := dayTimeRange.NextStop(a)

	if value != today {
		t.Errorf("expected: %s, got: %s", today, value)
	}

	value = dayTimeRange.NextStop(b)

	if value != tomorrow {
		t.Errorf("expected: %s, got: %s", tomorrow, value)
	}
}

func TestDayTimeRangeString(t *testing.T) {
	dayTimeRange := DayTimeRange{
		Start: NewDayTime(17, 0, 0),
		Stop:  NewDayTime(18, 0, 0),
	}

	expected := "between 17:00:00 and 18:00:00"
	value := dayTimeRange.String()

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}
