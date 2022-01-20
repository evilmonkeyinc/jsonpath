package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test rangeToken struct conforms to Token interface
var _ Token = &rangeToken{}

func Test_newRangeToken(t *testing.T) {
	assert.IsType(t, &rangeToken{}, newRangeToken(nil, nil, nil, nil))

	type input struct {
		to, from, step interface{}
		options        *Options
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
				options: &Options{},
			},
			expected: &rangeToken{
				allowMap:    false,
				allowString: false,
			},
		},
		{
			input: input{
				options: &Options{
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
				options: &Options{
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
				options: &Options{
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
				options: &Options{
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

	tests := []*tokenTest{
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
				from: &expressionToken{expression: ""},
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
				from: &expressionToken{expression: "'key'"},
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
				from: &expressionToken{expression: "@.length-1"},
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
				to:   &expressionToken{expression: ""},
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
				to:   &expressionToken{expression: "'key'"},
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
				to:   &expressionToken{expression: "@.length-2"},
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
				step: &expressionToken{expression: ""},
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
				step: &expressionToken{expression: "'key'"},
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
				step: &expressionToken{expression: "@.length-1"},
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
				from: &expressionToken{expression: "nil"},
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
				to:   &expressionToken{expression: "nil"},
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
				step: &expressionToken{expression: "nil"},
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
	}

	batchTokenTests(t, tests)
}

func Test_getRange(t *testing.T) {
	type input struct {
		token            *rangeToken
		obj              interface{}
		start, end, step *int64
	}

	type expected struct {
		obj interface{}
		err string
	}

	testArray := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13}

	intPtr := func(i int64) *int64 {
		return &i
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				token: &rangeToken{},
				obj:   nil,
			},
			expected: expected{
				err: "range: invalid token target. expected [array slice] got [nil]",
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   123,
			},
			expected: expected{
				err: "range: invalid token target. expected [array slice] got [int]",
			},
		},
		{
			input: input{
				token: &rangeToken{allowMap: true},
				obj:   123,
			},
			expected: expected{
				err: "range: invalid token target. expected [array slice map] got [int]",
			},
		},
		{
			input: input{
				token: &rangeToken{allowString: true},
				obj:   123,
			},
			expected: expected{
				err: "range: invalid token target. expected [array slice string] got [int]",
			},
		},
		{
			input: input{
				token: &rangeToken{allowMap: true, allowString: true},
				obj:   123,
			},
			expected: expected{
				err: "range: invalid token target. expected [array slice map string] got [int]",
			},
		},
		{
			input: input{
				token: &rangeToken{allowString: true},
				obj:   "return after this:result text",
				start: intPtr(18),
			},
			expected: expected{
				obj: "result text",
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   testArray,
				start: intPtr(15),
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   testArray,
				end:   intPtr(15),
			},
			expected: expected{
				obj: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   testArray,
				step:  intPtr(0),
			},
			expected: expected{
				err: "range: invalid token out of range",
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   testArray,
			},
			expected: expected{
				obj: testArray,
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   testArray,
				end:   nil,
			},
			expected: expected{
				obj: testArray[0:14],
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   testArray,
				end:   intPtr(-1),
			},
			expected: expected{
				obj: testArray[0:13],
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   testArray,
				end:   intPtr(-3),
			},
			expected: expected{
				obj: testArray[0:11],
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   testArray,
				start: intPtr(-3),
				end:   intPtr(-1),
			},
			expected: expected{
				obj: testArray[11:13],
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				step:  intPtr(2),
			},
			expected: expected{
				obj: []interface{}{"one", "three", "five"},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				start: intPtr(1),
				step:  intPtr(2),
			},
			expected: expected{
				obj: []interface{}{"two", "four"},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				start: intPtr(1),
				end:   intPtr(1),
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				start: intPtr(1),
				end:   intPtr(2),
			},
			expected: expected{
				obj: []interface{}{"two"},
			},
		},
		{
			input: input{
				token: &rangeToken{allowMap: true},
				obj: map[string]interface{}{
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
				obj: []interface{}{
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
			input: input{
				token: &rangeToken{allowMap: true},
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"f": "eff",
					"g": "gee",
					"d": "dee",
				},
				step: intPtr(2),
			},
			expected: expected{
				obj: []interface{}{
					"ae",
					"see",
					"ee",
					"gee",
				},
			},
		},
		{
			input: input{
				token: &rangeToken{allowMap: true},
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"f": "eff",
					"g": "gee",
					"d": "dee",
				},
				start: intPtr(1),
				end:   intPtr(-2),
			},
			expected: expected{
				obj: []interface{}{
					"bee",
					"see",
					"dee",
					"ee",
				},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				step:  intPtr(-1),
			},
			expected: expected{
				obj: []interface{}{"five", "four", "three", "two", "one"},
			},
		},
		{
			input: input{
				token: &rangeToken{allowString: true},
				obj:   "abcdef",
				step:  intPtr(-1),
			},
			expected: expected{
				obj: "fedcba",
			},
		},
		{
			input: input{
				token: &rangeToken{allowMap: true},
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"d": "dee",
				},
				step: intPtr(-1),
			},
			expected: expected{
				obj: []interface{}{"ee", "dee", "see", "bee", "ae"},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				start: intPtr(1),
				end:   intPtr(2),
				step:  intPtr(-1),
			},
			expected: expected{
				obj: []interface{}{"two"},
			},
		},
		{
			input: input{
				token: &rangeToken{allowString: true},
				obj:   "abcdef",
				step:  intPtr(-1),
				start: intPtr(1),
				end:   intPtr(2),
			},
			expected: expected{
				obj: "b",
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"d": "dee",
				},
				start: intPtr(1),
				end:   intPtr(2),
				step:  intPtr(-1),
			},
			expected: expected{
				err: "range: invalid token target. expected [array slice] got [map]",
			},
		},
		{
			input: input{
				token: &rangeToken{allowMap: true},
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"d": "dee",
				},
				start: intPtr(1),
				end:   intPtr(2),
				step:  intPtr(-1),
			},
			expected: expected{
				obj: []interface{}{"bee"},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				step:  intPtr(-1),
				start: intPtr(1),
				end:   intPtr(5),
			},
			expected: expected{
				obj: []interface{}{"five", "four", "three", "two"},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   "abcdef",
				step:  intPtr(-1),
				start: intPtr(1),
				end:   intPtr(5),
			},
			expected: expected{
				err: "range: invalid token target. expected [array slice] got [string]",
			},
		},
		{
			input: input{
				token: &rangeToken{allowString: true},
				obj:   "abcdef",
				step:  intPtr(-1),
				start: intPtr(1),
				end:   intPtr(5),
			},
			expected: expected{
				obj: "edcb",
			},
		},
		{
			input: input{
				token: &rangeToken{allowMap: true},
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"d": "dee",
				},
				step:  intPtr(-1),
				start: intPtr(1),
				end:   intPtr(5),
			},
			expected: expected{
				obj: []interface{}{"ee", "dee", "see", "bee"},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				start: intPtr(-10),
			},
			expected: expected{
				obj: []interface{}{"one", "two", "three", "four", "five"},
			},
		},
		{
			input: input{
				token: &rangeToken{},
				obj:   []string{"one", "two", "three", "four", "five"},
				start: intPtr(0),
				end:   intPtr(-10),
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := test.input.token.getRange(test.input.obj, test.input.start, test.input.end, test.input.step)

			if test.expected.obj == nil {
				assert.Nil(t, obj)
			} else {
				assert.Equal(t, test.expected.obj, obj)
			}

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}
}
