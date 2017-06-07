package configuration

import (
	"fmt"
	"reflect"
	"time"

	"github.com/intelux/gotomatic/executor"
	"github.com/mitchellh/mapstructure"
)

type executorDecl struct {
	Type    string
	Timeout time.Duration
}

type commandExecutorParams struct {
	Command string
	Args    []string
}

type httpExecutorParams struct {
	Method      string
	URL         string
	StatusCodes []int
}

func mapToExecutor() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t != reflect.TypeOf((*executor.Executor)(nil)).Elem() {
			return data, nil
		}

		declaration := executorDecl{
			Timeout: time.Second,
		}

		err := mapstructure.Decode(data, &declaration)

		if err != nil {
			return data, err
		}

		switch declaration.Type {
		case "cmd":
			var params commandExecutorParams

			err := mapstructure.Decode(data, &params)

			if err != nil {
				return data, err
			}

			return executor.CommandExecutor(params.Command, params.Args...), nil
		case "http":
			params := httpExecutorParams{
				Method: "GET",
				StatusCodes: []int{
					200,
					201,
				},
			}

			err := mapstructure.Decode(data, &params)

			if err != nil {
				return data, err
			}

			return executor.HTTPExecutor(params.Method, params.URL, params.StatusCodes, declaration.Timeout), nil
		}

		return data, fmt.Errorf("unknown command type \"%s\"", declaration.Type)
	}
}
