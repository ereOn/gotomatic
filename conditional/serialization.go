package conditional

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type conditionDeclaration struct {
	Name string
	Type string
}

type manualConditionParams struct {
	State bool
}

type inverseConditionParams struct {
	Condition interface{}
}

type delayConditionParams struct {
	Condition interface{}
	Delay     time.Duration
}

// Registry represents a condition registry.
type Registry interface {
	DecodeCondition(input interface{}) (Condition, error)
	DecodeConditions(input interface{}) ([]Condition, error)
	Close()
}

type registryImpl struct {
	index map[string]Condition
}

// NewRegistry instantiates a new condition registry.
func NewRegistry() Registry {
	return &registryImpl{
		index: make(map[string]Condition),
	}
}

func (r *registryImpl) decode(m interface{}, rawVal interface{}) error {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
	})

	return decoder.Decode(m)
}

// DecodeCondition decodes a condition from an arbitrary input structure, if
// possible.
func (r *registryImpl) DecodeCondition(input interface{}) (condition Condition, err error) {
	var name string

	if err = r.decode(input, &name); err == nil {
		var ok bool

		if condition, ok = r.index[name]; !ok {
			return nil, fmt.Errorf("no condition found with the name `%s`", name)
		}

		return Dereference(condition), err
	}

	var declaration conditionDeclaration

	if err = r.decode(input, &declaration); err != nil {
		return nil, err
	}

	switch declaration.Type {
	case "manual":
		var params manualConditionParams

		if err = r.decode(input, &params); err != nil {
			return nil, err
		}

		condition = NewManualCondition(params.State)
	case "inverse":
		var params inverseConditionParams

		// Parsing those specific params can never fail.
		r.decode(input, &params)

		var subcondition Condition

		if subcondition, err = r.DecodeCondition(params.Condition); err != nil {
			return nil, err
		}

		condition = Inverse(subcondition)
	case "delay":
		var params delayConditionParams

		if err = r.decode(input, &params); err != nil {
			return nil, err
		}

		var subcondition Condition

		if subcondition, err = r.DecodeCondition(params.Condition); err != nil {
			return nil, err
		}

		condition = Delay(subcondition, params.Delay)
	default:
		return nil, fmt.Errorf("unknown condition type: %s", declaration.Type)
	}

	if declaration.Name != "" {
		if _, ok := r.index[declaration.Name]; ok {
			condition.Close()

			return nil, fmt.Errorf("a condition with the name `%s` was already registered", declaration.Name)
		}

		r.index[declaration.Name] = condition

		condition = Dereference(condition)
	}

	return condition, nil
}

type conditionDeclarations []interface{}

// DecodeConditions decodes a list of conditions from an arbitrary input structure, if
// possible.
func (r *registryImpl) DecodeConditions(input interface{}) (conditions []Condition, err error) {
	var declarations conditionDeclarations

	if err = r.decode(input, &declarations); err != nil {
		return nil, err
	}

	for _, declaration := range declarations {
		condition, err := r.DecodeCondition(declaration)

		if err != nil {
			for _, condition = range conditions {
				condition.Close()
			}

			return nil, err
		}

		conditions = append(conditions, condition)
	}

	return conditions, nil
}

func (r *registryImpl) Close() {
	for _, condition := range r.index {
		condition.Close()
	}

	r.index = nil
}
