package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ScriptToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &scriptToken{},
			input: input{},
			expected: expected{
				err: "invalid parameter. expression is empty",
			},
		},
		{
			token: &scriptToken{
				expression: "length",
			},
			input: input{},
			expected: expected{
				err: "invalid expression. failed to parse expression. eval:1:1: undeclared name: length",
			},
		},
		{
			token: &scriptToken{
				expression: "2*10",
			},
			input: input{
				root:    nil,
				current: nil,
			},
			expected: expected{
				err: "cannot get index from nil array",
			},
		},
		{
			token: &scriptToken{
				expression: "\"key\"",
			},
			input: input{
				root:    nil,
				current: nil,
			},
			expected: expected{
				err: "cannot get key from nil map",
			},
		},
		{
			token: &scriptToken{
				expression: "true",
			},
			input: input{
				root:    nil,
				current: nil,
			},
			expected: expected{
				err: "unexpected script result. expected integer or string",
			},
		},
		{
			token: &scriptToken{
				expression:       "2*10",
				returnEvaluation: true,
			},
			input: input{
				root:    nil,
				current: "something",
			},
			expected: expected{
				value: int64(20),
			},
		},
		{
			token: &scriptToken{
				expression: "@.length-1",
			},
			input: input{
				root:    nil,
				current: []interface{}{"one", "two", "three"},
			},
			expected: expected{
				value: "three",
			},
		},
	}

	batchTokenTests(t, tests)
}

// TODO : this still needs expanded

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
				expression: "",
			},
			expected: expected{
				err: "invalid parameter. expression is empty",
			},
		},
		{
			input: input{
				expression: "@]",
			},
			expected: expected{
				err: "invalid expression. query must start with '$'",
			},
		},
		{
			input: input{
				expression: "@[]",
			},
			expected: expected{
				err: "invalid expression. invalid token. empty subscript",
			},
		},
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
				err: "failed to parse expression. eval:1:2: expected 'EOF', found '--'",
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
				assert.EqualError(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}

}