[![Build Status](https://travis-ci.org/intelux/gotomatic.svg?branch=master)](https://travis-ci.org/intelux/gotomatic)
[![GoDoc](https://godoc.org/github.com/intelux/gotomatic?status.svg)](https://godoc.org/github.com/intelux/gotomatic)

# gotomatic

Gotomatic is both a library and a set of tools to deal with conditions.

All the condition types provided by gotomatic implement
[Condition](https://godoc.org/github.com/intelux/gotomatic/conditional#Condition)
which is a very simple, yet powerful interface.

All condition types are thread-safe and have a low-memory footprint. They rely
on channels internally and are very effective.

Conditions can be composed logically (and, or, xor), negated, delayed or even
set from external sources (commands, HTTP requests).

Here is a simple example:

```
// Create a new manual condition, which can be set programmatically.
condition := conditional.NewManualCondition(false)
defer condition.Close()

// Set the condition in another goroutine.
go condition.Set(true)

// Wait for the condition to become true.
condition.Wait(true)
```

Here is a more complex one:

```
// Create a new time condition, which is only set between 10:00 and 11:00.
conditionA := conditional.NewTimeCondition(
	conditional.NewRecurrentMoment(
		time.Date(0, 1, 1, 10, 0, 0, 0, time.Local),
		time.Date(0, 1, 1, 11, 0, 0, 0, time.Local),
		conditional.FrequencyDay,
	)
)

// Create a new time condition, which is only set the first day of each month.
conditionB := conditional.NewTimeCondition(
	conditional.NewRecurrentMoment(
		time.Date(0, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(0, 1, 2, 0, 0, 0, 0, time.Local),
		conditional.FrequencyMonth,
	)
)

// Create a new condition, which is only set the first day of each month,
between 10:00 and 11:00.
condition := conditional.NewCompositeCondition(conditional.OperatorAnd, conditionA, conditionB)
defer condition.Close()

// Wait for the condition to become true.
condition.Wait(true)
```
