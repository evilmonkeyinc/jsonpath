package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test scriptToken struct conforms to Token interface
var _ Token = &scriptToken{}

func Test_newScriptToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		actual, err := newScriptToken("", &testEngine{}, nil)
		assert.Nil(t, err)
		assert.IsType(t, &scriptToken{}, actual)
	})
	t.Run("failed", func(t *testing.T) {
		actual, err := newScriptToken("", &testEngine{err: fmt.Errorf("failed")}, nil)
		assert.EqualError(t, err, "failed")
		assert.Nil(t, actual)
	})
}

func Test_ScriptToken_String(t *testing.T) {

	tests := []*tokenStringTest{
		{
			input:    &scriptToken{expression: ""},
			expected: "[()]",
		},
		{
			input:    &scriptToken{expression: "1+1"},
			expected: "[(1+1)]",
		},
		{
			input:    &scriptToken{expression: "true"},
			expected: "[(true)]",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_ScriptToken_Type(t *testing.T) {
	assert.Equal(t, "script", (&scriptToken{}).Type())
}

var scriptTests = []*tokenTest{
	{
		token: &scriptToken{},
		input: input{},
		expected: expected{
			err: "invalid expression. is empty",
		},
	},
	{
		token: &scriptToken{
			expression:         "engine error",
			compiledExpression: &testCompiledExpression{err: fmt.Errorf("engine error")},
		},
		input: input{},
		expected: expected{
			err: "invalid expression. engine error",
		},
	},
	{
		token: &scriptToken{
			expression:         "nil response",
			compiledExpression: &testCompiledExpression{response: nil},
		},
		input: input{},
		expected: expected{
			err: "unexpected expression result. expected [int string] got [nil]",
		},
	},
	{
		token: &scriptToken{
			expression:         "bool response",
			compiledExpression: &testCompiledExpression{response: true},
		},
		input: input{},
		expected: expected{
			err: "unexpected expression result. expected [int string] got [bool]",
		},
	},
	{
		token: &scriptToken{
			expression:         "string response",
			compiledExpression: &testCompiledExpression{response: "key"},
		},
		input: input{
			current: map[string]interface{}{
				"key": "value",
			},
		},
		expected: expected{
			value: "value",
		},
	},
	{
		token: &scriptToken{
			expression:         "int response",
			compiledExpression: &testCompiledExpression{response: 1},
		},
		input: input{
			current: []string{"one", "two", "three"},
		},
		expected: expected{
			value: "two",
		},
	},
}

func Test_ScriptToken_Apply(t *testing.T) {
	batchTokenTests(t, scriptTests)
}

func Benchmark_ScriptToken_Apply(b *testing.B) {
	batchTokenBenchmarks(b, scriptTests)
}
