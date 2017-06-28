package trigger

import (
	"context"
	"testing"
)

func TestWithConditionName(t *testing.T) {
	ctx := context.Background()

	if state := GetConditionName(ctx); state != nil {
		t.Errorf("expected no state, but got: %s", *state)
	}

	ctx = WithConditionName(ctx, "foo")

	if state := GetConditionName(ctx); state == nil {
		t.Errorf("expected a state")
	}
}

func TestWithConditionState(t *testing.T) {
	ctx := context.Background()

	if state := GetConditionState(ctx); state != nil {
		t.Errorf("expected no state, but got: %t", *state)
	}

	ctx = WithConditionState(ctx, true)

	if state := GetConditionState(ctx); state == nil {
		t.Errorf("expected a state")
	}
}
