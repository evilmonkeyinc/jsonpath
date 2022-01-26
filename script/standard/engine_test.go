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
		{
			input: input{
				expression: ".$",
			},
			expected: expected{
				compiled: &compiledExpression{
					expression:   ".$",
					rootOperator: nil,
					engine:       engine,
					options:      nil,
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "$[]",
			},
			expected: expected{
				err: "invalid token. '[]' does not match any token format",
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
				value: "'root'",
			},
		},
		{
			input: input{
				root:       "root",
				current:    "current",
				expression: "@",
			},
			expected: expected{
				value: "'current'",
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
		{
			input: input{
				expression: "@.name=~'hello.*'",
				current: map[string]interface{}{
					"name": "hello world",
				},
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				expression: ".@.name=~'hello.*'",
				current: map[string]interface{}{
					"name": "hello world",
				},
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				expression: "@[]=~'hello.*'",
				current: map[string]interface{}{
					"name": "hello world",
				},
			},
			expected: expected{
				err: "invalid token. '[]' does not match any token format",
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
		{
			input: input{
				expression: "@.key==true",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: &equalsOperator{
					arg1: currentKey,
					arg2: "true",
				},
			},
		},
		{
			input: input{
				expression: "@.key=~'hello.*'",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: &regexOperator{
					arg1: currentKey,
					arg2: "'hello.*'",
				},
			},
		},
		{
			input: input{
				expression: ".$",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: nil,
				err:      "",
			},
		},
		{
			input: input{
				expression: "$[]",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: nil,
				err:      "invalid token. '[]' does not match any token format",
			},
		},
		{
			input: input{
				expression: ".||",
				tokens:     []string{"||"},
			},
			expected: expected{
				operator: nil,
				err:      "",
			},
		},
		{
			input: input{
				expression: "||.",
				tokens:     []string{"||"},
			},
			expected: expected{
				operator: nil,
				err:      "",
			},
		},
		{
			input: input{
				expression: ".||.",
				tokens:     []string{"||"},
			},
			expected: expected{
				operator: &orOperator{arg1: ".", arg2: "."},
				err:      "",
			},
		},
		{
			input: input{
				expression: "||.",
				tokens:     []string{"||", "&&"},
			},
			expected: expected{
				operator: nil,
				err:      "",
			},
		},
		{
			input: input{
				expression: "!true",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: &notOperator{arg: "true"},
				err:      "",
			},
		},
		{
			input: input{
				expression: "!(0==1)",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: &notOperator{
					arg: &equalsOperator{arg1: "0", arg2: "1"},
				},
				err: "",
			},
		},
		{
			input: input{
				expression: "true!true",
				tokens:     defaultTokens,
			},
			expected: expected{
				operator: nil,
				err:      "",
			},
		},
		{
			input: input{
				expression: "true!true",
				tokens:     []string{"!"},
			},
			expected: expected{
				operator: nil,
				err:      "",
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

func Test_ScriptEngine_parseArgument(t *testing.T) {

	engine := &ScriptEngine{}

	type input struct {
		argument string
		tokens   []string
		options  *option.QueryOptions
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
				argument: ".$",
				tokens:   defaultTokens,
			},
			expected: expected{
				value: ".$",
			},
		},
		{
			input: input{
				argument: "$[]",
				tokens:   defaultTokens,
			},
			expected: expected{
				err: "invalid token. '[]' does not match any token format",
			},
		},
		{
			input: input{
				argument: "1||1",
				tokens:   defaultTokens,
			},
			expected: expected{
				value: &orOperator{arg1: "1", arg2: "1"},
			},
		},
		{
			input: input{
				argument: "[1,2,3]",
				tokens:   defaultTokens,
			},
			expected: expected{
				value: []interface{}{
					float64(1),
					float64(2),
					float64(3),
				},
			},
		},
		{
			input: input{
				argument: `{"key":"value"}`,
				tokens:   defaultTokens,
			},
			expected: expected{
				value: map[string]interface{}{
					"key": "value",
				},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := engine.parseArgument(test.input.argument, test.input.tokens, test.input.options)
			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
			assert.Equal(t, test.expected.value, actual)
		})
	}
}

func Test_Evaluate(t *testing.T) {

	tests := []struct {
		expression string
		expected   interface{}
	}{
		{
			expression: "true || false",
			expected:   true,
		},
		{
			expression: "true && true",
			expected:   true,
		},
		{
			expression: "false || true && true",
			expected:   true,
		},
		{
			expression: "true && true || false",
			expected:   true,
		},
		{
			expression: "true == true",
			expected:   true,
		},
		{
			expression: "'true' == true",
			expected:   false,
		},
		{
			expression: "'true' == 'true'",
			expected:   true,
		},
		{
			expression: "true != true",
			expected:   false,
		},
		{
			expression: "'true' != true",
			expected:   true,
		},
		{
			expression: "'true' != 'true'",
			expected:   false,
		},
		{
			expression: "1 < 2",
			expected:   true,
		},
		{
			expression: "2 < 2",
			expected:   false,
		},
		{
			expression: "1 <= 2",
			expected:   true,
		},
		{
			expression: "2 <= 2",
			expected:   true,
		},
		{
			expression: "2 > 1",
			expected:   true,
		},
		{
			expression: "2 > 2",
			expected:   false,
		},
		{
			expression: "2 >= 1",
			expected:   true,
		},
		{
			expression: "2 >= 2",
			expected:   true,
		},
		{
			expression: "1 + 2",
			expected:   float64(3),
		},
		{
			expression: "2 - 1",
			expected:   float64(1),
		},
		{
			expression: "1 + 2 - 3",
			expected:   float64(0),
		},
		{
			expression: "1 * 2",
			expected:   float64(2),
		},
		{
			expression: "1 / 2",
			expected:   float64(0.5),
		},
		{
			expression: "1 * 2 / 8",
			expected:   float64(0.25),
		},
		{
			expression: "2 ** 0",
			expected:   float64(1),
		},
		{
			expression: "2 ** 1",
			expected:   float64(2),
		},
		{
			expression: "2 ** 2",
			expected:   float64(4),
		},
		{
			expression: "4 % 2",
			expected:   int64(0),
		},
		{
			expression: "4 % 3",
			expected:   int64(1),
		},
		{
			expression: "15/5*(8-6+3)*5",
			expected:   float64(75),
		},
		{
			expression: "8+3-3*6+2*(2*3)",
			expected:   float64(5),
		},
		{
			expression: "!false",
			expected:   true,
		},
		{
			expression: "!true",
			expected:   false,
		},
		{
			expression: "!(1==1)",
			expected:   false,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			engine := &ScriptEngine{}
			actual, err := engine.Evaluate(nil, nil, test.expression, nil)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}

}
