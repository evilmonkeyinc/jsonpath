package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test keyToken struct conforms to Token interface
var _ Token = &keyToken{}

func Test_KeyToken_String(t *testing.T) {
	tests := []*tokenStringTest{
		{
			input:    &keyToken{},
			expected: "['']",
		},
		{
			input:    &keyToken{key: "with space"},
			expected: "['with space']",
		},
		{
			input:    &keyToken{key: "nospace"},
			expected: "['nospace']",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_KeyToken_Type(t *testing.T) {
	assert.Equal(t, "key", (&keyToken{}).Type())
}

func Test_KeyToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &keyToken{key: "key"},
			input: input{
				current: nil,
			},
			expected: expected{
				value: nil,
				err:   "key: invalid token target. expected [map] got [nil]",
			},
		},
		{
			token: &keyToken{key: "key"},
			input: input{
				current: "",
			},
			expected: expected{
				value: nil,
				err:   "key: invalid token target. expected [map] got [string]",
			},
		},
		{
			token: &keyToken{key: "key"},
			input: input{
				current: map[string]interface{}{
					"key": true,
				},
			},
			expected: expected{
				value: true,
				err:   "",
			},
		},
		{
			token: &keyToken{key: "missing"},
			input: input{
				current: map[string]interface{}{
					"key": true,
				},
			},
			expected: expected{
				value: nil,
				err:   "key: invalid token key 'missing' not found",
			},
		},
		{
			token: &keyToken{key: "key"},
			input: input{
				current: map[string]interface{}{
					"key": map[string]interface{}{
						"next": "nested target",
					},
				},
				tokens: []Token{
					&keyToken{
						key: "next",
					},
				},
			},
			expected: expected{
				value: "nested target",
				err:   "",
			},
		},
	}

	batchTokenTests(t, tests)
}
