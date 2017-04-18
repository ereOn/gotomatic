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

			err = encoder.Unmarshal([]byte{0}, &result)

			if err == nil {
				t.Errorf("expected a deserialization error")
			}
		}
	}
}

func TestGenericConditionSerializationErrorUnknownType(t *testing.T) {
	encoder := commonEncoder{}
	var result Condition
	document := conditionDocument{
		Type: "unknown",
	}
	err := encoder.decode(document, &result)

	if err == nil {
		t.Errorf("expected error")
	}
}

func TestManualConditionSerialization(t *testing.T) {
	reference := NewManualCondition(false)
	defer reference.Close()

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
			t.Errorf("deserialization should have succeeded for %s: %s. Data: %s\n", encoder, err, data)
		}

		if _, ok := result.(*ManualCondition); !ok {
			t.Errorf("condition was deserialized to the wrong type")
		}
	}
}
