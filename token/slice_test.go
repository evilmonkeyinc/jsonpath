package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test sliceToken struct conforms to Token interface
var _ Token = &sliceToken{}

func Test_SliceToken_String(t *testing.T) {
	tests := []*tokenStringTest{
		{
			input:    &sliceToken{number: 0},
			expected: "[:0]",
		},
		{
			input:    &sliceToken{number: 10},
			expected: "[:10]",
		},
		{
			input:    &sliceToken{number: -1},
			expected: "[:-1]",
		},
		{
			input:    &sliceToken{number: &expressionToken{expression: "@.length"}},
			expected: "[:(@.length)]",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_SliceToken_Type(t *testing.T) {
	assert.Equal(t, "slice", (&sliceToken{}).Type())
}

func Test_SliceToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &sliceToken{},
			input: input{},
			expected: expected{
				err: "slice: invalid token argument. expected [int] got [nil]",
			},
		},
		{
			token: &sliceToken{
				number: "1",
			},
			input: input{},
			expected: expected{
				err: "slice: invalid token argument. expected [int] got [string]",
			},
		},
		{
			token: &sliceToken{number: 2},
			input: input{},
			expected: expected{
				err: "slice: invalid token target. expected [array map slice string] got [nil]",
			},
		},
		{
			token: &sliceToken{number: 2},
			input: input{
				current: 123,
			},
			expected: expected{
				err: "slice: invalid token target. expected [array map slice string] got [int]",
			},
		},
		{
			token: &sliceToken{
				number: &expressionToken{expression: ""},
			},
			input: input{},
			expected: expected{
				err: "slice: invalid token invalid expression. is empty",
			},
		},
		{
			token: &sliceToken{
				number: &expressionToken{expression: "\"key\""},
			},
			input: input{},
			expected: expected{
				err: "slice: invalid token unexpected expression result. expected [int] got [string]",
			},
		},
		{
			token: &sliceToken{
				number: &expressionToken{expression: "2"},
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
				},
			},
		},
		{
			token: &sliceToken{
				number: -2,
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
					"four",
					"five",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
				},
			},
		},
		{
			token: &sliceToken{
				number: &expressionToken{expression: "-2"},
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
					"four",
					"five",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
				},
			},
		},
		{
			token: &sliceToken{
				number: 5,
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
					"four",
					"five",
				},
			},
			expected: expected{
				err: "slice: invalid token out of range",
			},
		},
		{
			token: &sliceToken{
				number: 2,
			},
			input: input{
				current: [3]interface{}{
					"one",
					"two",
					"three",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
				},
			},
		},
		{
			token: &sliceToken{
				number: 3,
			},
			input: input{
				current: "substring",
			},
			expected: expected{
				value: "sub",
			},
		},
		{
			token: &sliceToken{
				number: 3,
			},
			input: input{
				current: map[string]string{
					"a": "one",
					"e": "five",
					"d": "four",
					"c": "three",
					"b": "two",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
				},
			},
		},
		{
			token: &sliceToken{
				number: 3,
			},
			input: input{
				current: "substring",
				tokens: []Token{
					&indexToken{index: 1},
				},
			},
			expected: expected{
				value: "u",
			},
		},
		{
			token: &sliceToken{
				number: 2,
			},
			input: input{
				current: []string{
					"one",
					"two",
					"three",
				},
				tokens: []Token{
					&indexToken{index: 1},
				},
			},
			expected: expected{
				value: "two",
			},
		},
		{
			token: &sliceToken{
				number: 2,
			},
			input: input{
				current: []map[string]interface{}{
					{"name": "one"},
					{"name": "two"},
					{"name": "three"},
				},
				tokens: []Token{
					&keyToken{key: "name"},
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
				},
			},
		},
	}

	batchTokenTests(t, tests)

}
