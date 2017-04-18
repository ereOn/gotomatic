package conditional

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
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

type conditionDocument struct {
	Type   string      `json:"type" yaml:"type"`
	Name   string      `json:"name,omitempty" yaml:"name,omitempty"`
	Params interface{} `json:"params,omitempty" yaml:"params,omitempty"`
}

func (e commonEncoder) decode(document conditionDocument, condition *Condition) error {
	switch document.Type {
	case "manual":
		*condition = &ManualCondition{}
	default:
		return fmt.Errorf("unknown type `%s` for condition", document.Type)
	}

	return (*condition).decodeDocumentParams(document.Params)
}

// String returns a string that describes the encoder.
func (JSONEncoder) String() string { return "JSONEncoder" }

// Marshal serializes a Condition.
func (JSONEncoder) Marshal(condition Condition) ([]byte, error) {
	return json.Marshal(conditionDocument{
		Type:   condition.documentType(),
		Params: condition.documentParams(),
	})
}

// Unmarshal deserializes a Condition.
func (e JSONEncoder) Unmarshal(data []byte, condition *Condition) error {
	document := conditionDocument{}

	if err := json.Unmarshal(data, &document); err != nil {
		return err
	}

	return e.decode(document, condition)
}

// String returns a string that describes the encoder.
func (YAMLEncoder) String() string { return "YAMLEncoder" }

// Marshal serializes a Condition.
func (YAMLEncoder) Marshal(condition Condition) ([]byte, error) {
	return yaml.Marshal(conditionDocument{
		Type:   condition.documentType(),
		Params: condition,
	})
}

// Unmarshal deserializes a Condition.
func (e YAMLEncoder) Unmarshal(data []byte, condition *Condition) error {
	document := conditionDocument{}

	if err := yaml.Unmarshal(data, &document); err != nil {
		return err
	}

	return e.decode(document, condition)
}

// Manual conditions.

func (c *ManualCondition) documentType() string {
	return "manual"
}

type manualConditionParams struct{}

func (c *ManualCondition) documentParams() interface{} {
	return manualConditionParams{}
}

func (c *ManualCondition) decodeDocumentParams(params interface{}) error {
	// This is not efficient, but simplifies the code *A LOT*.
	data, _ := yaml.Marshal(params)
	p := manualConditionParams{}

	return yaml.Unmarshal(data, &p)
}
