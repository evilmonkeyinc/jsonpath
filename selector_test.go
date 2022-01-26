package jsonpath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Selector_String(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "$.store.book[*].author",
			expected: "$['store']['book'][*]['author']",
		},
		{
			input:    "$..author",
			expected: "$..['author']",
		},
		{
			input:    "$.store.*",
			expected: "$['store'][*]",
		},
		{
			input:    "$.store..price",
			expected: "$['store']..['price']",
		},
		{
			input:    "$..book[2]",
			expected: "$..['book'][2]",
		},
		{
			input:    "$..book[(@.length-1)]",
			expected: "$..['book'][(@.length-1)]",
		},
		{
			input:    "$..book[-1:]",
			expected: "$..['book'][-1:]",
		},
		{
			input:    "$..book[0,1]",
			expected: "$..['book'][0,1]",
		},
		{
			input:    "$..book[:2]",
			expected: "$..['book'][:2]",
		},
		{
			input:    "$..book[?(@.isbn)]",
			expected: "$..['book'][?(@.isbn)]",
		},
		{
			input:    "$..book[?(@.price<10)]",
			expected: "$..['book'][?(@.price<10)]",
		},
		{
			input:    "$..*",
			expected: "$..[*]",
		},
		{
			input:    "$.store. book[0].author",
			expected: "$['store']['book'][0]['author']",
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			compiled, _ := Compile(test.input)
			assert.Equal(t, test.expected, compiled.String())
		})
	}

}

func Test_Selector_QueryString(t *testing.T) {

	sampleQuery, _ := Compile("$.expensive")
	altSampleQuery, _ := Compile("$..author")
	lengthQuery, _ := Compile("$.length")
	rootQuery, _ := Compile("$")

	type input struct {
		selector *Selector
		jsonData string
	}

	type expected struct {
		value interface{}
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				selector: sampleQuery,
			},
			expected: expected{
				err: "invalid data. unexpected type or nil",
			},
		},
		{
			input: input{
				selector: rootQuery,
				jsonData: "42",
			},
			expected: expected{
				value: int64(42),
			},
		},
		{
			input: input{
				selector: rootQuery,
				jsonData: "3.14",
			},
			expected: expected{
				value: float64(3.14),
			},
		},
		{
			input: input{
				selector: rootQuery,
				jsonData: "true",
			},
			expected: expected{
				value: true,
			},
		},
		{
			input: input{
				selector: rootQuery,
				jsonData: "false",
			},
			expected: expected{
				value: false,
			},
		},
		{
			input: input{
				selector: rootQuery,
				jsonData: "not a json string",
			},
			expected: expected{
				err: "invalid data. unexpected type or nil",
			},
		},
		{
			input: input{
				selector: rootQuery,
				jsonData: `"json string"`,
			},
			expected: expected{
				value: "json string",
			},
		},
		{
			input: input{
				selector: &Selector{
					selector: "invalid",
				},
				jsonData: "{}",
			},
			expected: expected{
				err: "invalid JSONPath selector 'invalid'",
			},
		},
		{
			input: input{
				selector: sampleQuery,
				jsonData: `{"key"}`,
			},
			expected: expected{
				err: "invalid data. invalid character '}' after object key",
			},
		},
		{
			input: input{
				selector: sampleQuery,
				jsonData: "{}",
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				selector: sampleQuery,
				jsonData: `{"expensive": "test"}`,
			},
			expected: expected{
				value: "test",
			},
		},
		{
			input: input{
				selector: sampleQuery,
				jsonData: sampleDataString,
			},
			expected: expected{
				value: int64(10),
			},
		},
		{
			input: input{
				selector: altSampleQuery,
				jsonData: sampleDataString,
			},
			expected: expected{
				value: []interface{}{
					"Nigel Rees",
					"Evelyn Waugh",
					"Herman Melville",
					"J. R. R. Tolkien",
				},
			},
		},
		{
			input: input{
				selector: lengthQuery,
				jsonData: `[1,2,3]`,
			},
			expected: expected{
				value: int64(3),
			},
		},
		{
			input: input{
				selector: lengthQuery,
				jsonData: `[1,2,]`,
			},
			expected: expected{
				err: "invalid data. invalid character ']' looking for beginning of value",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			value, err := test.input.selector.QueryString(test.input.jsonData)

			if test.expected.err != "" {
				assert.EqualError(t, err, test.expected.err)
			} else {
				assert.Nil(t, err)
			}

			if expectArray, ok := test.expected.value.([]interface{}); ok {
				assert.ElementsMatch(t, expectArray, value)
			} else {
				assert.EqualValues(t, test.expected.value, value)
			}
		})
	}
}

func Test_Selector_Query(t *testing.T) {

	sampleSelector, _ := Compile("$.expensive")
	altSampleSelector, _ := Compile("$..author")

	type input struct {
		selector *Selector
		jsonData interface{}
	}

	type expected struct {
		value interface{}
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				selector: sampleSelector,
				jsonData: make(chan bool, 1),
			},
			expected: expected{
				err: "key: invalid token target. expected [map] got [chan]",
			},
		},
		{
			input: input{
				selector: sampleSelector,
				jsonData: "not something that can be marshaled",
			},
			expected: expected{
				err: "key: invalid token target. expected [map] got [string]",
			},
		},
		{
			input: input{
				selector: sampleSelector,
				jsonData: &storeData{},
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				selector: &Selector{
					selector: "invalid",
				},
				jsonData: &sampleData{},
			},
			expected: expected{
				err: "invalid JSONPath selector 'invalid'",
			},
		},
		{
			input: input{
				selector: sampleSelector,
				jsonData: &bookData{},
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				selector: sampleSelector,
				jsonData: &sampleData{
					Expensive: 15,
				},
			},
			expected: expected{
				value: float64(15),
			},
		},
		{
			input: input{
				selector: sampleSelector,
				jsonData: sampleDataObject,
			},
			expected: expected{
				value: float32(10),
			},
		},
		{
			input: input{
				selector: altSampleSelector,
				jsonData: sampleDataObject,
			},
			expected: expected{
				value: []interface{}{
					"Nigel Rees",
					"Evelyn Waugh",
					"Herman Melville",
					"J. R. R. Tolkien",
				},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			value, err := test.input.selector.Query(test.input.jsonData)

			if test.expected.err != "" {
				assert.EqualError(t, err, test.expected.err)
			} else {
				assert.Nil(t, err)
			}

			if expectArray, ok := test.expected.value.([]interface{}); ok {
				assert.NotNil(t, value)
				if value != nil {
					assert.ElementsMatch(t, expectArray, value)
				}
			} else {
				assert.EqualValues(t, test.expected.value, value)
			}
		})
	}
}
