package standard

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getInteger(t *testing.T) {

	type input struct {
		argument   interface{}
		parameters map[string]interface{}
	}

	type expected struct {
		value int64
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{},
			expected: expected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: input{
				argument: 3,
			},
			expected: expected{
				value: 3,
			},
		},
		{
			input: input{
				argument: 3.14,
			},
			expected: expected{
				err: "invalid argument. expected integer",
			},
		},
		{
			input: input{
				argument: "3.14",
			},
			expected: expected{
				err: "invalid argument. expected integer",
			},
		},
		{
			input: input{
				argument: "3",
			},
			expected: expected{
				value: 3,
			},
		},
		{
			input: input{
				argument: "@",
				parameters: map[string]interface{}{
					"@": "10",
				},
			},
			expected: expected{
				value: 10,
			},
		},
		{
			input: input{
				argument: &plusOperator{arg1: 1, arg2: 2},
			},
			expected: expected{
				value: 3,
			},
		},
		{
			input: input{
				argument: &plusOperator{arg1: nil, arg2: 2},
			},
			expected: expected{
				err: "invalid argument. is nil",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := getInteger(test.input.argument, test.input.parameters)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}
}

func Test_getNumber(t *testing.T) {

	type input struct {
		argument   interface{}
		parameters map[string]interface{}
	}

	type expected struct {
		value float64
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{},
			expected: expected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: input{
				argument: 3,
			},
			expected: expected{
				value: float64(3),
			},
		},
		{
			input: input{
				argument: 3.14,
			},
			expected: expected{
				value: float64(3.14),
			},
		},
		{
			input: input{
				argument: "3.14",
			},
			expected: expected{
				value: float64(3.14),
			},
		},
		{
			input: input{
				argument: "3",
			},
			expected: expected{
				value: float64(3),
			},
		},
		{
			input: input{
				argument: "@",
				parameters: map[string]interface{}{
					"@": "10",
				},
			},
			expected: expected{
				value: float64(10),
			},
		},
		{
			input: input{
				argument: &plusOperator{arg1: 1, arg2: 2},
			},
			expected: expected{
				value: float64(3),
			},
		},
		{
			input: input{
				argument: &plusOperator{arg1: nil, arg2: 2},
			},
			expected: expected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: input{
				argument: true,
			},
			expected: expected{
				err: "invalid argument. expected number",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := getNumber(test.input.argument, test.input.parameters)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}
}

func Test_getBoolean(t *testing.T) {

	currentSelector, _ := newSelectorOperator("@", &ScriptEngine{}, nil)

	type input struct {
		argument   interface{}
		parameters map[string]interface{}
	}

	type expected struct {
		value bool
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				argument: true,
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				argument: "true",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				argument: "invalid",
			},
			expected: expected{
				err: "invalid argument. expected boolean",
			},
		},
		{
			input: input{
				argument: "invalid",
				parameters: map[string]interface{}{
					"invalid": false,
				},
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				argument: &lessThanOperator{},
			},
			expected: expected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: input{
				argument: &lessThanOperator{
					arg1: 1,
					arg2: 2,
				},
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				argument: currentSelector,
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				argument: "null",
				parameters: map[string]interface{}{
					"null": nil,
				},
			},
			expected: expected{
				value: false,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := getBoolean(test.input.argument, test.input.parameters)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}
}

func Test_getString(t *testing.T) {

	type input struct {
		argument   interface{}
		parameters map[string]interface{}
	}

	type expected struct {
		value string
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{},
			expected: expected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: input{
				argument: "'single quotes'",
			},
			expected: expected{
				value: "'single quotes'",
			},
		},
		{
			input: input{
				argument: `"double quotes"`,
			},
			expected: expected{
				value: "'double quotes'",
			},
		},
		{
			input: input{
				argument: `@`,
				parameters: map[string]interface{}{
					"@": "value",
				},
			},
			expected: expected{
				value: "'value'",
			},
		},
		{
			input: input{
				argument: `@`,
			},
			expected: expected{
				value: "@",
			},
		},
		{
			input: input{
				argument: &lessThanOperator{},
			},
			expected: expected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: input{
				argument: &lessThanOperator{arg1: 1, arg2: 2},
			},
			expected: expected{
				value: "true",
			},
		},
		{
			input: input{
				argument: `@`,
				parameters: map[string]interface{}{
					"@": "'value'",
				},
			},
			expected: expected{
				value: "'value'",
			},
		},
		{
			input: input{
				argument: `@`,
				parameters: map[string]interface{}{
					"@": `"value"`,
				},
			},
			expected: expected{
				value: "'value'",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := getString(test.input.argument, test.input.parameters)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}
}
