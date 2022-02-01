package token

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/stretchr/testify/assert"
)

// Test rangeToken struct conforms to Token interface
var _ Token = &rangeToken{}

func Test_newRangeToken(t *testing.T) {
	assert.IsType(t, &rangeToken{}, newRangeToken(nil, nil, nil, nil))

	type input struct {
		to, from, step interface{}
		options        *option.QueryOptions
	}

	type expected *rangeToken

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				options: nil,
			},
			expected: &rangeToken{
				allowMap:    false,
				allowString: false,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{},
			},
			expected: &rangeToken{
				allowMap:    false,
				allowString: false,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					AllowMapReferenceByIndex:    false,
					AllowStringReferenceByIndex: false,

					AllowMapReferenceByIndexInRange:    true,
					AllowStringReferenceByIndexInRange: true,
				},
			},
			expected: &rangeToken{
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					AllowMapReferenceByIndex:    true,
					AllowStringReferenceByIndex: true,

					AllowMapReferenceByIndexInRange:    false,
					AllowStringReferenceByIndexInRange: false,
				},
			},
			expected: &rangeToken{
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					AllowMapReferenceByIndex:    true,
					AllowStringReferenceByIndex: true,

					AllowMapReferenceByIndexInRange:    true,
					AllowStringReferenceByIndexInRange: true,
				},
			},
			expected: &rangeToken{
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					AllowMapReferenceByIndex:    false,
					AllowStringReferenceByIndex: false,

					AllowMapReferenceByIndexInRange:    false,
					AllowStringReferenceByIndexInRange: true,
				},
			},
			expected: &rangeToken{
				allowMap:    false,
				allowString: true,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := newRangeToken(test.input.to, test.input.from, test.input.step, test.input.options)
			assert.EqualValues(t, test.expected, actual)
		})
	}
}

