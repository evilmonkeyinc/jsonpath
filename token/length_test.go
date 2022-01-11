package token

import (
	"testing"
)

// Test lengthToken struct conforms to Token interface
var _ Token = &lengthToken{}

func Test_LengthToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &lengthToken{},
			input: input{
				current: nil,
			},
			expected: expected{
				err: "cannot get elements from nil object",
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: 1000,
			},
			expected: expected{
				err: "invalid object. expected array, map, or string",
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: [3]string{"one", "two", "three"},
			},
			expected: expected{
				value: 3,
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: []interface{}{"one", "two", "three", 4, 5},
			},
			expected: expected{
				value: 5,
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: map[string]int64{
					"one":   1,
					"two":   2,
					"three": 3,
				},
			},
			expected: expected{
				value: 3,
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: map[string]string{
					"one":   "1",
					"two":   "2",
					"three": "3",
				},
			},
			expected: expected{
				value: 3,
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: "this is 26 characters long",
			},
			expected: expected{
				value: 26,
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: "this is 26 characters long",
				tokens: []Token{
					&currentToken{},
				},
			},
			expected: expected{
				value: 26,
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: map[string]string{
					"length": "this would be the length",
				},
			},
			expected: expected{
				value: "this would be the length",
			},
		},
	}

	batchTokenTests(t, tests)

}
