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

func TestWeekdaysContains(t *testing.T) {
	weekdays := Weekdays{time.Monday, time.Saturday}
	values := []struct {
		Weekday  time.Weekday
		Expected bool
	}{
		{
			Weekday:  time.Monday,
			Expected: true,
		},
		{
			Weekday:  time.Tuesday,
			Expected: false,
		},
		{
			Weekday:  time.Wednesday,
			Expected: false,
		},
		{
			Weekday:  time.Thursday,
			Expected: false,
		},
		{
			Weekday:  time.Friday,
			Expected: false,
		},
		{
			Weekday:  time.Saturday,
			Expected: true,
		},
		{
			Weekday:  time.Sunday,
			Expected: false,
		},
	}

	for _, value := range values {
		if value.Expected != weekdays.Contains(value.Weekday) {
			if value.Expected {
				t.Errorf("expected %s to be in %s", value.Weekday, weekdays)
			} else {
				t.Errorf("expected %s to not be in %s", value.Weekday, weekdays)
			}
		}
	}
}

func TestWeekdaysRangeContains(t *testing.T) {
	weekdaysRange := WeekdaysRange{
		Weekdays: Weekdays{time.Monday, time.Saturday},
	}

	makeWeekday := func(wd time.Weekday) time.Time {
		return time.Date(2017, 3, 12+int(wd), 12, 0, 0, 0, time.Local)
	}

	values := []struct {
		Time     time.Time
		Expected bool
	}{
		{
			Time:     makeWeekday(time.Monday),
			Expected: true,
		},
		{
			Time:     makeWeekday(time.Tuesday),
			Expected: false,
		},
		{
			Time:     makeWeekday(time.Wednesday),
			Expected: false,
		},
		{
			Time:     makeWeekday(time.Thursday),
			Expected: false,
		},
		{
			Time:     makeWeekday(time.Friday),
			Expected: false,
		},
		{
			Time:     makeWeekday(time.Saturday),
			Expected: true,
		},
		{
			Time:     makeWeekday(time.Sunday),
			Expected: false,
		},
	}

	for _, value := range values {
		if value.Expected != weekdaysRange.Contains(value.Time) {
			if value.Expected {
				t.Errorf("expected %s to be in %s", value.Time, weekdaysRange)
			} else {
				t.Errorf("expected %s to not be in %s", value.Time, weekdaysRange)
			}
		}
	}
}

func TestWeekdaysRangeNextStart(t *testing.T) {
	weekdaysRange := WeekdaysRange{
		Weekdays: Weekdays{time.Monday, time.Tuesday, time.Saturday},
	}

	makeWeekday := func(wd time.Weekday) time.Time {
		return time.Date(2017, 3, 12+int(wd), 0, 0, 0, 0, time.Local)
	}

	weekday := makeWeekday(time.Sunday)
	value := weekdaysRange.NextStart(weekday)
	expected := makeWeekday(time.Monday)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	weekday = makeWeekday(time.Monday)
	value = weekdaysRange.NextStart(weekday)
	expected = makeWeekday(time.Saturday)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestWeekdaysRangeNextStop(t *testing.T) {
	weekdaysRange := WeekdaysRange{
		Weekdays: Weekdays{time.Monday, time.Tuesday, time.Saturday},
	}

	makeWeekday := func(wd time.Weekday) time.Time {
		return time.Date(2017, 3, 12+int(wd), 0, 0, 0, 0, time.Local)
	}

	weekday := makeWeekday(time.Sunday)
	value := weekdaysRange.NextStop(weekday)
	expected := makeWeekday(time.Wednesday)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	weekday = makeWeekday(time.Monday)
	value = weekdaysRange.NextStop(weekday)
	expected = makeWeekday(time.Wednesday)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestWeekdaysRangeString(t *testing.T) {
	weekdaysRange := WeekdaysRange{
		Weekdays: Weekdays{time.Monday, time.Tuesday, time.Saturday},
	}
	expected := "Monday, Tuesday, Saturday"
	value := weekdaysRange.String()

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}
