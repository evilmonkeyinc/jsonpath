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

func Test_getInvalidJSONPathQueryWithReason(t *testing.T) {

	type input struct {
		query  string
		reason error
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input: input{
				query:  "",
				reason: fmt.Errorf("the reason"),
			},
			expected: "invalid JSONPath query '' the reason",
		},
		{
			input: input{
				query:  "test",
				reason: fmt.Errorf("other reason"),
			},
			expected: "invalid JSONPath query 'test' other reason",
		},
		{
			input: input{
				query:  "other",
				reason: getInvalidJSONPathQuery("inside"),
			},
			expected: "invalid JSONPath query 'inside'",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidJSONPathQueryWithReason(test.input.query, test.input.reason)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidJSONPathQuery))
		})
	}
}
