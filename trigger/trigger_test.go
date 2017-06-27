package trigger

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/intelux/gotomatic/conditional"
)

func TestWatch(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	condition := conditional.NewManualCondition(false)
	wasUp := false

	up := FuncAction(func(context.Context) error {
		fmt.Println("up")
		wasUp = true
		condition.Set(false)
		return nil
	})

	down := FuncAction(func(context.Context) error {
		fmt.Println("down")
		if wasUp {
			return errors.New("fail")
		}

		condition.Set(true)
		return nil
	})

	trigger := Trigger{
		Up:   up,
		Down: down,
	}

	go condition.Set(true)
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
