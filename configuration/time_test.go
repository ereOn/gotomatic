package configuration

import (
	"testing"
	"time"

	gtime "github.com/intelux/gotomatic/time"
)

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
		Expected gtime.Frequency
	}{
		{"year", gtime.FrequencyYear},
		{"month", gtime.FrequencyMonth},
		{"week", gtime.FrequencyWeek},
		{"day", gtime.FrequencyDay},
		{"hour", gtime.FrequencyHour},
		{"minute", gtime.FrequencyMinute},
		{"second", gtime.FrequencySecond},
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
