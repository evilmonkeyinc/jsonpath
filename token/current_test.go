package token

import (
	"testing"
)

// Test currentToken struct conforms to Token interface
var _ Token = &currentToken{}

func Test_CurrentToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &currentToken{},
			input: input{
				current: map[string]interface{}{
					"name": "first",
				},
			},
			expected: expected{
				value: map[string]interface{}{
					"name": "first",
				},
			},
		},
		{
			token: &currentToken{},
			input: input{
				current: map[string]interface{}{
					"name": "first",
				},
				tokens: []Token{&keyToken{
					key: "name",
				}},
			},
			expected: expected{
				value: "first",
			},
		},
	}

	batchTokenTests(t, tests)
}
