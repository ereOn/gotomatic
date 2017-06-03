package conditional

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

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

func TestRegistryDecodeCondition(t *testing.T) {
	testCases := []struct {
		Fixture       string
		ExpectFailure bool
	}{
		{"fixture/invalid.yaml", true},
		{"fixture/invalid-name.yaml", true},
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
			registry := NewRegistry()
			defer registry.Close()

			data := readYAMLFixture(testCase.Fixture)
			condition, err := registry.DecodeCondition(data)

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

func TestRegistryDecodeConditions(t *testing.T) {
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
			registry := NewRegistry()
			defer registry.Close()

			data := readYAMLFixture(testCase.Fixture)
			conditions, err := registry.DecodeConditions(data)

			if testCase.ExpectFailure {
				if err == nil {
					t.Error("expected a failure but didn't get one")
				}

				if conditions != nil {
					t.Errorf("expected no conditions but got: %v", conditions)
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

func TestParseTime(t *testing.T) {
	testCases := []struct {
		Value    string
		Expected string
	}{
		{"13:02", "0000-01-01T13:02:00+01:00"},
		{"13:02Z", "0000-01-01T13:02:00Z"},
		{"13:02:45", "0000-01-01T13:02:45+01:00"},
		{"13:02:45Z", "0000-01-01T13:02:45Z"},
		{"22/11", "0000-11-22T00:00:00+01:00"},
		{"22/11Z", "0000-11-22T00:00:00Z"},
		{"Mar 13", "0000-03-13T00:00:00+01:00"},
		{"Mar", "0000-03-01T00:00:00+01:00"},
		{"2008-11-03", "2008-11-03T00:00:00+01:00"},
		{"2008-11-03 13:07:34", "2008-11-03T13:07:34+01:00"},
		{"2008-11-03T13:07:34+01:00", "2008-11-03T13:07:34+01:00"},
		{"2008-11-03T13:07:34Z", "2008-11-03T13:07:34Z"},
		{"monday", "0000-01-01T00:00:00+01:00"},
		{"TUESDAY", "0000-01-02T00:00:00+01:00"},
		{"wedNESday", "0000-01-03T00:00:00+01:00"},
		{"thursday", "0000-01-04T00:00:00+01:00"},
		{"friday", "0000-01-05T00:00:00+01:00"},
		{"saturday", "0000-01-06T00:00:00+01:00"},
		{"sunday", "0000-01-07T00:00:00+01:00"},
	}

	zone := time.FixedZone("XXX", 3600)

	for _, testCase := range testCases {
		t.Run(testCase.Value, func(t *testing.T) {
			expected, _ := time.ParseInLocation(time.RFC3339, testCase.Expected, zone)
			value, err := parseTime(testCase.Value, zone)

			if err != nil {
				t.Errorf("expected no error but got: %s", err)
			}

			if value != expected {
				t.Errorf("expected %v, got %v", expected, value)
			}
		})
	}
}

func TestParseTimeFailure(t *testing.T) {
	testCases := []struct {
		Value string
	}{
		{""},
		{"xxxxxxxx"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Value, func(t *testing.T) {
			_, err := parseTime(testCase.Value, time.Local)

			if err == nil {
				t.Error("expected an error but didn't get one")
			}
		})
	}
}

func TestParseFrequency(t *testing.T) {
	testCases := []struct {
		Value    string
		Expected Frequency
	}{
		{"year", FrequencyYear},
		{"month", FrequencyMonth},
		{"week", FrequencyWeek},
		{"day", FrequencyDay},
		{"hour", FrequencyHour},
		{"minute", FrequencyMinute},
		{"second", FrequencySecond},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Value, func(t *testing.T) {
			value, err := parseFrequency(testCase.Value)

			if err != nil {
				t.Errorf("expected no error but got: %s", err)
			}

			if value != testCase.Expected {
				t.Errorf("expected %v, got %v", testCase.Expected, value)
			}
		})
	}
}

func TestParseFrequencyFailure(t *testing.T) {
	testCases := []struct {
		Value string
	}{
		{""},
		{"xxxxxxxx"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Value, func(t *testing.T) {
			_, err := parseFrequency(testCase.Value)

			if err == nil {
				t.Error("expected an error but didn't get one")
			}
		})
	}
}
