package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ScriptToken_Apply(t *testing.T) {

	type input struct {
		expression    string
		root, current interface{}
	}
	type expected struct {
		value interface{}
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{},
			expected: expected{
				err: "invalid parameter. expression is empty",
			},
		},
		{
			input: input{
				expression: "2*10",
				root:       nil,
				current:    nil,
			},
			expected: expected{
				value: int64(20),
			},
		},
		{
			input: input{
				expression: "@.length-1",
				root:       nil,
				current:    []interface{}{"one", "two", "three"},
			},
			expected: expected{
				value: "three",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			token := &scriptToken{expression: test.input.expression}
			actual, err := token.Apply(test.input.root, test.input.current, nil)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}

}

func Test_evaluateExpression(t *testing.T) {

	type input struct {
		root, current interface{}
		expression    string
	}

	type expected struct {
		value interface{}
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				expression: "\"key\"",
			},
			expected: expected{
				value: "\"key\"",
			},
		},
		{
			input: input{
				expression: "1--1",
			},
			expected: expected{
				err: "failed to parse expression : eval:1:2: expected 'EOF', found '--'",
			},
		},
		{
			input: input{
				expression: "1",
			},
			expected: expected{
				value: int64(1),
			},
		},
		{
			input: input{
				expression: "1+1",
			},
			expected: expected{
				value: int64(2),
			},
		},
		{
			input: input{
				expression: "1-1",
			},
			expected: expected{
				value: int64(0),
			},
		},
		{
			input: input{
				expression: "2*2",
			},
			expected: expected{
				value: int64(4),
			},
		},
		{
			input: input{
				expression: "10.0/4",
			},
			expected: expected{
				value: float64(2.5),
			},
		},
		{
			input: input{
				expression: "1==1",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "1 != 1",
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				expression: "1 >2",
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				expression: "1< 2",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "3 % 2",
			},
			expected: expected{
				value: int64(1),
			},
		},
		{
			input: input{
				root:       []string{"one"},
				current:    []string{"two", "three"},
				expression: "$.length",
			},
			expected: expected{
				value: int64(1),
			},
		},
		{
			input: input{
				root:       []string{"one"},
				current:    []string{"two", "three"},
				expression: "@.length",
			},
			expected: expected{
				value: int64(2),
			},
		},
		{
			input: input{
				root:       []string{"one"},
				current:    []string{"two", "three"},
				expression: "@.length-1",
			},
			expected: expected{
				value: int64(1),
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current:    nil,
				expression: "$.expensive",
			},
			expected: expected{
				value: int64(10),
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current:    nil,
				expression: "$.expensive < 10",
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current:    nil,
				expression: "$.expensive > 5",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current: map[string]interface{}{
					"price": 5,
				},
				expression: "$.expensive > @.price",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current: map[string]interface{}{
					"price": 5,
				},
				expression: "$.expensive < @.price",
			},
			expected: expected{
				value: false,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := evaluateExpression(test.input.root, test.input.current, test.input.expression)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}

}
