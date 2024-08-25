package utils

import (
	"strings"

	"github.com/mitchellh/mapstructure"
)

func SplitString(s string) []string {
	tags := []string{}
	if s != "" {
		tags = strings.Split(s, ",")
	}

	return tags
}

func Decode(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
