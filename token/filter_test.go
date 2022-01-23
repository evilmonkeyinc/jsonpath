package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test filterToken struct conforms to Token interface
var _ Token = &filterToken{}

func Test_newFilterToken(t *testing.T) {
	assert.IsType(t, &filterToken{}, newFilterToken("", nil, nil))
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

func Test_FilterToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &filterToken{},
			input: input{},
			expected: expected{
				err: "invalid expression. is empty",
			},
		},
		{
			token: &filterToken{
				expression: "nil current",
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{},
				},
			},
			input: input{},
			expected: expected{
				err: "filter: invalid token target. expected [array map slice] got [nil]",
			},
		},
		{
			token: &filterToken{
				expression: "invalid current",
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{},
				},
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
				expression: "fail compile",
				engine: &testEngine{
					err: fmt.Errorf("engine error"),
				},
			},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "invalid expression. engine error",
			},
		},
		{
			token: &filterToken{
				expression: "empty array",
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{},
				},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						err: fmt.Errorf("compiled failed"),
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: true,
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: false,
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: "",
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: "add this",
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: 3.14,
					},
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
				expression: "empty map",
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{},
				},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						err: fmt.Errorf("compiled failed"),
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: true,
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: false,
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: "",
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: "add this",
					},
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
				engine: &testEngine{
					compiledExpression: &testCompiledExpression{
						response: 3.14,
					},
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
				expression: "next is index",
				engine:     &testEngine{compiledExpression: &testCompiledExpression{response: true}},
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
				expression: "next is not index",
				engine:     &testEngine{compiledExpression: &testCompiledExpression{response: true}},
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
	}

	batchTokenTests(t, tests)
}
