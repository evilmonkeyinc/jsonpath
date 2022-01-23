package standard

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/script"
	"github.com/stretchr/testify/assert"
)

// Test ScriptEngine struct conforms to Engine interface
var _ script.Engine = &ScriptEngine{}

func Test_ScriptEngine_Compile(t *testing.T) {

	engine := &ScriptEngine{}

	type input struct {
		expression string
		options    *option.QueryOptions
	}

	type expected struct {
		compiled script.CompiledExpression
		err      string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				expression: "1 * 2 + 3",
			},
			expected: expected{
				compiled: &compiledExpression{
					expression: "1 * 2 + 3",
					rootOperator: &plusOperator{
						arg1: &multiplyOperator{
							arg1: "1",
							arg2: "2",
						},
						arg2: "3",
					},
					engine:  engine,
					options: nil,
				},
			},
		},
		{
			input: input{
				expression: "123",
			},
			expected: expected{
				compiled: &compiledExpression{
					expression:   "123",
					rootOperator: nil,
					engine:       engine,
					options:      nil,
				},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {

			actual, err := engine.Compile(test.input.expression, test.input.options)
			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
			assert.Equal(t, test.expected.compiled, actual)
		})
	}
}

func Test_ScriptEngine_buildOperators(t *testing.T) {

	engine := &ScriptEngine{}
	currentLengthSelector, _ := newSelectorOperator("@.length", engine, nil)
	currentKey, _ := newSelectorOperator("@.key", engine, nil)

	currentA, _ := newSelectorOperator("@.a", engine, nil)
	currentB, _ := newSelectorOperator("@.b", engine, nil)
	currentC, _ := newSelectorOperator("@.c", engine, nil)
	currentRangeZeroOne, _ := newSelectorOperator("@[0:1]", engine, nil)
	currentWildcard, _ := newSelectorOperator("@.*", engine, nil)

	type input struct {
		expression string
		tokens     []string
		options    *option.QueryOptions
	}

	type expected struct {
		operator operator
		err      string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				expression: " ",
				tokens:     []string{"+"},
			},
			expected: expected{
				operator: nil,
				err:      "",
			},
		},
		{
			input: input{
				expression: "@.length",
				tokens:     []string{"@"},
			},
			expected: expected{
				operator: currentLengthSelector,
				err:      "",
			},
		},
		{
			input: input{
				expression: "1=~\\d",
				tokens:     []string{"=~"},
			},
			expected: expected{
				operator: &regexOperator{arg1: "1", arg2: "\\d"},
				err:      "",
			},
		},
		{
			input: input{
				expression: "1||1",
				tokens:     []string{"||"},
			},
			expected: expected{
				operator: &orOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1&&1",
				tokens:     []string{"&&"},
			},
			expected: expected{
				operator: &andOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1==1",
				tokens:     []string{"=="},
			},
			expected: expected{
				operator: &equalsOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1!=1",
				tokens:     []string{"!="},
			},
			expected: expected{
				operator: &notEqualsOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1<=1",
				tokens:     []string{"<="},
			},
			expected: expected{
				operator: &lessThanOrEqualOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1<1",
				tokens:     []string{"<"},
			},
			expected: expected{
				operator: &lessThanOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1>=1",
				tokens:     []string{">="},
			},
			expected: expected{
				operator: &greaterThanOrEqualOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1>1",
				tokens:     []string{">"},
			},
			expected: expected{
				operator: &greaterThanOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1+1",
				tokens:     []string{"+"},
			},
			expected: expected{
				operator: &plusOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1-1",
				tokens:     []string{"-"},
			},
			expected: expected{
				operator: &subtractOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1**1",
				tokens:     []string{"**"},
			},
			expected: expected{
				operator: &powerOfOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1*1",
				tokens:     []string{"*"},
			},
			expected: expected{
				operator: &multiplyOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1/1",
				tokens:     []string{"/"},
			},
			expected: expected{
				operator: &divideOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1%1",
				tokens:     []string{"%"},
			},
			expected: expected{
				operator: &modulusOperator{
					arg1: "1",
					arg2: "1",
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "1===1",
				tokens:     []string{"==="},
			},
			expected: expected{
				err: "unsupported operator",
			},
		},
		{
			input: input{
				expression: "1+1*1",
				tokens:     []string{"*", "/", "+", "-"},
			},
			expected: expected{
				operator: &multiplyOperator{
					arg1: &plusOperator{
						arg1: "1",
						arg2: "1",
					},
					arg2: "1",
				},
			},
		},
		{
			input: input{
				expression: "1+1+1",
				tokens:     []string{"*", "/", "+", "-"},
			},
			expected: expected{
				operator: &plusOperator{
					arg1: "1",
					arg2: &plusOperator{
						arg1: "1",
						arg2: "1",
					},
				},
			},
		},
		{
			input: input{
				expression: "1+1===1",
				tokens:     []string{"+", "==="},
			},
			expected: expected{
				err: "unsupported operator",
			},
		},
		{
			input: input{
				expression: "1===1+1",
				tokens:     []string{"+", "==="},
			},
			expected: expected{
				err: "unsupported operator",
			},
		},
		{
			input: input{
				expression: `@.a && (@.b || @.c)`,
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: &andOperator{
					arg1: currentA,
					arg2: &orOperator{
						arg1: currentB,
						arg2: currentC,
					},
				},
			},
		},
		{
			input: input{
				expression: "@[0:1]==[1]",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: &equalsOperator{
					arg1: currentRangeZeroOne,
					arg2: []interface{}{float64(1)},
				},
			},
		},
		{
			input: input{
				expression: "@.*==[1,2]",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: &equalsOperator{
					arg1: currentWildcard,
					arg2: []interface{}{float64(1), float64(2)},
				},
			},
		},
		{
			input: input{
				expression: "@.key<3",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: &lessThanOperator{
					arg1: currentKey,
					arg2: "3",
				},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := engine.buildOperators(test.input.expression, test.input.tokens, test.input.options)
			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
			assert.Equal(t, test.expected.operator, actual)
		})
	}
}

func Test_ScriptEngine_Evaluate(t *testing.T) {

	engine := &ScriptEngine{}

	type input struct {
		root, current interface{}
		expression    string
		options       *option.QueryOptions
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
				root:       "root",
				current:    "current",
				expression: "nil",
			},
			expected: expected{
				value: nil,
			},
		},
		{
			input: input{
				root:       "root",
				current:    "current",
				expression: "null",
			},
			expected: expected{
				value: nil,
			},
		},
		{
			input: input{
				root:       "root",
				current:    "current",
				expression: "$",
			},
			expected: expected{
				value: "root",
			},
		},
		{
			input: input{
				root:       "root",
				current:    "current",
				expression: "@",
			},
			expected: expected{
				value: "current",
			},
		},
		{
			input: input{
				root:       "root",
				current:    "current",
				expression: "other",
			},
			expected: expected{
				value: "other",
			},
		},
		{
			input: input{
				expression: "1+fish",
			},
			expected: expected{
				err: "invalid argument. expected number",
			},
		},
		{
			input: input{
				expression: "@[-1]==2",
				current:    []interface{}{0, 2},
			},
			expected: expected{
				value: true,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := engine.Evaluate(test.input.root, test.input.current, test.input.expression, test.input.options)
			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
			assert.Equal(t, test.expected.value, actual)
		})
	}

}
