package time

import (
	"testing"
	"time"
)

func TestRecurrentMoment(t *testing.T) {
	edt, err := time.LoadLocation("Canada/Eastern")

	if err != nil {
		panic(err)
	}

	moment := NewRecurrentMoment(
		time.Date(1990, 1, 1, 1, 15, 10, 0, edt),
		time.Date(1990, 1, 1, 1, 40, 10, 0, edt),
		FrequencyHour,
	)

	now := time.Date(2017, 4, 10, 15, 22, 0, 0, edt)
	expectedState := true
	expectedTime := time.Date(2017, 4, 10, 15, 40, 10, 0, edt)
	state, vtime := moment.NextInterval(now)

	if expectedState != state {
		t.Errorf("expected: %t, got: %t", expectedState, state)
	}

	if expectedTime != vtime {
		t.Errorf("expected: %s, got: %s", expectedTime, vtime)
	}

	now = time.Date(2017, 4, 10, 15, 42, 0, 0, edt)
	expectedState = false
	expectedTime = time.Date(2017, 4, 10, 16, 15, 10, 0, edt)
	state, vtime = moment.NextInterval(now)

	if expectedState != state {
		t.Errorf("expected: %t, got: %t", expectedState, state)
	}

	if expectedTime != vtime {
		t.Errorf("expected: %s, got: %s", expectedTime, vtime)
	}
}
