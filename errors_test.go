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

func Test_getInvalidJSONPathSelector(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "invalid JSONPath selector ''",
		},
		{
			input:    "test",
			expected: "invalid JSONPath selector 'test'",
		},
		{
			input:    "other",
			expected: "invalid JSONPath selector 'other'",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidJSONPathSelector(test.input)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidJSONPathSelector))
		})
	}
}

func Test_getInvalidJSONPathSelectorWithReason(t *testing.T) {

	type input struct {
		selector string
		reason   error
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input: input{
				selector: "",
				reason:   fmt.Errorf("the reason"),
			},
			expected: "invalid JSONPath selector '' the reason",
		},
		{
			input: input{
				selector: "test",
				reason:   fmt.Errorf("other reason"),
			},
			expected: "invalid JSONPath selector 'test' other reason",
		},
		{
			input: input{
				selector: "other",
				reason:   getInvalidJSONPathSelector("inside"),
			},
			expected: "invalid JSONPath selector 'inside'",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidJSONPathSelectorWithReason(test.input.selector, test.input.reason)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidJSONPathSelector))
		})
	}
}
