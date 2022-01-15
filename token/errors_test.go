package token

import (
	goErr "errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/errors"
	"github.com/stretchr/testify/assert"
)

func Test_error(t *testing.T) {

	t.Run("getInvalidExpressionEmptyError", func(t *testing.T) {
		actual := getInvalidExpressionEmptyError()
		assert.EqualError(t, actual, "invalid expression. is empty")
		assert.True(t, goErr.Is(actual, errors.ErrInvalidExpression))
	})

	t.Run("getInvalidTokenEmpty", func(t *testing.T) {
		actual := getInvalidTokenEmpty()
		assert.EqualError(t, actual, "invalid token. token string is empty")
		assert.True(t, goErr.Is(actual, errors.ErrInvalidToken))
	})
}

func Test_isInvalidTokenTargetError(t *testing.T) {
	tests := []struct {
		input    error
		expected bool
	}{
		{
			input:    fmt.Errorf("invalid token"),
			expected: false,
		},
		{
			input:    fmt.Errorf("is %w", errors.ErrInvalidTokenTarget),
			expected: true,
		},
		{
			input:    fmt.Errorf("is not %v", errors.ErrInvalidTokenTarget),
			expected: false,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := isInvalidTokenTargetError(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_getInvalidExpressionError(t *testing.T) {

	tests := []struct {
		input    error
		expected string
	}{
		{
			input:    errors.ErrInvalidExpression,
			expected: "invalid expression",
		},
		{
			input:    fmt.Errorf("invalid expression. this"),
			expected: "invalid expression. invalid expression. this",
		},
		{
			input:    fmt.Errorf("%w target", errors.ErrInvalidExpression),
			expected: "invalid expression target",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidExpressionError(test.input)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidExpression))
		})
	}
}

func Test_getInvalidExpressionFormatError(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "invalid expression. invalid format ''",
		},
		{
			input:    "()",
			expected: "invalid expression. invalid format '()'",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidExpressionFormatError(test.input)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidExpression))
		})
	}
}

func Test_getInvalidTokenArgumentError(t *testing.T) {

	type input struct {
		tokenType string
		got       reflect.Kind
		expected  []reflect.Kind
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input:    input{},
			expected: ": invalid token argument. expected [] got [invalid]",
		},
		{
			input: input{
				tokenType: "test",
				got:       reflect.Bool,
				expected:  []reflect.Kind{reflect.Array, reflect.Slice},
			},
			expected: "test: invalid token argument. expected [array slice] got [bool]",
		},
		{
			input: input{
				tokenType: "test",
				got:       reflect.Bool,
				expected:  []reflect.Kind{reflect.Array},
			},
			expected: "test: invalid token argument. expected [array] got [bool]",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidTokenArgumentError(test.input.tokenType, test.input.got, test.input.expected...)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidToken))
		})
	}
}

func Test_getInvalidTokenArgumentNilError(t *testing.T) {

	type input struct {
		tokenType string
		expected  []reflect.Kind
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input:    input{},
			expected: ": invalid token argument. expected [] got [nil]",
		},
		{
			input: input{
				tokenType: "test",
				expected:  []reflect.Kind{reflect.Array},
			},
			expected: "test: invalid token argument. expected [array] got [nil]",
		},
		{
			input: input{
				tokenType: "test",
				expected:  []reflect.Kind{reflect.Array, reflect.Slice},
			},
			expected: "test: invalid token argument. expected [array slice] got [nil]",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidTokenArgumentNilError(test.input.tokenType, test.input.expected...)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidToken))
		})
	}
}

func Test_getInvalidTokenError(t *testing.T) {

	type input struct {
		tokenType string
		reason    error
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input: input{
				tokenType: "test",
				reason:    fmt.Errorf("the reason"),
			},
			expected: "test: invalid token the reason",
		},
		{
			input: input{
				tokenType: "other",
				reason:    fmt.Errorf("different error"),
			},
			expected: "other: invalid token different error",
		},
		{
			input: input{
				tokenType: "other",
				reason:    getInvalidTokenError("test", fmt.Errorf("the reason")),
			},
			expected: "test: invalid token the reason",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidTokenError(test.input.tokenType, test.input.reason)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidToken))
		})
	}
}

func Test_getInvalidTokenFormatError(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "invalid token. '' does not match any token format",
		},
		{
			input:    "()",
			expected: "invalid token. '()' does not match any token format",
		},
		{
			input:    "--",
			expected: "invalid token. '--' does not match any token format",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidTokenFormatError(test.input)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidToken))
		})
	}
}

func Test_getInvalidTokenKeyNotFoundError(t *testing.T) {

	type input struct {
		tokenType, key string
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input:    input{},
			expected: ": invalid token key '' not found",
		},
		{
			input: input{
				tokenType: "test",
			},
			expected: "test: invalid token key '' not found",
		},
		{
			input: input{
				key: "key",
			},
			expected: ": invalid token key 'key' not found",
		},
		{
			input: input{
				tokenType: "test",
				key:       "key",
			},
			expected: "test: invalid token key 'key' not found",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidTokenKeyNotFoundError(test.input.tokenType, test.input.key)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidToken))
		})
	}
}

