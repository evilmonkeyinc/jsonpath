package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test indexToken struct conforms to Token interface
var _ Token = &indexToken{}

func Test_newIndexToken(t *testing.T) {
	assert.IsType(t, &indexToken{}, newIndexToken(0, nil))

	type input struct {
		index   int64
		options *Options
	}

	type expected *indexToken

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				index:   0,
				options: nil,
			},
			expected: &indexToken{
				index:       0,
				allowMap:    false,
				allowString: false,
			},
		},
		{
			input: input{
				index:   0,
				options: &Options{},
			},
			expected: &indexToken{
				index:       0,
				allowMap:    false,
				allowString: false,
			},
		},
		{
			input: input{
				index: 0,
				options: &Options{
					AllowMapReferenceByIndex:    false,
					AllowStringReferenceByIndex: false,

					AllowMapReferenceByIndexInSubscript:    true,
					AllowStringReferenceByIndexInSubscript: true,
				},
			},
			expected: &indexToken{
				index:       0,
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				index: 0,
				options: &Options{
					AllowMapReferenceByIndex:    true,
					AllowStringReferenceByIndex: true,

					AllowMapReferenceByIndexInSubscript:    false,
					AllowStringReferenceByIndexInSubscript: false,
				},
			},
			expected: &indexToken{
				index:       0,
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				index: 0,
				options: &Options{
					AllowMapReferenceByIndex:    true,
					AllowStringReferenceByIndex: true,

					AllowMapReferenceByIndexInSubscript:    true,
					AllowStringReferenceByIndexInSubscript: true,
				},
			},
			expected: &indexToken{
				index:       0,
				allowMap:    true,
				allowString: true,
			},
		},
		{
			input: input{
				index: 0,
				options: &Options{
					AllowMapReferenceByIndex:    false,
					AllowStringReferenceByIndex: false,

					AllowMapReferenceByIndexInSubscript:    false,
					AllowStringReferenceByIndexInSubscript: true,
				},
			},
			expected: &indexToken{
				index:       0,
				allowMap:    false,
				allowString: true,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := newIndexToken(test.input.index, test.input.options)
			assert.EqualValues(t, test.expected, actual)
		})
	}
}

func Test_IndexToken_String(t *testing.T) {
	tests := []*tokenStringTest{
		{
			input:    &indexToken{},
			expected: "[0]",
		},
		{
			input:    &indexToken{index: -1},
			expected: "[-1]",
		},
		{
			input:    &indexToken{index: 10},
			expected: "[10]",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_IndexToken_Type(t *testing.T) {
	assert.Equal(t, "index", (&indexToken{}).Type())
}

func Test_IndexToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &indexToken{index: 0},
			input: input{
				current: nil,
			},
			expected: expected{
				err: "index: invalid token target. expected [array slice] got [nil]",
			},
		},
		{
			token: &indexToken{index: 0},
			input: input{
				current: 123,
			},
			expected: expected{
				err: "index: invalid token target. expected [array slice] got [int]",
			},
		},
		{
			token: &indexToken{index: 0, allowMap: true},
			input: input{
				current: 123,
			},
			expected: expected{
				err: "index: invalid token target. expected [array slice map] got [int]",
			},
		},
		{
			token: &indexToken{index: 0, allowString: true},
			input: input{
				current: 123,
			},
			expected: expected{
				err: "index: invalid token target. expected [array slice string] got [int]",
			},
		},
		{
			token: &indexToken{index: 0, allowMap: true, allowString: true},
			input: input{
				current: 123,
			},
			expected: expected{
				err: "index: invalid token target. expected [array slice map string] got [int]",
			},
		},
		{
			token: &indexToken{index: 5},
			input: input{
				current: "Find(X)",
			},
			expected: expected{
				err: "index: invalid token target. expected [array slice] got [string]",
			},
		},
		{
			token: &indexToken{index: 5, allowString: true},
			input: input{
				current: "Find(X)",
			},
			expected: expected{
				value: "X",
			},
		},
		{
			token: &indexToken{index: 0},
			input: input{
				current: [3]string{"one", "two", "three"},
			},
			expected: expected{
				value: "one",
			},
		},
		{
			token: &indexToken{index: 0},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				value: "one",
			},
		},
		{
			token: &indexToken{index: 2},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				value: "three",
			},
		},
		{
			token: &indexToken{index: 4},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				err: "index: invalid token out of range",
			},
		},
		{
			token: &indexToken{index: 1},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				value: 2,
			},
		},
		{
			token: &indexToken{index: -1},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				value: "three",
			},
		},
		{
			token: &indexToken{index: -2},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				value: 2,
			},
		},
		{
			token: &indexToken{index: -3},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				value: "one",
			},
		},
		{
			token: &indexToken{index: -4},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				err: "index: invalid token out of range",
			},
		},
		{
			token: &indexToken{index: 1},
			input: input{
				current: []interface{}{
					map[string]interface{}{
						"name":  "one",
						"value": 1,
					},
					map[string]interface{}{
						"name":  "two",
						"value": 2,
					},
					map[string]interface{}{
						"name":  "three",
						"value": 3,
					},
				},
				tokens: []Token{
					&keyToken{key: "name"},
				},
			},
			expected: expected{
				value: "two",
			},
		},
		{
			token: &indexToken{index: 1},
			input: input{
				current: map[string]interface{}{
					"a": map[string]interface{}{
						"name":  "one",
						"value": 1,
					},
					"c": map[string]interface{}{
						"name":  "three",
						"value": 3,
					},
					"b": map[string]interface{}{
						"name":  "two",
						"value": 2,
					},
				},
			},
			expected: expected{
				err: "index: invalid token target. expected [array slice] got [map]",
			},
		},
		{
			token: &indexToken{index: 1, allowMap: true},
			input: input{
				current: map[string]interface{}{
					"a": map[string]interface{}{
						"name":  "one",
						"value": 1,
					},
					"c": map[string]interface{}{
						"name":  "three",
						"value": 3,
					},
					"b": map[string]interface{}{
						"name":  "two",
						"value": 2,
					},
				},
			},
			expected: expected{
				value: map[string]interface{}{
					"name":  "two",
					"value": 2,
				},
			},
		},
	}

	batchTokenTests(t, tests)
}
