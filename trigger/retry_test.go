package trigger

import (
	"context"
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	c := 0
	f := func(ctx context.Context) error {
		if c < 2 {
			c++
			return errors.New("fail")
		}

		return nil
	}

	action := Retry(FuncAction(f), 3, 0)
	ctx := context.Background()
	err := action.run(ctx)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}

func TestRetryFailure(t *testing.T) {
	f := func(ctx context.Context) error {
		return errors.New("fail")
	}

	action := Retry(FuncAction(f), 3, 0)
	ctx := context.Background()
	err := action.run(ctx)

	if err == nil {
		t.Error("expected an error but didn't get one")
	}
}
