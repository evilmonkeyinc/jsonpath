package token

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/stretchr/testify/assert"
)

// Test unionToken struct conforms to Token interface
var _ Token = &unionToken{}

func Test_newUnionToken(t *testing.T) {
	assert.IsType(t, &unionToken{}, newUnionToken(nil, nil))

	type input struct {
		args    []interface{}
		options *option.QueryOptions
	}

	type expected *unionToken

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				options: nil,
			},
			expected: &unionToken{
				allowMap:    false,
				allowString: false,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{},
			},
			expected: &unionToken{
				allowMap:    false,
				allowString: false,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					FailUnionOnInvalidIdentifier: true,
				},
			},
			expected: &unionToken{
				allowMap:                     false,
				allowString:                  false,
				failUnionOnInvalidIdentifier: true,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					AllowMapReferenceByIndex:    false,
					AllowStringReferenceByIndex: false,

					AllowMapReferenceByIndexInUnion:    true,
					AllowStringReferenceByIndexInUnion: true,
				},
			},
			expected: &unionToken{
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					AllowMapReferenceByIndex:    true,
					AllowStringReferenceByIndex: true,

					AllowMapReferenceByIndexInUnion:    false,
					AllowStringReferenceByIndexInUnion: false,
				},
			},
			expected: &unionToken{
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					AllowMapReferenceByIndex:    true,
					AllowStringReferenceByIndex: true,

					AllowMapReferenceByIndexInUnion:    true,
					AllowStringReferenceByIndexInUnion: true,
				},
			},
			expected: &unionToken{
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				options: &option.QueryOptions{
					AllowMapReferenceByIndex:    false,
					AllowStringReferenceByIndex: false,

					AllowMapReferenceByIndexInUnion:    false,
					AllowStringReferenceByIndexInUnion: true,
				},
			},
			expected: &unionToken{
				allowMap:    false,
				allowString: true,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := newUnionToken(test.input.args, test.input.options)
			assert.EqualValues(t, test.expected, actual)
		})
	}
}

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
			input: &unionToken{arguments: []interface{}{
				1,
				&expressionToken{expression: "4%2", compiledExpression: &testCompiledExpression{response: 0}},
				"last",
			}},
			expected: "[1,(4%2),'last']",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_UnionToken_Type(t *testing.T) {
	assert.Equal(t, "union", (&unionToken{}).Type())
}

var unionTests = []*tokenTest{
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
				&expressionToken{expression: "nil", compiledExpression: &testCompiledExpression{response: nil}},
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
			current: []interface{}{1, 2, 3, 4, 5},
		},
		expected: expected{
			err: "union: invalid token argument. expected [int string] got [float64]",
		},
	},
	{
		token: &unionToken{
			arguments: []interface{}{
				&expressionToken{expression: "", compiledExpression: &testCompiledExpression{response: ""}},
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
				&expressionToken{expression: "1+1", compiledExpression: &testCompiledExpression{response: 2}},
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
				&expressionToken{expression: "1+1", compiledExpression: &testCompiledExpression{response: 2}},
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
			value: []interface{}{"one", "four"},
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
			allowMap: true,
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
			allowMap: true,
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
			value: []interface{}{"one", "four"},
		},
	},
	{
		token: &unionToken{
			arguments:   []interface{}{0, 2, 4},
			allowString: true,
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
			arguments:   []interface{}{0, 2, 4},
			allowString: true,
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

func Test_UnionToken_Apply(t *testing.T) {
	batchTokenTests(t, unionTests)
}

func Benchmark_UnionToken_Apply(b *testing.B) {
	batchTokenBenchmarks(b, unionTests)
}

func Test_UnionToken_getUnionByIndex(t *testing.T) {

	type input struct {
		token *unionToken
		obj   interface{}
		keys  []int64
		next  []Token
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
					allowMap: true,
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
					allowString: true,
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
					allowMap:    true,
					allowString: true,
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
				token: &unionToken{
					failUnionOnInvalidIdentifier: true,
				},
				obj:  []string{"one", "two", "three"},
				keys: []int64{4},
			},
			expected: expected{
				err: "union: invalid token out of range",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{4},
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
		{
			input: input{
				token: &unionToken{
					failUnionOnInvalidIdentifier: true,
				},
				obj:  []string{"one", "two", "three"},
				keys: []int64{-10},
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
				obj: []interface{}{},
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
				token: &unionToken{allowString: true},
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
				token: &unionToken{allowMap: true},
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
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{1, 1},
			},
			expected: expected{
				obj: []interface{}{"two", "two"},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{1, 1},
				next:  []Token{&indexToken{index: 1}},
			},
			expected: expected{
				obj: "two",
			},
		},
		{
			input: input{
				token: &unionToken{allowString: true},
				obj:   "abcdefghijklmnopqrstuvwxyz",
				keys:  []int64{0, 2, 4},
				next:  []Token{&indexToken{index: 2, allowString: true}},
			},
			expected: expected{
				obj: "e",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{1, 2},
				next:  []Token{&testToken{err: fmt.Errorf("fail")}},
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{1, 2},
				next:  []Token{&testToken{value: nil}},
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   []string{"one", "two", "three"},
				keys:  []int64{1, 2},
				next:  []Token{&testToken{value: "1"}},
			},
			expected: expected{
				obj: []interface{}{"1", "1"},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := test.input.token.getUnionByIndex(nil, test.input.obj, test.input.keys, test.input.next)

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
		next  []Token
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
				token: &unionToken{
					failUnionOnInvalidIdentifier: true,
				},
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
				keys: []string{"a", "b", "c", "f"},
			},
			expected: expected{
				obj: []interface{}{"one", "two", "three"},
			},
		},
		{
			input: input{
				token: &unionToken{
					failUnionOnInvalidIdentifier: true,
				},
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
				obj: []interface{}{"one", "two", "three"},
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
				token: &unionToken{
					failUnionOnInvalidIdentifier: true,
				},
				obj:  sampleStruct{},
				keys: []string{"missing", "gone"},
			},
			expected: expected{
				err: "union: invalid token key 'gone,missing' not found",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj:   sampleStruct{},
				keys:  []string{"missing", "gone"},
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: sampleStruct{
					One: "value",
				},
				keys: []string{"one", "one"},
			},
			expected: expected{
				obj: []interface{}{"value", "value"},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "value",
				},
				keys: []string{"a", "a", "a"},
			},
			expected: expected{
				obj: []interface{}{"value", "value", "value"},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "value",
				},
				keys: []string{"a", "a", "a"},
				next: []Token{&indexToken{index: 1}},
			},
			expected: expected{
				obj: "value",
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "value",
				},
				keys: []string{"a", "a", "a"},
				next: []Token{&testToken{err: fmt.Errorf("fail")}},
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "value",
				},
				keys: []string{"a", "a", "a"},
				next: []Token{&testToken{value: nil}},
			},
			expected: expected{
				obj: []interface{}{},
			},
		},
		{
			input: input{
				token: &unionToken{},
				obj: map[string]interface{}{
					"a": "value",
				},
				keys: []string{"a", "a", "a"},
				next: []Token{&testToken{value: "1"}},
			},
			expected: expected{
				obj: []interface{}{"1", "1", "1"},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := test.input.token.getUnionByKey(nil, test.input.obj, test.input.keys, test.input.next)

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
