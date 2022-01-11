package token

import (
	"testing"
)

// Test rootToken struct conforms to Token interface
var _ Token = &rootToken{}

func Test_RootToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &rootToken{},
			input: input{
				root: map[string]interface{}{
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
			token: &rootToken{},
			input: input{
				root: map[string]interface{}{
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
