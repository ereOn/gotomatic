package conditional

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeCondition(t *testing.T) {
	dayRange := DayRange{
		Start: ToTimeOfDay(time.Date(2000, 1, 1, 1, 15, 0, 0, time.Local)),
		Stop:  ToTimeOfDay(time.Date(2000, 1, 1, 1, 16, 0, 0, time.Local)),
	}
	condition := NewTimeCondition(dayRange)
	defer condition.Close()

	for {
		state, change := condition.GetAndWaitChange()
		fmt.Printf("%s: %v", time.Now(), state)
		<-change
	}
}
