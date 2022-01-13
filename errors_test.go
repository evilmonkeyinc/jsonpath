package jsonpath

import (
	"fmt"
	"testing"

	goErr "errors"

	"github.com/evilmonkeyinc/jsonpath/errors"
	"github.com/stretchr/testify/assert"
)

func Test_getInvalidJSONData(t *testing.T) {

	tests := []struct {
		input    error
		expected string
	}{
		{
			input:    fmt.Errorf("a reason"),
			expected: "invalid data. a reason",
		},
		{
			input:    fmt.Errorf("different reason"),
			expected: "invalid data. different reason",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidJSONData(test.input)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidJSONData))
		})
	}
}

func Test_getInvalidJSONPathQuery(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "invalid JSONPath query ''",
		},
		{
			input:    "test",
			expected: "invalid JSONPath query 'test'",
		},
		{
			input:    "other",
			expected: "invalid JSONPath query 'other'",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidJSONPathQuery(test.input)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidJSONPathQuery))
		})
	}
}
