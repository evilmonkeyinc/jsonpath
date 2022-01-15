package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test lengthToken struct conforms to Token interface
var _ Token = &lengthToken{}

func Test_LengthToken_String(t *testing.T) {
	assert.Equal(t, ".length", (&lengthToken{}).String())
}

func Test_LengthToken_Type(t *testing.T) {
	assert.Equal(t, "length", (&lengthToken{}).Type())
}

func Test_LengthToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &lengthToken{},
			input: input{
				current: nil,
			},
			expected: expected{
				err: "length: invalid token target. expected [array map slice string] got [nil]",
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: 1000,
			},
			expected: expected{
				err: "length: invalid token target. expected [array map slice string] got [int]",
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: [3]string{"one", "two", "three"},
			},
			expected: expected{
				value: int64(3),
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: []interface{}{"one", "two", "three", 4, 5},
			},
			expected: expected{
				value: int64(5),
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
				value: int64(3),
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
				value: int64(3),
			},
		},
		{
			token: &lengthToken{},
			input: input{
				current: "this is 26 characters long",
			},
			expected: expected{
				value: int64(26),
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
				value: int64(26),
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
