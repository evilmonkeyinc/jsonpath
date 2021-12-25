package jsonpath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getKey(t *testing.T) {

	type input struct {
		obj interface{}
		key string
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
				obj: nil,
				err: "cannot get key from nil map",
			},
		},
		{
			input: input{
				obj: "",
			},
			expected: expected{
				obj: nil,
				err: "invalid object. expected map",
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"key": true,
				},
				key: "key",
			},
			expected: expected{
				obj: true,
				err: "",
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"key": true,
				},
				key: "missing",
			},
			expected: expected{
				obj: nil,
				err: "'missing' key not found in object",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := getKey(test.input.obj, test.input.key)

			assert.Equal(t, test.expected.obj, obj)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}
}

func Test_getIndex(t *testing.T) {

	type input struct {
		obj interface{}
		idx int
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
				err: "cannot get index from nil slice",
			},
		},
		{
			input: input{
				obj: "not a slice",
			},
			expected: expected{
				err: "invalid object. expected slice",
			},
		},
		{
			input: input{
				obj: []string{"one", "two", "three"},
				idx: 0,
			},
			expected: expected{
				obj: "one",
			},
		},
		{
			input: input{
				obj: []string{"one", "two", "three"},
				idx: 2,
			},
			expected: expected{
				obj: "three",
			},
		},
		{
			input: input{
				obj: []string{"one", "two", "three"},
				idx: 4,
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: 1,
			},
			expected: expected{
				obj: 2,
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: -1,
			},
			expected: expected{
				obj: "three",
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: -2,
			},
			expected: expected{
				obj: 2,
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: -3,
			},
			expected: expected{
				obj: "one",
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: -4,
			},
			expected: expected{
				err: "index out of range",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := getIndex(test.input.obj, test.input.idx)

			assert.Equal(t, test.expected.obj, obj)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}

}

func Test_getRange(t *testing.T) {
	type input struct {
		obj              interface{}
		start, end, step *int
	}

	type expected struct {
		obj interface{}
		err string
	}

	testArray := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, "ten", "eleven", "twelve", 13}

	intPtr := func(i int) *int {
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
				err: "cannot get range from nil slice",
			},
		},
		{
			input: input{
				obj: "not a slice",
			},
			expected: expected{
				err: "invalid object. expected slice",
			},
		},
		{
			input: input{
				obj:   testArray,
				start: intPtr(15),
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			input: input{
				obj: testArray,
				end: intPtr(15),
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			input: input{
				obj:  testArray,
				step: intPtr(0),
			},
			expected: expected{
				err: "invalid parameter. step should be greater than 1",
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
				obj: testArray[0:13],
			},
		},
		{
			input: input{
				obj: testArray,
				end: intPtr(-3),
			},
			expected: expected{
				obj: testArray[0:11],
			},
		},
		{
			input: input{
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
				start: intPtr(1),
				step:  intPtr(2),
			},
			expected: expected{
				obj: []interface{}{"two", "four"},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := getRange(test.input.obj, test.input.start, test.input.end, test.input.step)

			assert.Equal(t, test.expected.obj, obj)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}

}
