package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test rangeToken struct conforms to Token interface
var _ Token = &rangeToken{}

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
				err: "range: invalid token argument. expected [int] got [nil]",
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
				to:   1,
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
					"three",
				},
			},
		},
		{
			token: &rangeToken{
				from: 0,
				to:   1,
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
				to:   1,
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
				current: []interface{}{},
			},
			expected: expected{
				err: "range: invalid token range: invalid token out of range",
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
				from: &expressionToken{expression: "\"key\""},
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
				to:   &expressionToken{expression: "\"key\""},
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
					"two",
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
				step: &expressionToken{expression: "\"key\""},
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
					"four",
				},
			},
		},
		{
			token: &rangeToken{from: 10},
			input: input{
				current: "this is a substring",
			},
			expected: expected{
				value: "substring",
			},
		},
		{
			token: &rangeToken{from: -9},
			input: input{
				current: "this is a substring",
				tokens: []Token{
					&indexToken{
						index: 0,
					},
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
				err: "range: invalid token target. expected [array map slice string] got [int]",
			},
		},
	}

	batchTokenTests(t, tests)
}

func Test_getRange(t *testing.T) {
	type input struct {
		obj       interface{}
		start     int64
		end, step *int64
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
				obj: nil,
			},
			expected: expected{
				err: "range: invalid token target. expected [array map slice string] got [nil]",
			},
		},
		{
			input: input{
				obj: 123,
			},
			expected: expected{
				err: "range: invalid token target. expected [array map slice string] got [int]",
			},
		},
		{
			input: input{
				obj:   "return after this:result text",
				start: 18,
			},
			expected: expected{
				obj: "result text",
			},
		},
		{
			input: input{
				obj:   testArray,
				start: 15,
			},
			expected: expected{
				err: "range: invalid token out of range",
			},
		},
		{
			input: input{
				obj: testArray,
				end: intPtr(15),
			},
			expected: expected{
				err: "range: invalid token out of range",
			},
		},
		{
			input: input{
				obj:  testArray,
				step: intPtr(0),
			},
			expected: expected{
				err: "range: invalid token out of range",
			},
		},
		{
			input: input{
				obj: testArray,
			},
			expected: expected{
				obj: testArray,
			},
		},
		{
			input: input{
				obj: testArray,
				end: intPtr(-1),
			},
			expected: expected{
				obj: testArray[0:14],
			},
		},
		{
			input: input{
				obj: testArray,
				end: intPtr(-3),
			},
			expected: expected{
				obj: testArray[0:12],
			},
		},
		{
			input: input{
				obj:   testArray,
				start: -3,
				end:   intPtr(-1),
			},
			expected: expected{
				obj: testArray[11:14],
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three", "four", "five"},
				step: intPtr(2),
			},
			expected: expected{
				obj: []interface{}{"one", "three", "five"},
			},
		},
		{
			input: input{
				obj:   []string{"one", "two", "three", "four", "five"},
				start: 1,
				step:  intPtr(2),
			},
			expected: expected{
				obj: []interface{}{"two", "four"},
			},
		},
		{
			input: input{
				obj:   []string{"one", "two", "three", "four", "five"},
				start: 1,
				end:   intPtr(1),
			},
			expected: expected{
				obj: []interface{}{"two"},
			},
		},
		{
			input: input{
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
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"f": "eff",
					"g": "gee",
					"d": "dee",
				},
				start: 1,
				end:   intPtr(-2),
			},
			expected: expected{
				obj: []interface{}{
					"bee",
					"see",
					"dee",
					"ee",
					"eff",
				},
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three", "four", "five"},
				step: intPtr(-1),
			},
			expected: expected{
				obj: []interface{}{"five", "four", "three", "two", "one"},
			},
		},
		{
			input: input{
				obj:  "abcdef",
				step: intPtr(-1),
			},
			expected: expected{
				obj: "fedcba",
			},
		},
		{
			input: input{
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
				obj:   []string{"one", "two", "three", "four", "five"},
				start: 1,
				end:   intPtr(1),
				step:  intPtr(-1),
			},
			expected: expected{
				obj: []interface{}{"two"},
			},
		},
		{
			input: input{
				obj:   "abcdef",
				step:  intPtr(-1),
				start: 1,
				end:   intPtr(1),
			},
			expected: expected{
				obj: "b",
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"d": "dee",
				},
				start: 1,
				end:   intPtr(1),
				step:  intPtr(-1),
			},
			expected: expected{
				obj: []interface{}{"bee"},
			},
		},
		{
			input: input{
				obj:   []string{"one", "two", "three", "four", "five"},
				step:  intPtr(-1),
				start: 1,
				end:   intPtr(4),
			},
			expected: expected{
				obj: []interface{}{"five", "four", "three", "two"},
			},
		},
		{
			input: input{
				obj:   "abcdef",
				step:  intPtr(-1),
				start: 1,
				end:   intPtr(4),
			},
			expected: expected{
				obj: "edcb",
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"b": "bee",
					"a": "ae",
					"c": "see",
					"e": "ee",
					"d": "dee",
				},
				step:  intPtr(-1),
				start: 1,
				end:   intPtr(4),
			},
			expected: expected{
				obj: []interface{}{"ee", "dee", "see", "bee"},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := getRange(&rangeToken{}, test.input.obj, test.input.start, test.input.end, test.input.step)

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
