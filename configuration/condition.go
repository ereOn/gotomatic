package configuration

import (
	"fmt"
	"reflect"
	"time"

	"github.com/intelux/gotomatic/conditional"
	"github.com/intelux/gotomatic/executor"
	gtime "github.com/intelux/gotomatic/time"
	"github.com/mitchellh/mapstructure"
)

type conditionDecl struct {
	Name string
	Type string
}

type manualConditionParams struct {
	State bool
}

type inverseConditionParams struct {
	Condition conditional.Condition
}

type delayConditionParams struct {
	Condition conditional.Condition
	Delay     time.Duration
}

type compositeConditionParams struct {
	Conditions []conditional.Condition
}

type timeConditionParams struct {
	Start     time.Time
	Stop      time.Time
	Frequency gtime.Frequency
}

type cutOffConditionParams struct {
	Up       uint
	Down     uint
	Period   time.Duration
	Executor executor.Executor
}

func (c *configurationImpl) stringToCondition() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf((*conditional.Condition)(nil)).Elem() {
			return data, nil
		}

		name := data.(string)
		condition := c.GetCondition(name)

		if condition == nil {
			return nil, fmt.Errorf("no condition found with the name \"%s\"", name)
		}

		return condition, nil
	}
}

func (c *configurationImpl) mapToCondition() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t != reflect.TypeOf((*conditional.Condition)(nil)).Elem() {
			return data, nil
		}

		var declaration conditionDecl

		if err := c.decode(data, &declaration); err != nil {
			return data, err
		}

		var condition conditional.Condition

		switch declaration.Type {
		case "manual":
			var params manualConditionParams

			if err := c.decode(data, &params); err != nil {
				return data, err
			}

			condition = conditional.NewManualCondition(params.State)
		case "inverse":
			var params inverseConditionParams

			if err := c.decode(data, &params); err != nil {
				return data, err
			}

			condition = conditional.Inverse(params.Condition)
		case "delay":
			var params delayConditionParams

			if err := c.decode(data, &params); err != nil {
				return data, err
			}

			condition = conditional.Delay(params.Condition, params.Delay)
		case "and":
			var params compositeConditionParams

			if err := c.decode(data, &params); err != nil {
				return data, err
			}

			condition = conditional.NewCompositeCondition(conditional.OperatorAnd, params.Conditions...)
		case "or":
			var params compositeConditionParams

			if err := c.decode(data, &params); err != nil {
				return data, err
			}

			condition = conditional.NewCompositeCondition(conditional.OperatorOr, params.Conditions...)
		case "xor":
			var params compositeConditionParams

			if err := c.decode(data, &params); err != nil {
				return data, err
			}

			condition = conditional.NewCompositeCondition(conditional.OperatorXor, params.Conditions...)
		case "time":
			params := timeConditionParams{
				Frequency: gtime.FrequencyYear,
			}

			if err := c.decode(data, &params); err != nil {
				return data, err
			}

			condition = conditional.NewTimeCondition(gtime.NewRecurrentMoment(params.Start, params.Stop, params.Frequency))
		case "cut-off":
			params := cutOffConditionParams{
				Up:       0,
				Down:     3,
				Period:   time.Second * 5,
				Executor: executor.FalseExecutor,
			}

			if err := c.decode(data, &params); err != nil {
				return data, err
			}

			condition = conditional.NewCutOffCondition(params.Up, params.Down, params.Period, params.Executor)
		default:
			return data, fmt.Errorf("unknown condition type: %s", declaration.Type)
		}

		if declaration.Name != "" {
			if err := c.AddCondition(declaration.Name, condition); err != nil {
				return data, err
			}

			condition = conditional.Dereference(condition)
		}

		return condition, nil
	}
}
