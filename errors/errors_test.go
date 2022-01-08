package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsKeyNotFoundError(t *testing.T) {

	tests := []struct {
		input    error
		expected bool
	}{
		{
			input:    nil,
			expected: false,
		},
		{
			input:    errKeyNotFound,
			expected: true,
		},
		{
			input:    GetKeyNotFoundError("key"),
			expected: true,
		},
		{
			input:    fmt.Errorf("'key' key not found in object"),
			expected: false,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := IsKeyNotFoundError(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}

}

func Test_IsFailedToParseExpressionError(t *testing.T) {

	tests := []struct {
		input    error
		expected bool
	}{
		{
			input:    nil,
			expected: false,
		},
		{
			input:    errFailedToParseExpression,
			expected: true,
		},
		{
			input:    GetFailedToParseExpressionError(fmt.Errorf("reason")),
			expected: true,
		},
		{
			input:    fmt.Errorf("failed to parse expression. reason"),
			expected: false,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := IsFailedToParseExpressionError(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}

}

func Test_IsJSONMarshalFailedError(t *testing.T) {

	tests := []struct {
		input    error
		expected bool
	}{
		{
			input:    nil,
			expected: false,
		},
		{
			input:    errJSONMarshalFailed,
			expected: true,
		},
		{
			input:    GetJSONMarshalFailedError(fmt.Errorf("reason")),
			expected: true,
		},
		{
			input:    fmt.Errorf("json marshal failed. reason"),
			expected: false,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := IsJSONMarshalFailedError(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}

}

func Test_IsInvalidExpressionError(t *testing.T) {

	tests := []struct {
		input    error
		expected bool
	}{
		{
			input:    nil,
			expected: false,
		},
		{
			input:    errInvalidExpression,
			expected: true,
		},
		{
			input:    GetInvalidExpressionError(fmt.Errorf("reason")),
			expected: true,
		},
		{
			input:    fmt.Errorf("invalid expression. reason"),
			expected: false,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := IsInvalidExpressionError(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}

}
