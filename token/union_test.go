package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test unionToken struct conforms to Token interface
var _ Token = &unionToken{}

func Test_UnionToken_String(t *testing.T) {
	tests := []*tokenStringTest{
		{
			input:    &unionToken{},
			expected: "[]",
		},
		{
			input:    &unionToken{arguments: []interface{}{"one"}},
			expected: "['one']",
		},
		{
			input:    &unionToken{arguments: []interface{}{1, 3, 4}},
			expected: "[1,3,4]",
		},
		{
			input:    &unionToken{arguments: []interface{}{1, &expressionToken{expression: "4%2"}, "last"}},
			expected: "[1,(4%2),'last']",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_UnionToken_Type(t *testing.T) {
	assert.Equal(t, "union", (&unionToken{}).Type())
}

func Test_UnionToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &unionToken{},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "union: invalid token argument. expected [array slice] got [nil]",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					&expressionToken{expression: "nil"},
				},
			},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "union: invalid token argument. expected [int string] got [nil]",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{"one", 2},
			},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "union: invalid token argument. expected [string] got [int]",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{2, "one"},
			},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "union: invalid token argument. expected [int] got [string]",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{3.14},
			},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "union: invalid token argument. expected [int string] got [float64]",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					&expressionToken{
						expression: "",
					},
					"one",
				},
			},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "union: invalid token invalid expression. is empty",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					&expressionToken{
						expression: "1+1",
					},
					"one",
				},
			},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "union: invalid token argument. expected [int] got [string]",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					&expressionToken{
						expression: "1+1",
					},
					3,
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
				value: []interface{}{
					"three",
					"four",
				},
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					0,
					3,
					4,
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
				err: "union: invalid token out of range",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					0,
					3,
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
				value: []interface{}{
					"one",
					"four",
				},
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					"a",
					"d",
				},
				allowMapIndex: true,
			},
			input: input{
				current: map[string]interface{}{
					"a": "one",
					"b": "two",
					"c": "three",
					"d": "four",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"four",
				},
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					"a",
					"d",
					"e",
				},
				allowMapIndex: true,
			},
			input: input{
				current: map[string]interface{}{
					"a": "one",
					"b": "two",
					"c": "three",
					"d": "four",
				},
			},
			expected: expected{
				err: "union: invalid token key 'e' not found",
			},
		},
		{
			token: &unionToken{
				arguments:        []interface{}{0, 2, 4},
				allowStringIndex: true,
			},
			input: input{
				current: "abcdefghijkl",
			},
			expected: expected{
				value: "ace",
			},
		},
		{
			token: &unionToken{
				arguments:        []interface{}{0, 2, 4},
				allowStringIndex: true,
			},
			input: input{
				current: "abcdefghijkl",
				tokens: []Token{
					&indexToken{index: 1, allowString: true},
				},
			},
			expected: expected{
				value: "c",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					0,
					3,
				},
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
					"four",
				},
				tokens: []Token{
					&indexToken{
						index: -1,
					},
				},
			},
			expected: expected{
				value: "four",
			},
		},
		{
			token: &unionToken{
				arguments: []interface{}{
					0,
					3,
				},
			},
			input: input{
				current: []map[string]interface{}{
					{"name": "one"},
					{"name": "two"},
					{"name": "three"},
					{"name": "four"},
				},
				tokens: []Token{
					&keyToken{
						key: "name",
					},
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"four",
				},
			},
		},
	}

	batchTokenTests(t, tests)
}

