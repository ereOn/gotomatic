package trigger

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/intelux/gotomatic/conditional"
)

func TestWatch(t *testing.T) {
	condition := conditional.NewManualCondition(false)
	result := make(chan error)

	triggered := make(chan bool)
	defer close(triggered)

	trigger := TriggerFunc(func(w io.Writer, name string, state bool) error {
		triggered <- state
		return nil
	})

	go func() {
		defer close(result)
		result <- Watch(condition, trigger, ioutil.Discard, "name", true)
	}()

	condition.Set(true)

	if state := <-triggered; !state {
		t.Error("expected state to be true")
	}

	condition.Close()

	if err := <-result; err != conditional.ErrConditionClosed {
		t.Errorf("expected %s error but got: %s", conditional.ErrConditionClosed, err)
	}
}

func TestWatchTriggerFailure(t *testing.T) {
	condition := conditional.NewManualCondition(false)
	defer condition.Close()

	result := make(chan error)
	fail := errors.New("fail")

	trigger := TriggerFunc(func(w io.Writer, name string, state bool) error {
		return fail
	})

	go func() {
		defer close(result)
		result <- Watch(condition, trigger, ioutil.Discard, "name", true)
	}()

	condition.Set(true)

	if err := <-result; err != fail {
		t.Errorf("expected %s error but got: %s", fail, err)
	}
}
