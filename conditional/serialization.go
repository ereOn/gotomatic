package conditional

import (
	"encoding/json"
	"errors"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

var (
	// ErrTypeMissing is returned when a type is missing from a serialized condition.
	ErrTypeMissing = errors.New("no type specified")
	// ErrTypeInvalid is returned when a type is invalid.
	ErrTypeInvalid = errors.New("type must be a string")
)

// ConditionEncoder is the interface for all condition encoders.
type ConditionEncoder interface {
	Marshal(Condition) ([]byte, error)
	Unmarshal([]byte, *Condition) error
}

type commonEncoder struct{}

// JSONEncoder is the JSON encoder class for conditions.
type JSONEncoder struct{ commonEncoder }

// YAMLEncoder is the YAML encoder class for conditions.
type YAMLEncoder struct{ commonEncoder }

func (commonEncoder) extract(values map[string]interface{}, condition *Condition) error {
	var conditionType string

	if v, ok := values["type"]; !ok {
		return ErrTypeMissing
	} else if conditionType, ok = v.(string); !ok {
		return ErrTypeInvalid
	}

	switch conditionType {
	case "manual":
		*condition = NewManualCondition(false)
	default:
		return fmt.Errorf("unknown type `%s` for condition", conditionType)
	}

	if v, ok := values["name"]; ok {
		if name, ok := v.(string); ok {
			(*condition).SetName(name)
		}
	}

	return nil
}

// String returns a string that describes the encoder.
func (JSONEncoder) String() string { return "JSONEncoder" }

// Marshal serializes a Condition.
func (JSONEncoder) Marshal(condition Condition) ([]byte, error) {
	return json.Marshal(condition)
}

// Unmarshal deserializes a Condition.
func (e JSONEncoder) Unmarshal(data []byte, condition *Condition) error {
	values := make(map[string]interface{})

	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	return e.extract(values, condition)
}

// String returns a string that describes the encoder.
func (YAMLEncoder) String() string { return "YAMLEncoder" }

// Marshal serializes a Condition.
func (YAMLEncoder) Marshal(condition Condition) ([]byte, error) {
	return yaml.Marshal(condition)
}

// Unmarshal deserializes a Condition.
func (e YAMLEncoder) Unmarshal(data []byte, condition *Condition) error {
	values := make(map[string]interface{})

	if err := yaml.Unmarshal(data, &values); err != nil {
		return err
	}

	return e.extract(values, condition)
}

func (c *ManualCondition) serializedData() interface{} {
	return struct {
		Name string `json:"name" yaml:"name"`
		Type string `json:"type" yaml:"type"`
	}{
		Name: c.name,
		Type: "manual",
	}
}

// MarshalJSON serializes a condition in JSON.
func (c *ManualCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.serializedData())
}

// MarshalYAML serializes a condition in YAML.
func (c *ManualCondition) MarshalYAML() (interface{}, error) {
	return c.serializedData(), nil
}