func Test_RangeToken_String(t *testing.T) {
	tests := []*tokenStringTest{
		{
			input:    &rangeToken{},
			expected: "[:]",
		},
		{
			input:    &rangeToken{from: 1},
			expected: "[1:]",
		},
		{
			input:    &rangeToken{from: 1, to: 2},
			expected: "[1:2]",
		},
		{
			input:    &rangeToken{from: 1, to: &expressionToken{expression: "@.length-1"}},
			expected: "[1:(@.length-1)]",
		},
		{
			input:    &rangeToken{from: 1, to: &expressionToken{expression: "@.length-1"}, step: 2},
			expected: "[1:(@.length-1):2]",
		},
		{
			input:    &rangeToken{from: 1, to: 2, step: 3},
			expected: "[1:2:3]",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_RangeToken_Type(t *testing.T) {
	assert.Equal(t, "range", (&rangeToken{}).Type())
}

func Test_RangeToken_Apply(t *testing.T) {
	batchTokenTests(t, rangeTests)
}

func Benchmark_RangeToken_Apply(b *testing.B) {
	batchTokenBenchmarks(b, rangeTests)
}

var rangeTests = []*tokenTest{
	{
		token: &rangeToken{
			from: nil,
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
				"three",
			},
		},
	},
	{
		token: &rangeToken{
			from: 0,
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
				"three",
			},
		},
	},
	{
		token: &rangeToken{
			from: 1,
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
				"two",
				"three",
			},
		},
	},
	{
		token: &rangeToken{
			from: 1,
			to:   2,
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
				"two",
			},
		},
	},
	{
		token: &rangeToken{
			from: 1,
			step: 3,
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
				"two",
			},
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to:   3,
			step: 2,
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
				"three",
			},
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to:   2,
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
		token: &rangeToken{
			from: 0,
			to:   2,
			step: 2,
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
			},
		},
	},
	{
		token: &rangeToken{
			from: 0,
			step: 2,
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
				"three",
			},
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to:   100,
		},
		input: input{
			current: []interface{}{1, 2, 3},
		},
		expected: expected{
			value: []interface{}{1, 2, 3},
		},
	},
	{
		token: &rangeToken{
			from: "string",
			to:   2,
			step: 2,
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token argument. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to:   "string",
			step: 2,
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token argument. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to:   1,
			step: "string",
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token argument. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: &expressionToken{
				expression:         "",
				compiledExpression: &testCompiledExpression{response: ""},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token invalid expression. is empty",
		},
	},
	{
		token: &rangeToken{
			from: &expressionToken{
				expression:         "'key'",
				compiledExpression: &testCompiledExpression{response: "key"},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: &indexToken{index: 0},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: &expressionToken{
				expression:         "@.length-1",
				compiledExpression: &testCompiledExpression{response: 2},
			},
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
				"three",
			},
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to: &expressionToken{
				expression:         "",
				compiledExpression: &testCompiledExpression{response: ""},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token invalid expression. is empty",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to: &expressionToken{
				expression:         "'key'",
				compiledExpression: &testCompiledExpression{response: "key"},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to:   &indexToken{index: 0},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to: &expressionToken{
				expression:         "@.length-2",
				compiledExpression: &testCompiledExpression{response: 1},
			},
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
			},
		},
	},
	{
		token: &rangeToken{
			from: 0,
			step: &expressionToken{
				expression:         "",
				compiledExpression: &testCompiledExpression{response: ""},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token invalid expression. is empty",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			step: &expressionToken{
				expression:         "'key'",
				compiledExpression: &testCompiledExpression{response: "key"},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			step: &indexToken{index: 0},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [string]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			step: &expressionToken{
				expression:         "@.length-1",
				compiledExpression: &testCompiledExpression{response: 2},
			},
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
				"three",
			},
		},
	},
	{
		token: &rangeToken{
			from: 0,
		},
		input: input{
			tokens: []Token{&keyToken{key: "name"}},
			current: []map[string]interface{}{
				{
					"name": "one",
				},
				{
					"name": "two",
				},
				{
					"name": "three",
				},
				{
					"name": "four",
				},
				{
					"name": "five",
				},
			},
		},
		expected: expected{
			value: []interface{}{
				"one",
				"two",
				"three",
				"four",
				"five",
			},
		},
	},
	{
		token: &rangeToken{
			from: 1,
			to:   -2,
		},
		input: input{
			tokens: []Token{&keyToken{key: "name"}},
			current: []map[string]interface{}{
				{
					"name": "one",
				},
				{
					"name": "two",
				},
				{
					"name": "three",
				},
				{
					"name": "four",
				},
				{
					"name": "five",
				},
			},
		},
		expected: expected{
			value: []interface{}{
				"two",
				"three",
			},
		},
	},
	{
		token: &rangeToken{from: 10},
		input: input{
			current: "this is a substring",
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice] got [string]",
		},
	},
	{
		token: &rangeToken{from: 10, allowString: true},
		input: input{
			current: "this is a substring",
		},
		expected: expected{
			value: "substring",
		},
	},
	{
		token: &rangeToken{from: -9, allowString: true},
		input: input{
			current: "this is a substring",
			tokens: []Token{
				&indexToken{index: 0, allowString: true},
			},
		},
		expected: expected{
			value: "s",
		},
	},
	{
		token: &rangeToken{from: 1},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
				"four",
			},
			tokens: []Token{
				&indexToken{
					index: 0,
				},
			},
		},
		expected: expected{
			value: "two",
		},
	},
	{
		token: &rangeToken{
			from: &expressionToken{
				expression:         "nil",
				compiledExpression: &testCompiledExpression{response: nil},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
				"four",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [nil]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to: &expressionToken{
				expression:         "nil",
				compiledExpression: &testCompiledExpression{response: nil},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
				"four",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [nil]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to:   1,
			step: &expressionToken{
				expression:         "nil",
				compiledExpression: &testCompiledExpression{response: nil},
			},
		},
		input: input{
			current: []interface{}{
				"one",
				"two",
				"three",
				"four",
			},
		},
		expected: expected{
			err: "range: invalid token unexpected expression result. expected [int] got [nil]",
		},
	},
	{
		token: &rangeToken{
			from: 0,
			to:   1,
		},
		input: input{
			current: 123,
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice] got [int]",
		},
	},
	{
		token: &rangeToken{},
		input: input{
			current: nil,
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice] got [nil]",
		},
	},
	{
		token: &rangeToken{},
		input: input{
			current: 123,
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice] got [int]",
		},
	},
	{
		token: &rangeToken{allowMap: true},
		input: input{
			current: 123,
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice map] got [int]",
		},
	},
	{
		token: &rangeToken{allowString: true},
		input: input{
			current: 123,
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice string] got [int]",
		},
	},
	{
		token: &rangeToken{allowMap: true, allowString: true},
		input: input{
			current: 123,
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice map string] got [int]",
		},
	},
	{
		token: &rangeToken{
			from:        18,
			allowString: true,
		},
		input: input{
			current: "return after this:result text",
		},
		expected: expected{
			value: "result text",
		},
	},
	{
		token: &rangeToken{from: 15},
		input: input{
			current: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &rangeToken{to: 15},
		input: input{
			current: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
		expected: expected{
			value: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
	},
	{
		token: &rangeToken{step: 0},
		input: input{
			current: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
		expected: expected{
			err: "range: invalid token out of range",
		},
	},
	{
		token: &rangeToken{},
		input: input{
			current: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
		expected: expected{
			value: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
	},
	{
		token: &rangeToken{to: nil},
		input: input{
			current: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
		expected: expected{
			value: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
	},
	{
		token: &rangeToken{to: -1},
		input: input{
			current: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
		expected: expected{
			value: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve"},
		},
	},
	{
		token: &rangeToken{to: -3},
		input: input{
			current: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
		expected: expected{
			value: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten"},
		},
	},
	{
		token: &rangeToken{
			from: -3,
			to:   -1,
		},
		input: input{
			current: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
		},
		expected: expected{
			value: []interface{}{"eleven", "twelve"},
		},
	},
	{
		token: &rangeToken{
			step: 2,
		},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{"one", "three", "five"},
		},
	},
	{
		token: &rangeToken{from: 1, step: 2},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{"two", "four"},
		},
	},
	{
		token: &rangeToken{from: 1, to: 1},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &rangeToken{
			from: 1,
			to:   2,
		},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{"two"},
		},
	},
	{
		token: &rangeToken{allowMap: true},
		input: input{
			current: map[string]interface{}{
				"b": "bee",
				"a": "ae",
				"c": "see",
				"e": "ee",
				"f": "eff",
				"g": "gee",
				"d": "dee",
			},
		},
		expected: expected{
			value: []interface{}{
				"ae",
				"bee",
				"see",
				"dee",
				"ee",
				"eff",
				"gee",
			},
		},
	},
	{
		token: &rangeToken{allowMap: true, step: 2},
		input: input{
			current: map[string]interface{}{
				"b": "bee",
				"a": "ae",
				"c": "see",
				"e": "ee",
				"f": "eff",
				"g": "gee",
				"d": "dee",
			},
		},
		expected: expected{
			value: []interface{}{
				"ae",
				"see",
				"ee",
				"gee",
			},
		},
	},
	{
		token: &rangeToken{allowMap: true, from: 1, to: -2},
		input: input{
			current: map[string]interface{}{
				"b": "bee",
				"a": "ae",
				"c": "see",
				"e": "ee",
				"f": "eff",
				"g": "gee",
				"d": "dee",
			},
		},
		expected: expected{
			value: []interface{}{
				"bee",
				"see",
				"dee",
				"ee",
			},
		},
	},
	{
		token: &rangeToken{step: -1},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{"five", "four", "three", "two", "one"},
		},
	},
	{
		token: &rangeToken{allowString: true, step: -1},
		input: input{
			current: "abcdef",
		},
		expected: expected{
			value: "fedcba",
		},
	},
	{
		token: &rangeToken{allowMap: true, step: -1},
		input: input{
			current: map[string]interface{}{
				"b": "bee",
				"a": "ae",
				"c": "see",
				"e": "ee",
				"d": "dee",
			},
		},
		expected: expected{
			value: []interface{}{"ee", "dee", "see", "bee", "ae"},
		},
	},
	{
		token: &rangeToken{from: 1, to: 2, step: -1},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{"two"},
		},
	},
	{
		token: &rangeToken{allowString: true, from: 1, to: 2, step: -1},
		input: input{
			current: "abcdef",
		},
		expected: expected{
			value: "b",
		},
	},
	{
		token: &rangeToken{from: 1, to: 2, step: -1},
		input: input{
			current: map[string]interface{}{
				"b": "bee",
				"a": "ae",
				"c": "see",
				"e": "ee",
				"d": "dee",
			},
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice] got [map]",
		},
	},
	{
		token: &rangeToken{allowMap: true, from: 1, to: 2, step: -1},
		input: input{
			current: map[string]interface{}{
				"b": "bee",
				"a": "ae",
				"c": "see",
				"e": "ee",
				"d": "dee",
			},
		},
		expected: expected{
			value: []interface{}{"bee"},
		},
	},
	{
		token: &rangeToken{from: 1, to: 5, step: -1},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{"five", "four", "three", "two"},
		},
	},
	{
		token: &rangeToken{from: 1, to: 5, step: -1},
		input: input{
			current: "abcdef",
		},
		expected: expected{
			err: "range: invalid token target. expected [array slice] got [string]",
		},
	},
	{
		token: &rangeToken{allowString: true, from: 1, to: 5, step: -1},
		input: input{
			current: "abcdef",
		},
		expected: expected{
			value: "edcb",
		},
	},
	{
		token: &rangeToken{allowMap: true, from: 1, to: 5, step: -1},
		input: input{
			current: map[string]interface{}{
				"b": "bee",
				"a": "ae",
				"c": "see",
				"e": "ee",
				"d": "dee",
			},
		},
		expected: expected{
			value: []interface{}{"ee", "dee", "see", "bee"},
		},
	},
	{
		token: &rangeToken{from: -10},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{"one", "two", "three", "four", "five"},
		},
	},
	{
		token: &rangeToken{from: 0, to: -10},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &rangeToken{},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
			tokens:  []Token{&keyToken{"key"}},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &rangeToken{},
		input: input{
			current: []string{"one", "two", "three", "four", "five"},
			tokens:  []Token{&testToken{value: nil}},
		},
		expected: expected{
			value: []interface{}{},
		},
	},
}
