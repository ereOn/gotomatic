package configuration

import (
	"testing"
)

func TestMapToConditionTrigger(t *testing.T) {
	testCases := []struct {
		Fixture       string
		ExpectFailure bool
	}{
		{"fixture/invalid.yaml", true},
		{"fixture/condition-trigger-invalid-params.yaml", true},
		{"fixture/condition-trigger-invalid-params-invalid-action.yaml", true},
		{"fixture/condition-trigger-invalid-params-invalid-action-command.yaml", true},
		{"fixture/condition-trigger-invalid-params-no-command.yaml", true},
		{"fixture/condition-trigger-unknown-type.yaml", true},
		{"fixture/condition-trigger-invalid-condition.yaml", true},
		{"fixture/condition-trigger-missing-condition.yaml", true},
		{"fixture/condition-trigger.yaml", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Fixture, func(t *testing.T) {
			configuration := newConfigurationImpl()
			defer configuration.Close()

			var conditionTrigger ConditionTrigger

			data := readYAMLFixture(testCase.Fixture)
			err := configuration.decode(data, &conditionTrigger)

			if testCase.ExpectFailure {
				if err == nil {
					t.Error("expected a failure but didn't get one")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %s", err)
				}
			}
		})
	}
}
