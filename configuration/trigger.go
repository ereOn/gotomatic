package configuration

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/intelux/gotomatic/conditional"
	"github.com/intelux/gotomatic/trigger"
	"github.com/mitchellh/mapstructure"
)

func (c *configurationImpl) mapToConditionTrigger() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t != reflect.TypeOf((*ConditionTrigger)(nil)).Elem() {
			return data, nil
		}

		var trigger trigger.Trigger

		err := c.decode(data, &trigger)

		if err != nil {
			return data, err
		}

		decl := struct {
			Condition conditional.Condition
		}{}

		err = c.decode(data, &decl)

		if err != nil {
			return data, err
		}

		if decl.Condition == nil {
			return data, errors.New("a condition is mandatory")
		}

		return ConditionTrigger{
			Trigger:   trigger,
			Condition: decl.Condition,
		}, nil
	}
}

func (c *configurationImpl) mapToTrigger() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t != reflect.TypeOf((*trigger.Trigger)(nil)).Elem() {
			return data, nil
		}

		declaration := struct {
			Type string
		}{}

		err := c.decode(data, &declaration)

		if err != nil {
			return data, err
		}

		var trig trigger.Trigger

		switch declaration.Type {
		case "command":
			params := struct {
				Command string
				Args    []string
				Env     map[string]string
			}{}

			err := c.decode(data, &params)

			if err != nil {
				return data, err
			}

			if params.Command == "" {
				return data, errors.New("a command is mandatory for that trigger type")
			}

			env := make([]string, 0)

			for key, value := range params.Env {
				env = append(env, fmt.Sprintf("%s=%s", key, value))
			}

			trig = trigger.NewCommandTrigger(params.Command, params.Args, env)

		default:
			return data, fmt.Errorf("unknown trigger type: %s", declaration.Type)
		}

		return trig, nil
	}
}
