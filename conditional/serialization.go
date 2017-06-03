package conditional

import (
	"fmt"
	"reflect"
	"strings"
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

type compositeConditionParams struct {
	Conditions []interface{}
}

type timeConditionParams struct {
	Start     time.Time
	Stop      time.Time
	Frequency Frequency
}

type cutOffConditionParams struct {
	Up       uint
	Down     uint
	Period   time.Duration
	Executor Executor
}

type externalCommandDeclaration struct {
	Type    string
	Timeout time.Duration
}

type cmdExternalCommandDeclaration struct {
	Command string
	Args    []string
}

type httpExternalCommandDeclaration struct {
	Method      string
	URL         string
	StatusCodes []int
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
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			StringToTimeHookFunc(time.Local),
			StringToFrequencyFunc(),
			MapToExecutor(),
		),
		Result: rawVal,
	})

	return decoder.Decode(m)
}

func (r *registryImpl) decodeCompositeCondition(operator CompositeOperator, input interface{}) (Condition, error) {
	var err error
	var params compositeConditionParams

	if err = r.decode(input, &params); err != nil {
		return nil, err
	}

	subconditions := make([]Condition, len(params.Conditions))

	for i, x := range params.Conditions {
		if subconditions[i], err = r.DecodeCondition(x); err != nil {
			return nil, err
		}
	}

	return NewCompositeCondition(operator, subconditions...), nil
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
	case "and":
		if condition, err = r.decodeCompositeCondition(OperatorAnd, input); err != nil {
			return nil, err
		}
	case "or":
		if condition, err = r.decodeCompositeCondition(OperatorOr, input); err != nil {
			return nil, err
		}
	case "xor":
		if condition, err = r.decodeCompositeCondition(OperatorXor, input); err != nil {
			return nil, err
		}
	case "time":
		params := timeConditionParams{
			Frequency: FrequencyYear,
		}

		if err := r.decode(input, &params); err != nil {
			return nil, err
		}

		condition = NewTimeCondition(NewRecurrentMoment(params.Start, params.Stop, params.Frequency))
	case "cut-off":
		params := cutOffConditionParams{
			Up:       0,
			Down:     3,
			Period:   time.Second * 5,
			Executor: FalseExecutor,
		}

		if err := r.decode(input, &params); err != nil {
			return nil, err
		}

		condition = NewCutOffCondition(params.Up, params.Down, params.Period, params.Executor)
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

func parseTime(s string, loc *time.Location) (t time.Time, err error) {
	formats := []string{
		"15:04",
		"15:04Z07:00",
		"15:04:05",
		"15:04:05Z07:00",
		"02/01",
		"02/01Z07:00",
		"Jan 02",
		"Jan",
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05Z07:00",
		time.RFC3339,
	}

	for _, format := range formats {
		t, err = time.ParseInLocation(format, s, loc)

		if err == nil {
			return
		}
	}

	for i, day := range []time.Weekday{
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
		time.Sunday,
	} {
		if strings.EqualFold(s, day.String()) {
			return time.Date(0, 1, 1+i, 0, 0, 0, 0, loc), nil
		}
	}

	return t, fmt.Errorf("could not parse time \"%s\" as one of \"%s\"", s, strings.Join(formats, "\", \""))
}

// StringToTimeHookFunc transforms a string into a time.Time.
func StringToTimeHookFunc(loc *time.Location) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		return parseTime(data.(string), loc)
	}
}

func parseFrequency(s string) (Frequency, error) {
	switch s {
	case "year":
		return FrequencyYear, nil
	case "month":
		return FrequencyMonth, nil
	case "week":
		return FrequencyWeek, nil
	case "day":
		return FrequencyDay, nil
	case "hour":
		return FrequencyHour, nil
	case "minute":
		return FrequencyMinute, nil
	case "second":
		return FrequencySecond, nil
	}

	return nil, fmt.Errorf("unknown frequency \"%s\"", s)
}

// StringToFrequencyFunc transforms a string into a frequency.
func StringToFrequencyFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf((*Frequency)(nil)).Elem() {
			return data, nil
		}

		return parseFrequency(data.(string))
	}
}

// MapToExecutor transforms a dict into an external command.
func MapToExecutor() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}

		if t != reflect.TypeOf((*Executor)(nil)).Elem() {
			return data, nil
		}

		declaration := externalCommandDeclaration{
			Timeout: time.Second,
		}

		err := mapstructure.Decode(data, &declaration)

		if err != nil {
			return data, err
		}

		switch declaration.Type {
		case "cmd":
			var params cmdExternalCommandDeclaration

			err := mapstructure.Decode(data, &params)

			if err != nil {
				return data, err
			}

			return CommandExecutor(params.Command, params.Args...), nil
		case "http":
			params := httpExternalCommandDeclaration{
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

			return HTTPExecutor(params.Method, params.URL, params.StatusCodes, declaration.Timeout), nil
		}

		return data, fmt.Errorf("unknown command type \"%s\"", declaration.Type)
	}
}
