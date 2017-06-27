package trigger

import (
	"context"
	"fmt"
	"os/exec"
)

type commandAction struct {
	cmd  string
	args []string
	env  []string
}

// NewCommanAction instantiates a new action that executes a command.
//
// Two environment variables are added before the command gets executed:
//
// - `GOTOMATIC_CONDITION_NAME`: The name of the condition whose state changed,
// if it has one.
//
// - `GOTOMATIC_CONDITION_STATE`: The state of the condition, as 0 or 1.
func NewCommandAction(cmd string, args []string, env []string) Action {
	return &commandAction{
		cmd:  cmd,
		args: args,
		env:  env,
	}
}

func (t *commandAction) run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, t.cmd, t.args...)
	cmd.Env = t.env

	if name := GetConditionName(ctx); name != nil {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOTOMATIC_CONDITION_NAME=%s", *name))
	}

	if state := GetConditionState(ctx); state != nil {
		var stateInt int

		if *state {
			stateInt = 1
		}

		cmd.Env = append(cmd.Env, fmt.Sprintf("GOTOMATIC_CONDITION_STATE=%d", stateInt))
	}

	return cmd.Run()
}
