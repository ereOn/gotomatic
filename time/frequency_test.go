package time

import (
	"testing"
	"time"
)

func TestFrequencyYear(t *testing.T) {
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	start := time.Date(2016, 2, 29, 11, 15, 30, 50, edt)
	now := time.Date(2020, 4, 10, 15, 45, 10, 20, edt)
	expected := time.Date(2020, 2, 29, 11, 15, 30, 50, edt)
	value := FrequencyYear.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2021, 3, 1, 11, 15, 30, 50, edt)
	value = FrequencyYear.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	now = time.Date(2020, 1, 39, 15, 45, 10, 20, edt)
	expected = time.Date(2019, 2, 29, 11, 15, 30, 50, edt)
	value = FrequencyYear.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2020, 2, 29, 11, 15, 30, 50, edt)
	value = FrequencyYear.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestFrequencyMonth(t *testing.T) {
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	start := time.Date(2016, 2, 29, 11, 15, 30, 50, edt)
	now := time.Date(2020, 4, 10, 15, 45, 10, 20, edt)
	expected := time.Date(2020, 3, 29, 11, 15, 30, 50, edt)
	value := FrequencyMonth.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2020, 4, 29, 11, 15, 30, 50, edt)
	value = FrequencyMonth.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	now = time.Date(2020, 4, 30, 15, 45, 10, 20, edt)
	expected = time.Date(2020, 4, 29, 11, 15, 30, 50, edt)
	value = FrequencyMonth.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2020, 5, 29, 11, 15, 30, 50, edt)
	value = FrequencyMonth.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestFrequencyWeek(t *testing.T) {
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	start := time.Date(2017, 4, 25, 11, 15, 30, 50, edt) // Tuesday.
	now := time.Date(2017, 4, 10, 15, 45, 10, 20, edt)
	expected := time.Date(2017, 4, 4, 11, 15, 30, 50, edt)
	value := FrequencyWeek.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 11, 11, 15, 30, 50, edt)
	value = FrequencyWeek.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	now = time.Date(2017, 4, 11, 15, 45, 10, 20, edt)
	expected = time.Date(2017, 4, 11, 11, 15, 30, 50, edt)
	value = FrequencyWeek.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 18, 11, 15, 30, 50, edt)
	value = FrequencyWeek.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestFrequencyDay(t *testing.T) {
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	start := time.Date(2017, 4, 25, 11, 15, 30, 50, edt) // Tuesday.
	now := time.Date(2017, 4, 10, 15, 45, 10, 20, edt)
	expected := time.Date(2017, 4, 10, 11, 15, 30, 50, edt)
	value := FrequencyDay.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 11, 11, 15, 30, 50, edt)
	value = FrequencyDay.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	now = time.Date(2017, 4, 10, 10, 45, 10, 20, edt)
	expected = time.Date(2017, 4, 9, 11, 15, 30, 50, edt)
	value = FrequencyDay.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 10, 11, 15, 30, 50, edt)
	value = FrequencyDay.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestFrequencyHour(t *testing.T) {
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	start := time.Date(2017, 4, 25, 11, 15, 30, 50, edt) // Tuesday.
	now := time.Date(2017, 4, 10, 15, 45, 10, 20, edt)
	expected := time.Date(2017, 4, 10, 15, 15, 30, 50, edt)
	value := FrequencyHour.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 10, 16, 15, 30, 50, edt)
	value = FrequencyHour.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	now = time.Date(2017, 4, 10, 10, 05, 10, 20, edt)
	expected = time.Date(2017, 4, 10, 9, 15, 30, 50, edt)
	value = FrequencyHour.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 10, 10, 15, 30, 50, edt)
	value = FrequencyHour.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestFrequencyMinute(t *testing.T) {
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	start := time.Date(2017, 4, 25, 11, 15, 30, 50, edt) // Tuesday.
	now := time.Date(2017, 4, 10, 15, 45, 10, 20, edt)
	expected := time.Date(2017, 4, 10, 15, 44, 30, 50, edt)
	value := FrequencyMinute.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 10, 15, 45, 30, 50, edt)
	value = FrequencyMinute.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	now = time.Date(2017, 4, 10, 10, 5, 35, 20, edt)
	expected = time.Date(2017, 4, 10, 10, 5, 30, 50, edt)
	value = FrequencyMinute.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 10, 10, 6, 30, 50, edt)
	value = FrequencyMinute.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}

func TestFrequencySecond(t *testing.T) {
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	start := time.Date(2017, 4, 25, 11, 15, 30, 50, edt) // Tuesday.
	now := time.Date(2017, 4, 10, 15, 45, 10, 20, edt)
	expected := time.Date(2017, 4, 10, 15, 45, 9, 50, edt)
	value := FrequencySecond.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 10, 15, 45, 10, 50, edt)
	value = FrequencySecond.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	now = time.Date(2017, 4, 10, 15, 45, 10, 51, edt)
	expected = time.Date(2017, 4, 10, 15, 45, 10, 50, edt)
	value = FrequencySecond.Previous(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}

	expected = time.Date(2017, 4, 10, 15, 45, 11, 50, edt)
	value = FrequencySecond.Next(start, now)

	if value != expected {
		t.Errorf("expected: %s, got: %s", expected, value)
	}
}
