package trigger

import (
	"context"
	"os"
	"testing"
)

func TestCommandAction(t *testing.T) {
	action := NewCommandAction("echo", nil, os.Environ())
	ctx := context.Background()
	ctx = WithConditionName(ctx, "foo")
	ctx = WithConditionState(ctx, true)
	err := action.run(ctx)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}
