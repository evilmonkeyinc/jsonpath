package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test scriptToken struct conforms to Token interface
var _ Token = &scriptToken{}

func Test_newScriptToken(t *testing.T) {
	assert.IsType(t, &scriptToken{}, newScriptToken("", nil))
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

func Test_ScriptToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &scriptToken{},
			input: input{},
			expected: expected{
				err: "invalid expression. is empty",
			},
		},
		{
			token: &scriptToken{
				expression: "length",
			},
			input: input{},
			expected: expected{
				err: "invalid expression. eval:1:1: undeclared name: length",
			},
		},
		{
			token: &scriptToken{
				expression: "nil",
			},
			input: input{},
			expected: expected{
				err: "unexpected expression result. expected [int string] got [nil]",
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
				err: "index: invalid token target. expected [array slice] got [nil]",
			},
		},
		{
			token: &scriptToken{
				expression: "'key'",
			},
			input: input{
				root:    nil,
				current: nil,
			},
			expected: expected{
				err: "key: invalid token target. expected [map] got [nil]",
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
				err: "unexpected expression result. expected [int string] got [bool]",
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
