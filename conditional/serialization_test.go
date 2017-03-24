package conditional

import (
	"testing"
)

func TestJSONEncoderString(t *testing.T) {
	value := JSONEncoder{}.String()
	expected := "JSONEncoder"

	if value != expected {
		t.Errorf("expected `%s`, got `%s`", expected, value)
	}
}

func TestYAMLEncoderString(t *testing.T) {
	value := YAMLEncoder{}.String()
	expected := "YAMLEncoder"

	if value != expected {
		t.Errorf("expected `%s`, got `%s`", expected, value)
	}
}

func TestGenericConditionSerialization(t *testing.T) {
	for _, reference := range []Condition{
		NewManualCondition(false),
	} {
		defer reference.Close()
		reference.SetName("foo")

		for _, encoder := range []ConditionEncoder{
			JSONEncoder{},
			YAMLEncoder{},
		} {
			data, err := encoder.Marshal(reference)

			if err != nil {
				t.Errorf("serialization should have succeeded for %s: %s", encoder, err)
			}

			var result Condition
			err = encoder.Unmarshal(data, &result)

			if err != nil {
				t.Errorf("deserialization should have succeeded for %s: %s", encoder, err)
			}

			if result == nil {
				t.Errorf("deserialized condition should not be nil for %s", encoder)
			}

			if result.Name() != "foo" {
				t.Errorf("deserialized condition should be named `foo` but is named `%s`", result.Name())
			}

			err = encoder.Unmarshal([]byte{0}, &result)

			if err == nil {
				t.Errorf("expected a deserialization error")
			}
		}
	}
}

func TestGenericConditionSerializationErrorNoType(t *testing.T) {
	encoder := commonEncoder{}
	var result Condition
	values := map[string]interface{}{}
	expectedError := ErrTypeMissing
	err := encoder.extract(values, &result)

	if err != expectedError {
		t.Errorf("expected error to be %s, but got %s", expectedError, err)
	}
}

func TestGenericConditionSerializationErrorInvalidType(t *testing.T) {
	encoder := commonEncoder{}
	var result Condition
	values := map[string]interface{}{
		"type": true,
	}
	expectedError := ErrTypeInvalid
	err := encoder.extract(values, &result)

	if err != expectedError {
		t.Errorf("expected error to be %s, but got %s", expectedError, err)
	}
}

func TestGenericConditionSerializationErrorUnknownType(t *testing.T) {
	encoder := commonEncoder{}
	var result Condition
	values := map[string]interface{}{
		"type": "unknown",
	}
	err := encoder.extract(values, &result)

	if err == nil {
		t.Errorf("expected error")
	}
}

func TestManualConditionSerialization(t *testing.T) {
	reference := NewManualCondition(false)
	defer reference.Close()
	reference.SetName("foo")

	for _, encoder := range []ConditionEncoder{
		JSONEncoder{},
		YAMLEncoder{},
	} {
		data, err := encoder.Marshal(reference)

		if err != nil {
			t.Errorf("serialization should have succeeded for %s: %s", encoder, err)
		}

		var result Condition
		err = encoder.Unmarshal(data, &result)

		if err != nil {
			t.Errorf("deserialization should have succeeded for %s: %s", encoder, err)
		}

		if _, ok := result.(*ManualCondition); !ok {
			t.Errorf("condition was deserialized to the wrong type")
		}
	}
}
