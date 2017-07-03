package trigger

import (
	"context"
	"errors"
	"testing"

	"github.com/intelux/gotomatic/conditional"
)

func TestWatch(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	condition := conditional.NewManualCondition(false)

	up := FuncAction(func(context.Context) error {
		return errors.New("fail")
	})

	down := FuncAction(func(context.Context) error {
		return errors.New("fail")
	})

	trigger := Trigger{
		Up:   up,
		Down: down,
	}

	err := Watch(ctx, condition, trigger)

	if err == nil {
		t.Errorf("expected an error")
	}

	cancel()

	err = Watch(ctx, condition, trigger)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}
