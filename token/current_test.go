package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test currentToken struct conforms to Token interface
var _ Token = &currentToken{}

func Test_CurrentToken_Type(t *testing.T) {
	assert.Equal(t, "current", (&currentToken{}).Type())
}

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
				tokens: []Token{
					&keyToken{
						key: "name",
					},
				},
			},
			expected: expected{
				value: "first",
			},
		},
	}

	batchTokenTests(t, tests)
}
