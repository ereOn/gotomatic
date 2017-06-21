package trigger

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestTriggerFunc(t *testing.T) {
	f := func(w io.Writer, name string, state bool) error {
		w.Write([]byte(name))

		if state {
			w.Write([]byte("1"))
		} else {
			w.Write([]byte("0"))
		}

		return nil
	}

	trigger := Func(f)

	buf := &bytes.Buffer{}
	err := trigger.run(buf, "foo", true)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	expected := "foo1"
	output, _ := ioutil.ReadAll(buf)

	if string(output) != expected {
		t.Errorf("expected `%s`, got `%s`", expected, string(output))
	}
}
