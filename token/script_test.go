package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test scriptToken struct conforms to Token interface
var _ Token = &scriptToken{}

func Test_newScriptToken(t *testing.T) {
	assert.IsType(t, &scriptToken{}, newScriptToken("", nil, nil))
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
				expression: "engine error",
				engine:     &testEngine{err: fmt.Errorf("engine error")},
			},
			input: input{},
			expected: expected{
				err: "invalid expression. engine error",
			},
		},
		{
			token: &scriptToken{
				expression: "nil response",
				engine:     &testEngine{response: nil},
			},
			input: input{},
			expected: expected{
				err: "unexpected expression result. expected [int string] got [nil]",
			},
		},
		{
			token: &scriptToken{
				expression: "bool response",
				engine:     &testEngine{response: true},
			},
			input: input{},
			expected: expected{
				err: "unexpected expression result. expected [int string] got [bool]",
			},
		},
		{
			token: &scriptToken{
				expression: "string response",
				engine:     &testEngine{response: "key"},
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
				expression: "int response",
				engine:     &testEngine{response: 1},
			},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				value: "two",
			},
		},
	}

	batchTokenTests(t, tests)
}
