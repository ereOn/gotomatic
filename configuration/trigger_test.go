package configuration

import (
	"testing"

	"github.com/intelux/gotomatic/trigger"
)

func TestMapToAction(t *testing.T) {
	testCases := []struct {
		Fixture       string
		ExpectFailure bool
	}{
		{"fixture/invalid.yaml", true},
		{"fixture/action-invalid-type.yaml", true},
		{"fixture/action-unknown-type.yaml", true},
		{"fixture/action-command-invalid-params.yaml", true},
		{"fixture/action-command-empty-command.yaml", true},
		{"fixture/action-command.yaml", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Fixture, func(t *testing.T) {
			configuration := newConfigurationImpl()
			defer configuration.Close()

			var action trigger.Action

			data := readYAMLFixture(testCase.Fixture)
			err := configuration.decode(data, &action)

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
