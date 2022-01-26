package standard

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_compiledExpression_Evaluate(t *testing.T) {

	compile := func(engine *ScriptEngine, expression string) *compiledExpression {
		generic, _ := engine.Compile(expression, nil)
		specific, _ := generic.(*compiledExpression)
		return specific
	}

	engine := &ScriptEngine{}
	rootExpression := compile(engine, "$")
	currentExpression := compile(engine, "@")
	currentKeyExpression := compile(engine, "@.key")
	nilExpression := compile(engine, "nil")
	nullExpression := compile(engine, "null")

	type input struct {
		compiled      *compiledExpression
		root, current interface{}
	}
	type expected struct {
		err   string
		value interface{}
	}

	tests := []struct {
		input
		expected
	}{
		{
			input: input{
				compiled: &compiledExpression{
					expression: "",
					engine:     engine,
				},
			},
			expected: expected{
				err: "invalid expression. is empty",
			},
		},
		{
			input: input{
				compiled: &compiledExpression{
					expression: "true",
					engine:     engine,
				},
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				compiled: &compiledExpression{
					expression: "false",
					engine:     engine,
				},
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				compiled: &compiledExpression{
					expression: "3.14",
					engine:     engine,
				},
			},
			expected: expected{
				value: float64(3.14),
			},
		},
		{
			input: input{
				compiled: &compiledExpression{
					expression: "3",
					engine:     engine,
				},
			},
			expected: expected{
				value: float64(3),
			},
		},
		{
			input: input{
				compiled: &compiledExpression{
					expression: "[]",
					engine:     engine,
				},
			},
			expected: expected{
				value: "[]",
			},
		},
		{
			input: input{
				compiled: rootExpression,
				root:     123,
			},
			expected: expected{
				value: int(123),
			},
		},
		{
			input: input{
				compiled: currentExpression,
				current:  int64(321),
			},
			expected: expected{
				value: int64(321),
			},
		},
		{
			input: input{
				compiled: nilExpression,
			},
			expected: expected{
				value: nil,
			},
		},
		{
			input: input{
				compiled: nullExpression,
			},
			expected: expected{
				value: nil,
			},
		},
		{
			input: input{
				compiled: currentKeyExpression,
				current:  map[string]interface{}{},
			},
			expected: expected{
				err: "key: invalid token key 'key' not found",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := test.input.compiled.Evaluate(test.root, test.current)
			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
			assert.Equal(t, test.expected.value, actual)
		})
	}
}
