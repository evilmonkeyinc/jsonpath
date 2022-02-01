package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test filterToken struct conforms to Token interface
var _ Token = &filterToken{}

func Test_newFilterToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		actual, err := newFilterToken("", &testEngine{}, nil)
		assert.Nil(t, err)
		assert.IsType(t, &filterToken{}, actual)
	})
	t.Run("failed", func(t *testing.T) {
		actual, err := newFilterToken("", &testEngine{err: fmt.Errorf("failed")}, nil)
		assert.EqualError(t, err, "failed")
		assert.Nil(t, actual)
	})
}

func Test_FilterToken_String(t *testing.T) {
	tests := []*tokenStringTest{
		{
			input:    &filterToken{},
			expected: "[?()]",
		},
		{
			input:    &filterToken{expression: "true"},
			expected: "[?(true)]",
		},
		{
			input:    &filterToken{expression: "@.length<0"},
			expected: "[?(@.length<0)]",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_FilterToken_Type(t *testing.T) {
	assert.Equal(t, "filter", (&filterToken{}).Type())
}

func getNilPointer() *sampleStruct {
	return nil
}

var filterTests = []*tokenTest{
	{
		token: &filterToken{},
		input: input{},
		expected: expected{
			err: "invalid expression. is empty",
		},
	},
	{
		token: &filterToken{
			expression:         "nil current",
			compiledExpression: &testCompiledExpression{},
		},
		input: input{},
		expected: expected{
			err: "filter: invalid token target. expected [array map slice] got [nil]",
		},
	},
	{
		token: &filterToken{
			expression:         "invalid current",
			compiledExpression: &testCompiledExpression{},
		},
		input: input{
			current: "string",
		},
		expected: expected{
			err: "filter: invalid token target. expected [array map slice] got [string]",
		},
	},
	{
		token: &filterToken{
			expression:         "empty array",
			compiledExpression: &testCompiledExpression{},
		},
		input: input{
			current: []interface{}{},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &filterToken{
			expression: "failed evaluate array",
			compiledExpression: &testCompiledExpression{
				err: fmt.Errorf("compiled failed"),
			},
		},
		input: input{
			current: [3]interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &filterToken{
			expression: "true evaluate array",
			compiledExpression: &testCompiledExpression{
				response: true,
			},
		},
		input: input{
			current: [3]interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
		},
	},
	{
		token: &filterToken{
			expression: "false evaluate array",
			compiledExpression: &testCompiledExpression{
				response: false,
			},
		},
		input: input{
			current: [3]interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &filterToken{
			expression: "empty string evaluate array",
			compiledExpression: &testCompiledExpression{
				response: "",
			},
		},
		input: input{
			current: [3]interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &filterToken{
			expression: "non-empty string evaluate array",
			compiledExpression: &testCompiledExpression{
				response: "add this",
			},
		},
		input: input{
			current: [3]interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
		},
	},
	{
		token: &filterToken{
			expression: "other evaluate array",
			compiledExpression: &testCompiledExpression{
				response: 3.14,
			},
		},
		input: input{
			current: [3]interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
		},
	},
	{
		token: &filterToken{
			expression:         "empty map",
			compiledExpression: &testCompiledExpression{},
		},
		input: input{
			current: map[string]interface{}{},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &filterToken{
			expression: "failed evaluate map",
			compiledExpression: &testCompiledExpression{
				err: fmt.Errorf("compiled failed"),
			},
		},
		input: input{
			current: map[string]interface{}{"key": "value"},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &filterToken{
			expression: "true evaluate map",
			compiledExpression: &testCompiledExpression{
				response: true,
			},
		},
		input: input{
			current: map[string]interface{}{"key": "value"},
		},
		expected: expected{
			value: []interface{}{"value"},
		},
	},
	{
		token: &filterToken{
			expression: "false evaluate map",
			compiledExpression: &testCompiledExpression{
				response: false,
			},
		},
		input: input{
			current: map[string]interface{}{"key": "value"},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &filterToken{
			expression: "empty string evaluate map",
			compiledExpression: &testCompiledExpression{
				response: "",
			},
		},
		input: input{
			current: map[string]interface{}{"key": "value"},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &filterToken{
			expression: "non-empty string evaluate map",
			compiledExpression: &testCompiledExpression{
				response: "add this",
			},
		},
		input: input{
			current: map[string]interface{}{"key": "value"},
		},
		expected: expected{
			value: []interface{}{"value"},
		},
	},
	{
		token: &filterToken{
			expression: "other evaluate map",
			compiledExpression: &testCompiledExpression{
				response: 3.14,
			},
		},
		input: input{
			current: map[string]interface{}{"key": "value"},
		},
		expected: expected{
			value: []interface{}{"value"},
		},
	},
	{
		token: &filterToken{
			expression:         "next is index",
			compiledExpression: &testCompiledExpression{response: true},
		},
		input: input{
			current: []interface{}{1, 2, 3, 4, 5},
			tokens:  []Token{&indexToken{index: 1}},
		},
		expected: expected{
			value: 2,
		},
	},
	{
		token: &filterToken{
			expression:         "next is not index",
			compiledExpression: &testCompiledExpression{response: true},
		},
		input: input{
			current: []interface{}{
				map[string]interface{}{"key": "one"},
				map[string]interface{}{"key": "two"},
				map[string]interface{}{"key": "three"},
			},
			tokens: []Token{&keyToken{key: "key"}},
		},
		expected: expected{
			value: []interface{}{"one", "two", "three"},
		},
	},
	{
		token: &filterToken{
			expression: "array",
			compiledExpression: &testCompiledExpression{
				response: [1]string{"one"},
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "slice",
			compiledExpression: &testCompiledExpression{
				response: []string{"one"},
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "empty array",
			compiledExpression: &testCompiledExpression{
				response: [0]string{},
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "map",
			compiledExpression: &testCompiledExpression{
				response: map[string]interface{}{"key": "value"},
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "empty map",
			compiledExpression: &testCompiledExpression{
				response: map[string]interface{}{},
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "nil pointer",
			compiledExpression: &testCompiledExpression{
				response: getNilPointer(),
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "single quotes empty",
			compiledExpression: &testCompiledExpression{
				response: "''",
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "single quotes not empty",
			compiledExpression: &testCompiledExpression{
				response: "' '",
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "double quotes empty",
			compiledExpression: &testCompiledExpression{
				response: `""`,
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{},
			err:   "",
		},
	},
	{
		token: &filterToken{
			expression: "double quotes not empty",
			compiledExpression: &testCompiledExpression{
				response: `" "`,
			},
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
			err:   "",
		},
	},
}

func Test_FilterToken_Apply(t *testing.T) {
	batchTokenTests(t, filterTests)
}

func Benchmark_FilterToken_Apply(b *testing.B) {
	batchTokenBenchmarks(b, filterTests)
}
