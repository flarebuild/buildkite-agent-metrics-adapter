package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToSnakeCase(t *testing.T) {
	// Based on https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"camelCase", "camel_case"},
		{"PascalCase", "pascal_case"},
		{"snake_case", "snake_case"},
		{"Pascal_Snake", "pascal_snake"},
		{"SCREAMING_SNAKE", "screaming_snake"},
		{"kebab-case", "kebab_case"},
		{"Pascal-Kebab", "pascal_kebab"},
		{"SCREAMING-KEBAB", "screaming_kebab"},
		{"A", "a"},
		{"AA", "aa"},
		{"AAA", "aaa"},
		{"AAAA", "aaaa"},
		{"AaAa", "aa_aa"},
		{"HTTPRequest", "http_request"},
		{"BatteryLifeValue", "battery_life_value"},
		{"Id0Value", "id0_value"},
		{"ID0Value", "id0_value"},
	}

	for _, test := range tests {
		result := ToSnakeCase(test.input)
		assert.Equal(t, test.expected, result)
	}
}