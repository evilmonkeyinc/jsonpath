package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test expressionToken struct conforms to Token interface
var _ Token = &expressionToken{}

func Test_newExpressionToken(t *testing.T) {
	assert.IsType(t, &expressionToken{}, newExpressionToken("", nil))
}

func Test_ExpressionToken_String(t *testing.T) {

	tests := []*tokenStringTest{
		{
			input:    &expressionToken{expression: ""},
			expected: "()",
		},
		{
			input:    &expressionToken{expression: "1+1"},
			expected: "(1+1)",
		},
		{
			input:    &expressionToken{expression: "true"},
			expected: "(true)",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_ExpressionToken_Type(t *testing.T) {
	assert.Equal(t, "expression", (&expressionToken{}).Type())
}

func Test_ExpressionToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &expressionToken{},
			input: input{},
			expected: expected{
				err: "invalid expression. is empty",
			},
		},
		{
			token: &expressionToken{
				expression: "length",
			},
			input: input{},
			expected: expected{
				err: "invalid expression. eval:1:1: undeclared name: length",
			},
		},
		{
			token: &expressionToken{
				expression: "2*10",
			},
			input: input{
				root:    nil,
				current: nil,
			},
			expected: expected{
				value: int64(20),
			},
		},
		{
			token: &expressionToken{
				expression: "true",
			},
			input: input{
				root:    nil,
				current: nil,
			},
			expected: expected{
				value: true,
			},
		},
		{
			token: &expressionToken{
				expression: "@.length-1",
			},
			input: input{
				root:    nil,
				current: []interface{}{"one", "two", "three"},
			},
			expected: expected{
				value: int64(2),
			},
		},
		{
			token: &expressionToken{
				expression: "'abcdefg'",
			},
			input: input{
				tokens: []Token{
					&indexToken{index: 1, allowString: true},
				},
			},
			expected: expected{
				value: "a",
			},
		},
	}

	batchTokenTests(t, tests)
}

func Test_evaluateExpression(t *testing.T) {

	type input struct {
		root, current interface{}
		expression    string
		options       *Options
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
				err: "invalid expression. is empty",
			},
		},
		{
			input: input{
				expression: "@]",
				current:    true,
			},
			expected: expected{
				err: "invalid expression. eval:1:1: illegal character U+0040 '@'",
			},
		},
		{
			input: input{
				expression: "@[]",
			},
			expected: expected{
				err: "invalid expression. invalid token. '[]' does not match any token format",
			},
		},
		{
			input: input{
				expression: "'key'",
			},
			expected: expected{
				value: "'key'",
			},
		},
		{
			input: input{
				expression: "1--1",
			},
			expected: expected{
				err: "invalid expression. eval:1:2: expected 'EOF', found '--'",
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
				expression: "10/4",
			},
			expected: expected{
				value: int64(2),
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
				expression: "2 > 2",
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				expression: "2 >= 2",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "1 < 2",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "2 < 2",
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				expression: "2 <= 2",
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
		{
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current: map[string]interface{}{
					"price": 5,
				},
				expression: "$.missing < @.price",
			},
			expected: expected{
				err: "invalid expression. key: invalid token key 'missing' not found",
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"expensive": 9.99,
				},
				expression: "$.expensive == 9.99",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"expensive": 9.99,
				},
				expression: "$.expensive == float64(9.99)",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"name": "target",
				},
				expression: "$.name == 'target'",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				root: map[string]interface{}{
					"valid": true,
				},
				expression: "$.valid",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "true && true",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "true && false",
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				expression: "true || true",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "true || false",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "(true || false) && true",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "(true && false) && true",
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				expression: "(true && false) || true",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: "1 + 2 * 3",
			},
			expected: expected{
				value: int64(7),
			},
		},
		{
			input: input{
				expression: "(1 + 2) * 3",
			},
			expected: expected{
				value: int64(9),
			},
		},
		{
			input: input{
				expression: "nil",
			},
			expected: expected{
				value: nil,
			},
		},
		{
			input: input{
				expression: "'key' + '\\'s'",
			},
			expected: expected{
				value: "'key\\'s'",
			},
		},
		{
			input: input{
				expression: `@.key=="hi@example.com"`,
				current: map[string]interface{}{
					"key": "hi@example.com",
				},
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: `"hi@example.com"==@.key`,
				current: map[string]interface{}{
					"key": "hi@example.com",
				},
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: `$.key=="hi@example.com"`,
				root: map[string]interface{}{
					"key": "hi@example.com",
				},
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: `"hi@example.com"==$.key`,
				root: map[string]interface{}{
					"key": "hi@example.com",
				},
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: `"hi@example.com"==$`,
				root:       "hi@example.com",
			},
			expected: expected{
				value: true,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := evaluateExpression(test.input.root, test.input.current, test.input.expression, test.input.options)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}

}
