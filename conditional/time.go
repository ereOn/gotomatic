package conditional

import cron "gopkg.in/robfig/cron.v2"

// TimeCondition represents a condition that is met as long as the current time
// is in a given range.
//
// For recurrent time-based tasks, see CronCondition instead.
type TimeCondition struct {
	Condition
	Cron *cron.Cron
}

// NewTimeCondition instantiates a new TimeCondition.
func NewTimeCondition(on cron.Schedule, off cron.Schedule) *TimeCondition {
	condition := &TimeCondition{
		Condition: NewManualCondition(false),
		Cron:      cron.New(),
	}

	onJob := timeConditionJob{
		Condition: condition.Condition.(*ManualCondition),
		Satisfied: true,
	}
	offJob := timeConditionJob{
		Condition: condition.Condition.(*ManualCondition),
		Satisfied: false,
	}

	condition.Cron.Schedule(on, onJob)
	condition.Cron.Schedule(off, offJob)
	condition.Cron.Start()

	return condition
}

func (condition *TimeCondition) Close() error {
	condition.Cron.Stop()

	return condition.Condition.Close()
}

type timeConditionJob struct {
	Condition *ManualCondition
	Satisfied bool
}

func (j timeConditionJob) Run() {
	j.Condition.Set(j.Satisfied)
}
