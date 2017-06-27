package trigger

import (
	"context"
)

type funcAction struct {
	f func(ctx context.Context) error
}

// FuncAction creates a trigger from a function.
func FuncAction(f func(ctx context.Context) error) Action {
	return funcAction{f: f}
}

func (t funcAction) run(ctx context.Context) error {
	return t.f(ctx)
}
