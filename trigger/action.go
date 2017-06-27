package trigger

import "context"

// Action represents an action performed by a Trigger.
type Action interface {
	run(ctx context.Context) error
}

type actionKey int

const (
	conditionNameKey actionKey = iota
	conditionStateKey
)

// WithConditionName injects a condition name in the specified context.
func WithConditionName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, conditionNameKey, name)
}

// GetConditionName gets the condition name from a context.
func GetConditionName(ctx context.Context) *string {
	if name, ok := ctx.Value(conditionNameKey).(string); ok {
		return &name
	}

	return nil
}

// WithConditionState injects a condition state in the specified context.
func WithConditionState(ctx context.Context, state bool) context.Context {
	return context.WithValue(ctx, conditionStateKey, state)
}

// GetConditionState gets the condition state from a context.
func GetConditionState(ctx context.Context) *bool {
	if state, ok := ctx.Value(conditionStateKey).(bool); ok {
		return &state
	}

	return nil
}
