package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnionToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &unionToken{},
			input: input{
				current: []interface{}{},
			},
			expected: expected{
				err: "invalid parameter. expected arguments",
			},
		},
	}

	batchTokenTests(t, tests)
}

func Test_getUnionByIndex(t *testing.T) {

	type input struct {
		obj  interface{}
		keys []int64
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
				obj: nil,
			},
			expected: expected{
				err: "cannot get union from nil object",
			},
		},
		{
			input: input{
				obj: 123,
			},
			expected: expected{
				err: "invalid object. expected array, map, or string",
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []int64{4},
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []int64{-10},
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []int64{-1, -2},
			},
			expected: expected{
				obj: []interface{}{"three", "two"},
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []int64{0, 2},
			},
			expected: expected{
				obj: []interface{}{"three", "one"},
			},
		},
		{
			input: input{
				obj:  []interface{}{"one", "two", 3},
				keys: []int64{0, 2},
			},
			expected: expected{
				obj: []interface{}{"one", 3},
			},
		},
		{
			input: input{
				obj:  [3]int64{1, 2, 3},
				keys: []int64{0, 2},
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
				obj:  "abcdefghijklmnopqrstuvwxyz",
				keys: []int64{0, 2, 4},
			},
			expected: expected{
				obj: "ace",
			},
		},
		{
			input: input{
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
			obj, err := getUnionByIndex(test.input.obj, test.input.keys)

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

func Test_getUnionByKey(t *testing.T) {

	type input struct {
		obj  interface{}
		keys []string
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
			input: input{},
			expected: expected{
				err: "cannot get union from nil object",
			},
		},
		{
			input: input{
				obj: "string",
			},
			expected: expected{
				err: "invalid object. expected map",
			},
		},
		{
			input: input{
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
				err: "'f' key not found in object",
			},
		},
		{
			input: input{
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
				err: "'blah,f,one' key not found in object",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := getUnionByKey(test.input.obj, test.input.keys)

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
