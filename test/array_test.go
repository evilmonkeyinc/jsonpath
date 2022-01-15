package test

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/jsonpath"
	"github.com/stretchr/testify/assert"
)

var consensusNone string = "none"

func Test_Array(t *testing.T) {

	type test struct {
		query         string
		data          string
		expected      interface{}
		consensus     interface{}
		expectedError string
	}

	tests := []test{
		{
			query:     "$[1:3]",
			data:      `["first", "second", "third", "forth", "fifth"]`,
			expected:  []interface{}{"second", "third"},
			consensus: []interface{}{"second", "third"},
		},
		{
			query:     "$[0:5]",
			data:      `["first", "second", "third", "forth", "fifth"]`,
			expected:  []interface{}{"first", "second", "third", "forth", "fifth"},
			consensus: []interface{}{"first", "second", "third", "forth", "fifth"},
		},
		{
			query:         "$[7:10]",
			data:          `["first", "second", "third"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     []interface{}{},
		},
		{
			query:         "$[1:3]",
			data:          `{":": 42, "more": "string", "a": 1, "b": 2, "c": 3, "1:3": "nice"}`,
			expected:      []interface{}{float64(42), float64(1)},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[1:10]",
			data:          `["first", "second", "third"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     []interface{}{"second", "third"},
		},
		{
			query:         "$[2:113667776004]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     []interface{}{"third", "forth", "fifth"},
		},
		{
			query:         "$[2:-113667776004:-1]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     []interface{}{},
		},
		{
			query:         "$[-113667776004:2]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     []interface{}{"first", "second"},
		},
		{
			query:         "$[113667776004:2:-1]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     []interface{}{},
		},
		{
			query:         "$[-4:-5]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[-4:-4]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[-4:-3]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{float64(4)},
			expectedError: "",
			consensus:     []interface{}{float64(4)},
		},
		{
			query:         "$[-4:1]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[-4:2]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[-4:3]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{float64(4)},
			expectedError: "",
			consensus:     []interface{}{float64(4)},
		},
		{
			query:         "$[3:0:-2]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[7:3:-1]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     []interface{}{},
		},
		{
			query:         "$[0:3:-2]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      []interface{}{"third", "first"},
			expectedError: "",
			consensus:     consensusNone,
		},
		{
			query:         "$[::-2]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "invalid JSONPath query '$[::-2]' invalid token. '[::-2]' does not match any token format",
			consensus:     consensusNone,
		},
		{
			query:         "$[1:]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      []interface{}{"second", "third", "forth", "fifth"},
			expectedError: "",
			consensus:     []interface{}{"second", "third", "forth", "fifth"},
		},
		{
			query:         "$[3::-1]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "invalid JSONPath query '$[3::-1]' invalid token. '[3::-1]' does not match any token format",
			consensus:     consensusNone,
		},
		{
			query:         "$[:2]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      []interface{}{"first", "second"},
			expectedError: "",
			consensus:     []interface{}{"first", "second"},
		},
		{
			query:         "$[:]",
			data:          `["first","second"]`,
			expected:      nil,
			expectedError: "invalid JSONPath query '$[:]' invalid token. '[:]' does not match any token format",
			consensus:     []interface{}{"first", "second"},
		},
		{
			query:         "$[:]",
			data:          `{":": 42, "more": "string"}`,
			expected:      nil,
			expectedError: "invalid JSONPath query '$[:]' invalid token. '[:]' does not match any token format",
			consensus:     []interface{}{},
		},
		{
			query:         "$[::]",
			data:          `["first","second"]`,
			expected:      nil,
			expectedError: "invalid JSONPath query '$[::]' invalid token. '[::]' does not match any token format",
			consensus:     []interface{}{"first", "second"},
		},
		{
			query:         "$[:2:-1]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "invalid JSONPath query '$[:2:-1]' invalid token. '[:2:-1]' does not match any token format",
			consensus:     consensusNone,
		},
		{
			query:         "$[3:-4]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[3:-3]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[3:-2]",
			data:          `[2, "a", 4, 5, 100, "nice"]`,
			expected:      []interface{}{float64(5)},
			expectedError: "",
			consensus:     []interface{}{float64(5)},
		},
		{
			query:         "$[2:1]",
			data:          `["first", "second", "third", "forth"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[0:0]",
			data:          `["first", "second"]`,
			expected:      []interface{}{},
			expectedError: "",
			consensus:     []interface{}{},
		},
		{
			query:         "$[0:1]",
			data:          `["first", "second"]`,
			expected:      []interface{}{"first"},
			expectedError: "",
			consensus:     []interface{}{"first"},
		},
		{
			query:         "$[-1:]",
			data:          `["first", "second", "third"]`,
			expected:      []interface{}{"third"},
			expectedError: "",
			consensus:     []interface{}{"third"},
		},
		{
			query:         "$[-2:]",
			data:          `["first", "second", "third"]`,
			expected:      []interface{}{"second", "third"},
			expectedError: "",
			consensus:     []interface{}{"second", "third"},
		},
		{
			query:         "$[-4:]",
			data:          `["first", "second", "third"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     []interface{}{"first", "second", "third"},
		},
		{
			query:         "$[0:3:2]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      []interface{}{"first", "third"},
			expectedError: "",
			consensus:     []interface{}{"first", "third"},
		},
		{
			query:         "$[0:3:0]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "range: invalid token out of range",
			consensus:     consensusNone,
		},
		{
			query:         "$[0:3:1]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      []interface{}{"first", "second", "third"},
			expectedError: "",
			consensus:     []interface{}{"first", "second", "third"},
		},
		{
			query:         "$[010:024:010]",
			data:          `[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]`,
			expected:      []interface{}{float64(10), float64(20)},
			expectedError: "",
			consensus:     []interface{}{float64(10), float64(20)},
		},
		{
			query:         "$[0:4:2]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      []interface{}{"first", "third"},
			expectedError: "",
			consensus:     []interface{}{"first", "third"},
		},
		{
			query:         "$[1:3:]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      []interface{}{"second", "third"},
			expectedError: "",
			consensus:     []interface{}{"second", "third"},
		},
		{
			query:         "$[::2]",
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			expectedError: "invalid JSONPath query '$[::2]' invalid token. '[::2]' does not match any token format",
			consensus:     []interface{}{"first", "third", "fifth"},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := jsonpath.QueryString(test.query, test.data)
			if test.expectedError == "" {
				assert.Nil(t, err, fmt.Sprintf("%s error should be nil", test.query))
			} else {
				assert.EqualError(t, err, test.expectedError, fmt.Sprintf("%s invalid error", test.query))
			}
			assert.EqualValues(t, test.expected, actual, fmt.Sprintf("%s unexpected value", test.query))
		})
	}

}
