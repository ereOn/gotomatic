package configuration

import (
	"time"

	"github.com/mitchellh/mapstructure"
)

func (c *configurationImpl) decode(m interface{}, rawVal interface{}) error {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			stringToTimeHookFunc(time.Local),
			stringToFrequencyFunc(),
			mapToExecutor(),
			c.mapToCondition(),
			c.stringToCondition(),
		),
		Result: rawVal,
	})

	return decoder.Decode(m)
}
