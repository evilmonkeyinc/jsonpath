package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getUnion(t *testing.T) {

	type input struct {
		obj  interface{}
		keys []interface{}
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
				obj: "not a array or map",
			},
			expected: expected{
				err: "invalid object. expected array or map",
			},
		},
		{
			input: input{
				obj:  map[string]interface{}{},
				keys: []interface{}{"one", 2},
			},
			expected: expected{
				err: "invalid parameter. expected string keys",
			},
		},
		{
			input: input{
				obj:  map[string]interface{}{},
				keys: []interface{}{"one", "two"},
			},
			expected: expected{
				err: "'one,two' key not found in object",
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"one":   1,
					"two":   2,
					"three": 3,
					"four":  4,
				},
				keys: []interface{}{"one"},
			},
			expected: expected{
				obj: []interface{}{1},
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"one":   1,
					"two":   2,
					"three": 3,
					"four":  4,
				},
				keys: []interface{}{"one", "two"},
			},
			expected: expected{
				obj: []interface{}{1, 2},
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"one":   1,
					"two":   2,
					"three": 3,
					"four":  4,
				},
				keys: []interface{}{"one", "three"},
			},
			expected: expected{
				obj: []interface{}{1, 3},
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []interface{}{"one"},
			},
			expected: expected{
				err: "invalid parameter. expected integer keys",
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []interface{}{4},
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []interface{}{-10},
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []interface{}{-1, -2},
			},
			expected: expected{
				obj: []interface{}{"three", "two"},
			},
		},
		{
			input: input{
				obj:  []string{"one", "two", "three"},
				keys: []interface{}{0, 2},
			},
			expected: expected{
				obj: []interface{}{"three", "one"},
			},
		},
		{
			input: input{
				obj:  []interface{}{"one", "two", 3},
				keys: []interface{}{0, 2},
			},
			expected: expected{
				obj: []interface{}{"one", 3},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			obj, err := getUnion(test.input.obj, test.input.keys)

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
