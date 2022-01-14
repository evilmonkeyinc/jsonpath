package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test rootToken struct conforms to Token interface
var _ Token = &rootToken{}

func Test_RootToken_String(t *testing.T) {
	assert.Equal(t, "$", (&rootToken{}).String())
}

func Test_RootToken_Type(t *testing.T) {
	assert.Equal(t, "root", (&rootToken{}).Type())
}

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