func Test_UnionToken_getUnionByIndex(t *testing.T) {

	type input struct {
		token *unionToken
		obj   interface{}
		keys  []int64
	}

	type expected struct {
		obj interface{}
		err string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				token: &unionToken{},
				obj:   nil,
			},
			expected: expected{
				err: "union: invalid token target. expected [array slice] got [nil]",
			},
		},
		{
			input: input{
				token: &unionToken{
					allowMapIndex: true,
				},
				obj: nil,
			},
			expected: expected{
				err: "union: invalid token target. expected [array slice map] got [nil]",
			},
		},
		{
			input: input{
				token: &unionToken{
					allowStringIndex: true,
				},
				obj: nil,
			},
			expected: expected{
				err: "union: invalid token target. expected [array slice string] got [nil]",
			},
		},
		{
			input: input{
				token: &unionToken{
					allowMapIndex:    true,
					allowStringIndex: true,
				},
				obj: nil,
			},
			expected: expected{
				err: "union: invalid token target. expected [array slice map string] got [nil]",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   123,
			},
			expected: expected{
				err: "union: invalid token target. expected [array slice] got [int]",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{4},
			},
			expected: expected{
				err: "union: invalid token out of range",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{-10},
			},
			expected: expected{
				err: "union: invalid token out of range",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{-1, -2},
			},
			expected: expected{
				obj: []interface{}{"three", "two"},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{0, 2},
			},
			expected: expected{
				obj: []interface{}{"three", "one"},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []interface{}{"one", "two", 3},
				keys:  []int64{0, 2},
			},
			expected: expected{
				obj: []interface{}{"one", 3},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   [3]int64{1, 2, 3},
				keys:  []int64{0, 2},
			},
			expected: expected{
				obj: []interface{}{
					int64(1),
					int64(3),
				},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   "abcdefghijklmnopqrstuvwxyz",
				keys:  []int64{0, 2, 4},
			},
			expected: expected{
				err: "union: invalid token target. expected [array slice] got [string]",
			},
		},
		{
			input: input{
				token: &unionToken{allowStringIndex: true},
				obj:   "abcdefghijklmnopqrstuvwxyz",
				keys:  []int64{0, 2, 4},
			},
			expected: expected{
				obj: "ace",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "one",
					"d": "four",
					"e": "five",
					"c": "three",
					"b": "two",
				},
				keys: []int64{0, 1, 3},
			},
			expected: expected{
				err: "union: invalid token target. expected [array slice] got [map]",
			},
		},
		{
			input: input{
				token: &unionToken{allowMapIndex: true},
				obj: map[string]interface{}{
					"a": "one",
					"d": "four",
					"e": "five",
					"c": "three",
					"b": "two",
				},
				keys: []int64{0, 1, 3},
			},
			expected: expected{
				obj: []interface{}{
					"one",
					"two",
					"four",
				},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := test.input.token.getUnionByIndex(test.input.obj, test.input.keys)

			if test.expected.obj == nil {
				assert.Nil(t, obj)
			} else {
				assert.NotNil(t, obj)
				if array, ok := obj.([]interface{}); ok {
					assert.ElementsMatch(t, test.expected.obj, array)
				} else {
					assert.Equal(t, test.expected.obj, obj)
				}
			}

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}

}

func Test_UnionToken_getUnionByKey(t *testing.T) {

	type input struct {
		token *unionToken
		obj   interface{}
		keys  []string
	}

	type expected struct {
		obj []interface{}
		err string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				token: &unionToken{},
			},
			expected: expected{
				err: "union: invalid token target. expected [map] got [nil]",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   "string",
			},
			expected: expected{
				err: "union: invalid token target. expected [map] got [string]",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "one",
					"b": "two",
					"c": "three",
					"d": "four",
					"e": "five",
				},
				keys: []string{"a", "b", "c"},
			},
			expected: expected{
				obj: []interface{}{
					"one",
					"two",
					"three",
				},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "one",
					"b": "two",
					"c": "three",
					"d": "four",
					"e": "five",
				},
				keys: []string{"a", "b", "c", "f"},
			},
			expected: expected{
				err: "union: invalid token key 'f' not found",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "one",
					"b": "two",
					"c": "three",
					"d": "four",
					"e": "five",
				},
				keys: []string{"a", "b", "c", "f", "one", "blah"},
			},
			expected: expected{
				err: "union: invalid token key 'blah,f,one' not found",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: sampleStruct{
					One:   "one",
					Two:   "two",
					Three: 3,
					Four:  4,
					Five:  "five",
					Six:   "six",
				},
				keys: []string{"one", "three", "Six"},
			},
			expected: expected{
				obj: []interface{}{
					"one", int64(4), "six",
				},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   sampleStruct{},
				keys:  []string{"one", "three", "Six"},
			},
			expected: expected{
				obj: []interface{}{
					"", int64(0), "",
				},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   sampleStruct{},
				keys:  []string{"missing", "gone"},
			},
			expected: expected{
				err: "union: invalid token key 'gone,missing' not found",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := test.input.token.getUnionByKey(test.input.obj, test.input.keys)

			if test.expected.obj == nil {
				assert.Nil(t, obj)
			} else {
				assert.NotNil(t, obj)
				if obj != nil {
					assert.ElementsMatch(t, test.expected.obj, obj)
				}
			}

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}
}
