package configuration

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	gtime "github.com/intelux/gotomatic/time"
	"github.com/mitchellh/mapstructure"
)

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

func stringToTimeHookFunc(loc *time.Location) mapstructure.DecodeHookFunc {
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

func parseFrequency(s string) (gtime.Frequency, error) {
	switch s {
	case "year":
		return gtime.FrequencyYear, nil
	case "month":
		return gtime.FrequencyMonth, nil
	case "week":
		return gtime.FrequencyWeek, nil
	case "day":
		return gtime.FrequencyDay, nil
	case "hour":
		return gtime.FrequencyHour, nil
	case "minute":
		return gtime.FrequencyMinute, nil
	case "second":
		return gtime.FrequencySecond, nil
	}

	return nil, fmt.Errorf("unknown frequency \"%s\"", s)
}

// stringToFrequencyFunc transforms a string into a frequency.
func stringToFrequencyFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf((*gtime.Frequency)(nil)).Elem() {
			return data, nil
		}

		return parseFrequency(data.(string))
	}
}
