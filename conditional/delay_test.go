package conditional

import (
	"testing"
	"time"
)

func TestRealTimer(t *testing.T) {
	timer := realTimer{
		timer: time.NewTimer(time.Second),
	}
	timer.Stop()
}

func TestForeverTimer(t *testing.T) {
	timer := foreverTimer{
		channel: make(chan time.Time),
	}
	timer.Stop()
}

func TestDelay(t *testing.T) {
	m := NewManualCondition(false)
	condition := Delay(m, 10*time.Millisecond)
	defer condition.Close()

	assertConditionState(t, condition, false, "initialization to false")
	m.Set(true)
	assertConditionState(t, condition, true, "waiting for a while")

	// Let's do a fast change and back.
	m.Set(false)
	assertConditionState(t, condition, false, "waiting again for a while")

	m.Set(true)
	assertConditionState(t, condition, true, "waiting again for a while")
}
