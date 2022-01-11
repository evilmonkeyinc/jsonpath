package token

import (
	"testing"
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
