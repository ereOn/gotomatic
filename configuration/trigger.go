package configuration

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/intelux/gotomatic/trigger"
	"github.com/mitchellh/mapstructure"
)

func (c *configurationImpl) mapToAction() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t != reflect.TypeOf((*trigger.Action)(nil)).Elem() {
			return data, nil
		}

		declaration := struct {
			Type string
		}{}

		err := c.decode(data, &declaration)

		if err != nil {
			return data, err
		}

		var action trigger.Action

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
				return data, errors.New("a command is mandatory for that action type")
			}

			env := os.Environ()

			for key, value := range params.Env {
				env = append(env, fmt.Sprintf("%s=%s", key, value))
			}

			action = trigger.NewCommandAction(params.Command, params.Args, env)

		default:
			return data, fmt.Errorf("unknown action type: %s", declaration.Type)
		}

		return action, nil
	}
}
