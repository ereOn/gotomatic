package trigger

import (
	"fmt"
	"io"
	"os/exec"
)

type commandTrigger struct {
	cmd  string
	args []string
	env  []string
}

// NewCommandTrigger instantiates a new trigger that executes a command.
//
// Two environment variables are added before the command gets executed:
//
// - `GOTOMATIC_CONDITION_NAME`: The name of the condition whose state changed,
// if it has one.
//
// - `GOTOMATIC_CONDITION_STATE`: The state of the condition, as 0 or 1.
func NewCommandTrigger(cmd string, args []string, env []string) Trigger {
	return &commandTrigger{
		cmd:  cmd,
		args: args,
		env:  env,
	}
}

func (t *commandTrigger) run(w io.Writer, name string, state bool) error {
	cmd := exec.Command(t.cmd, t.args...)
	cmd.Env = t.env

	if name != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOTOMATIC_CONDITION_NAME=%s", name))
	}

	var stateInt int

	if state {
		stateInt = 1
	}

	cmd.Env = append(cmd.Env, fmt.Sprintf("GOTOMATIC_CONDITION_STATE=%d", stateInt))

	cmd.Stderr = w
	cmd.Stdout = w

	return cmd.Run()
}