func Test_getInvalidTokenOutOfRangeError(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: ": invalid token out of range",
		},
		{
			input:    "test",
			expected: "test: invalid token out of range",
		},
		{
			input:    "other",
			expected: "other: invalid token out of range",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidTokenOutOfRangeError(test.input)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidToken))
		})
	}
}

func Test_getInvalidTokenTargetError(t *testing.T) {

	type input struct {
		tokenType string
		got       reflect.Kind
		expected  []reflect.Kind
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input:    input{},
			expected: ": invalid token target. expected [] got [invalid]",
		},
		{
			input: input{
				tokenType: "test",
			},
			expected: "test: invalid token target. expected [] got [invalid]",
		},
		{
			input: input{
				got: reflect.Bool,
			},
			expected: ": invalid token target. expected [] got [bool]",
		},
		{
			input: input{
				expected: []reflect.Kind{reflect.Array},
			},
			expected: ": invalid token target. expected [array] got [invalid]",
		},
		{
			input: input{
				tokenType: "test",
				got:       reflect.Bool,
				expected:  []reflect.Kind{reflect.Array},
			},
			expected: "test: invalid token target. expected [array] got [bool]",
		},
		{
			input: input{
				tokenType: "test",
				got:       reflect.Bool,
				expected:  []reflect.Kind{reflect.Array, reflect.Slice},
			},
			expected: "test: invalid token target. expected [array slice] got [bool]",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidTokenTargetError(test.input.tokenType, test.input.got, test.input.expected...)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidTokenTarget))
		})
	}
}

func Test_getInvalidTokenTargetNilError(t *testing.T) {

	type input struct {
		tokenType string
		expected  []reflect.Kind
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input:    input{},
			expected: ": invalid token target. expected [] got [nil]",
		},
		{
			input: input{
				tokenType: "test",
			},
			expected: "test: invalid token target. expected [] got [nil]",
		},
		{
			input: input{
				expected: []reflect.Kind{reflect.Array},
			},
			expected: ": invalid token target. expected [array] got [nil]",
		},
		{
			input: input{
				tokenType: "test",
				expected:  []reflect.Kind{reflect.Array},
			},
			expected: "test: invalid token target. expected [array] got [nil]",
		},
		{
			input: input{
				tokenType: "test",
				expected:  []reflect.Kind{reflect.Array, reflect.Slice},
			},
			expected: "test: invalid token target. expected [array slice] got [nil]",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getInvalidTokenTargetNilError(test.input.tokenType, test.input.expected...)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrInvalidTokenTarget))
		})
	}
}

func Test_getUnexpectedExpressionResultError(t *testing.T) {

	type input struct {
		got      reflect.Kind
		expected []reflect.Kind
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input:    input{},
			expected: "unexpected expression result. expected [] got [invalid]",
		},
		{
			input: input{
				got: reflect.Bool,
			},
			expected: "unexpected expression result. expected [] got [bool]",
		},
		{
			input: input{
				expected: []reflect.Kind{reflect.Array},
			},
			expected: "unexpected expression result. expected [array] got [invalid]",
		},
		{
			input: input{
				got:      reflect.Bool,
				expected: []reflect.Kind{reflect.Array},
			},
			expected: "unexpected expression result. expected [array] got [bool]",
		},
		{
			input: input{
				got:      reflect.Bool,
				expected: []reflect.Kind{reflect.Array, reflect.Slice},
			},
			expected: "unexpected expression result. expected [array slice] got [bool]",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getUnexpectedExpressionResultError(test.input.got, test.input.expected...)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrUnexpectedExpressionResult))
		})
	}
}

func Test_getUnexpectedExpressionResultNilError(t *testing.T) {

	tests := []struct {
		input    []reflect.Kind
		expected string
	}{
		{
			input:    nil,
			expected: "unexpected expression result. expected [] got [nil]",
		},
		{
			input:    []reflect.Kind{},
			expected: "unexpected expression result. expected [] got [nil]",
		},
		{
			input:    []reflect.Kind{reflect.Array},
			expected: "unexpected expression result. expected [array] got [nil]",
		},
		{
			input:    []reflect.Kind{reflect.Array, reflect.Slice},
			expected: "unexpected expression result. expected [array slice] got [nil]",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getUnexpectedExpressionResultNilError(test.input...)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrUnexpectedExpressionResult))
		})
	}
}

func Test_getUnexpectedTokenError(t *testing.T) {

	type input struct {
		tokenType string
		index     int
	}

	tests := []struct {
		input    input
		expected string
	}{
		{
			input:    input{},
			expected: "unexpected token '' at index 0",
		},
		{
			input: input{
				tokenType: "test",
			},
			expected: "unexpected token 'test' at index 0",
		},
		{
			input: input{
				index: 1,
			},
			expected: "unexpected token '' at index 1",
		},
		{
			input: input{
				tokenType: "test",
				index:     2,
			},
			expected: "unexpected token 'test' at index 2",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := getUnexpectedTokenError(test.input.tokenType, test.input.index)
			assert.EqualError(t, actual, test.expected)
			assert.True(t, goErr.Is(actual, errors.ErrUnexpectedToken))
		})
	}
}
