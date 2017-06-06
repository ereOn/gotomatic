// +build !race

package conditional

import (
	"testing"
	"time"

	gtime "github.com/intelux/gotomatic/time"
)

func TestSleep(t *testing.T) {
	interrupt := make(chan struct{})
	sleep(0, interrupt)
	close(interrupt)
	sleep(time.Second, interrupt)
}

func TestTimeCondition(t *testing.T) {
	now := time.Date(2017, 3, 9, 11, 59, 59, 0, time.Local)
	moment := gtime.NewRecurrentMoment(
		time.Date(1900, 1, 1, 12, 0, 0, 0, time.Local),
		time.Date(1900, 1, 1, 12, 0, 1, 0, time.Local),
		gtime.FrequencyDay,
	)
	condition := NewTimeCondition(
		moment,
		TimeFunctionOption{
			TimeFunction: func() time.Time { return now },
		},
		SleepFunctionOption{
			SleepFunction: func(time.Duration, <-chan struct{}) bool { return true },
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
