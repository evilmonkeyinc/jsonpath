package token

import (
	"testing"
)

// Test keyToken struct conforms to Token interface
var _ Token = &keyToken{}

func Test_KeyToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &keyToken{key: "key"},
			input: input{
				current: nil,
			},
			expected: expected{
				value: nil,
				err:   "cannot get key from nil map",
			},
		},
		{
			token: &keyToken{key: "key"},
			input: input{
				current: "",
			},
			expected: expected{
				value: nil,
				err:   "invalid object. expected map",
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
				err:   "'missing' key not found in object",
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
