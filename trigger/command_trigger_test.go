package trigger

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestCommandTrigger(t *testing.T) {
	trigger := NewCommandTrigger("sh", []string{"-c", "echo $GOTOMATIC_CONDITION_NAME"}, os.Environ())
	buf := &bytes.Buffer{}
	err := trigger.run(buf, "foo", true)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	expected := "foo\n"
	output, _ := ioutil.ReadAll(buf)

	if string(output) != expected {
		t.Errorf("expected `%s`, got `%s`", expected, string(output))
	}
}
