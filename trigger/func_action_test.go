package trigger

import (
	"context"
	"testing"
)

func TestActionFunc(t *testing.T) {
	f := func(ctx context.Context) error {
		return nil
	}

	action := FuncAction(f)
	ctx := context.Background()
	err := action.run(ctx)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}
