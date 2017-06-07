package configuration

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/intelux/gotomatic/conditional"

	yaml "gopkg.in/yaml.v2"
)

func readYAMLFixture(path string) interface{} {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(f)

	if err != nil {
		panic(err)
	}

	var result interface{}

	err = yaml.Unmarshal(data, &result)

	if err != nil {
		panic(err)
	}

	return result
}

func TestMapToCondition(t *testing.T) {
	testCases := []struct {
		Fixture       string
		ExpectFailure bool
	}{
		{"fixture/invalid.yaml", true},
		{"fixture/invalid-name.yaml", true},
		{"fixture/invalid-type.yaml", true},
		{"fixture/invalid-manual-condition.yaml", true},
		{"fixture/invalid-inverse-condition.yaml", true},
		{"fixture/invalid-delay-condition.yaml", true},
		{"fixture/invalid-delay-condition-subcondition.yaml", true},
		{"fixture/invalid-composite-condition-and.yaml", true},
		{"fixture/invalid-composite-condition-or.yaml", true},
		{"fixture/invalid-composite-condition-xor.yaml", true},
		{"fixture/invalid-composite-condition-and-subcondition.yaml", true},
		{"fixture/invalid-composite-condition-or-subcondition.yaml", true},
		{"fixture/invalid-composite-condition-xor-subcondition.yaml", true},
		{"fixture/invalid-time-condition.yaml", true},
		{"fixture/invalid-cut-off-condition.yaml", true},
		{"fixture/invalid-cut-off-condition-invalid-executor.yaml", true},
		{"fixture/invalid-cut-off-condition-invalid-cmd-executor.yaml", true},
		{"fixture/invalid-cut-off-condition-invalid-http-executor.yaml", true},
		{"fixture/invalid-cut-off-condition-unknown-executor.yaml", true},
		{"fixture/unknown-type.yaml", true},
		{"fixture/manual-condition.yaml", false},
		{"fixture/inverse-condition.yaml", false},
		{"fixture/delay-condition.yaml", false},
		{"fixture/composite-condition-and.yaml", false},
		{"fixture/composite-condition-or.yaml", false},
		{"fixture/composite-condition-xor.yaml", false},
		{"fixture/time-condition.yaml", false},
		{"fixture/cut-off-condition-cmd.yaml", false},
		{"fixture/cut-off-condition-http.yaml", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Fixture, func(t *testing.T) {
			configuration := newConfigurationImpl()
			defer configuration.Clear()

			var condition conditional.Condition

			data := readYAMLFixture(testCase.Fixture)
			err := configuration.decode(data, &condition)

			if testCase.ExpectFailure {
				if err == nil {
					t.Error("expected a failure but didn't get one")
				}

				if condition != nil {
					t.Errorf("expected no condition but got: %v", condition)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %s", err)
				}

				if condition == nil {
					t.Error("expected a condition but didn't get one")
				}
			}
		})
	}
}

func TestMapToConditionList(t *testing.T) {
	testCases := []struct {
		Fixture       string
		ExpectFailure bool
	}{
		{"fixture/invalid.yaml", true},
		{"fixture/invalid-list.yaml", true},
		{"fixture/duplicate-list.yaml", true},
		{"fixture/complete.yaml", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Fixture, func(t *testing.T) {
			configuration := newConfigurationImpl()
			defer configuration.Clear()

			var conditions []conditional.Condition

			data := readYAMLFixture(testCase.Fixture)
			err := configuration.decode(data, &conditions)

			if testCase.ExpectFailure {
				if err == nil {
					t.Error("expected a failure but didn't get one")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %s", err)
				}

				if conditions == nil {
					t.Error("expected a list of conditions but didn't get one")
				}
			}
		})
	}
}
