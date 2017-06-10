package trigger

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"
)

func TestRetry(t *testing.T) {
	c := 0
	f := func(w io.Writer, name string, state bool) error {
		if c < 2 {
			c++
			return errors.New("fail")
		}

		return nil
	}

	trigger := Retry(TriggerFunc(f), 3, 0)
	err := trigger.run(ioutil.Discard, "", true)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}

func TestRetryFailure(t *testing.T) {
	f := func(w io.Writer, name string, state bool) error {
		return errors.New("fail")
	}

	trigger := Retry(TriggerFunc(f), 3, 0)
	err := trigger.run(ioutil.Discard, "", true)

	if err == nil {
		t.Error("expected an error but didn't get one")
	}
}
