package jsonpath

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/script"
	"github.com/stretchr/testify/assert"
)

type sampleData struct {
	Expensive float64    `json:"expensive"`
	Store     *storeData `json:"store"`
}

type storeData struct {
	Book []*bookData `json:"book"`
}

type bookData struct {
	Author   string  `json:"author"`
	Category string  `json:"category"`
	ISBN     string  `json:"isbn"`
	Price    float64 `json:"price"`
	Title    string  `json:"title"`
}

var sampleDataObject *sampleData = &sampleData{
	Expensive: 10,
	Store: &storeData{
		Book: []*bookData{
			{
				Category: "reference",
				Author:   "Nigel Rees",
				Title:    "Sayings of the Century",
				Price:    8.95,
			},
			{
				Category: "fiction",
				Author:   "Evelyn Waugh",
				Title:    "Sword of Honour",
				Price:    12.99,
			},
			{
				Category: "fiction",
				Author:   "Herman Melville",
				Title:    "Moby Dick",
				ISBN:     "0-553-21311-3",
				Price:    8.99,
			},
			{
				Category: "fiction",
				Author:   "J. R. R. Tolkien",
				Title:    "The Lord of the Rings",
				ISBN:     "0-395-19395-8",
				Price:    22.99,
			},
		},
	},
}

var sampleDataString string = `
{
	"store": {
		"book": [{
				"category": "reference",
				"author": "Nigel Rees",
				"title": "Sayings of the Century",
				"price": 8.95
			},
			{
				"category": "fiction",
				"author": "Evelyn Waugh",
				"title": "Sword of Honour",
				"price": 12.99
			},
			{
				"category": "fiction",
				"author": "Herman Melville",
				"title": "Moby Dick",
				"isbn": "0-553-21311-3",
				"price": 8.99
			},
			{
				"category": "fiction",
				"author": "J. R. R. Tolkien",
				"title": "The Lord of the Rings",
				"isbn": "0-395-19395-8",
				"price": 22.99
			}
		],
		"bicycle": {
			"color": "red",
			"price": 19.95
		}
	},
	"expensive": 10
}
`

// Tests designed after the examples in the specification document
//
// https://goessner.net/articles/JsonPath/
func Test_SpecificationTests(t *testing.T) {

	type expected struct {
		target []interface{}
		err    error
	}

	tests := []struct {
		name     string
		selector string
		expected expected
	}{
		{
			name:     "the authors of all books in the store",
			selector: "$.store.book[*].author",
			expected: expected{
				target: []interface{}{
					"Nigel Rees",
					"Evelyn Waugh",
					"Herman Melville",
					"J. R. R. Tolkien",
				},
			},
		},
		{
			name:     "all authors",
			selector: "$..author",
			expected: expected{
				target: []interface{}{
					"Nigel Rees",
					"Evelyn Waugh",
					"Herman Melville",
					"J. R. R. Tolkien",
				},
			},
		},
		{
			name:     "all things in store, which are some books and a red bicycle.",
			selector: "$.store.*",
			expected: expected{
				target: []interface{}{
					[]interface{}{
						map[string]interface{}{
							"category": "reference",
							"author":   "Nigel Rees",
							"title":    "Sayings of the Century",
							"price":    8.95,
						},
						map[string]interface{}{
							"category": "fiction",
							"author":   "Evelyn Waugh",
							"title":    "Sword of Honour",
							"price":    12.99,
						},
						map[string]interface{}{
							"category": "fiction",
							"author":   "Herman Melville",
							"title":    "Moby Dick",
							"isbn":     "0-553-21311-3",
							"price":    8.99,
						},
						map[string]interface{}{
							"category": "fiction",
							"author":   "J. R. R. Tolkien",
							"title":    "The Lord of the Rings",
							"isbn":     "0-395-19395-8",
							"price":    22.99,
						},
					},
					map[string]interface{}{
						"color": "red",
						"price": 19.95,
					},
				},
			},
		},
		{
			name:     "the price of everything in the store.",
			selector: "$.store..price",
			expected: expected{
				target: []interface{}{
					8.95,
					12.99,
					8.99,
					22.99,
					19.95,
				},
			},
		},
		{
			name:     "the third book",
			selector: "$..book[2]",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"category": "fiction",
						"author":   "Herman Melville",
						"title":    "Moby Dick",
						"isbn":     "0-553-21311-3",
						"price":    8.99,
					},
				},
			},
		},
		{
			name:     "the last book in order.",
			selector: "$..book[(@.length-1)]",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"category": "fiction",
						"author":   "J. R. R. Tolkien",
						"title":    "The Lord of the Rings",
						"isbn":     "0-395-19395-8",
						"price":    22.99,
					},
				},
			},
		},
		{
			name:     "the last book in order alt.",
			selector: "$..book[-1:]",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"category": "fiction",
						"author":   "J. R. R. Tolkien",
						"title":    "The Lord of the Rings",
						"isbn":     "0-395-19395-8",
						"price":    22.99,
					},
				},
			},
		},
		{
			name:     "the first two books",
			selector: "$..book[0,1]",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"category": "reference",
						"author":   "Nigel Rees",
						"title":    "Sayings of the Century",
						"price":    8.95,
					},
					map[string]interface{}{
						"category": "fiction",
						"author":   "Evelyn Waugh",
						"title":    "Sword of Honour",
						"price":    12.99,
					},
				},
			},
		},
		{
			name:     "the first two books alt",
			selector: "$..book[:2]",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"category": "reference",
						"author":   "Nigel Rees",
						"title":    "Sayings of the Century",
						"price":    8.95,
					},
					map[string]interface{}{
						"category": "fiction",
						"author":   "Evelyn Waugh",
						"title":    "Sword of Honour",
						"price":    12.99,
					},
				},
			},
		},
		{
			name:     "filter all books with isbn number",
			selector: "$..book[?(@.isbn)]",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"category": "fiction",
						"author":   "Herman Melville",
						"title":    "Moby Dick",
						"isbn":     "0-553-21311-3",
						"price":    8.99,
					},
					map[string]interface{}{
						"category": "fiction",
						"author":   "J. R. R. Tolkien",
						"title":    "The Lord of the Rings",
						"isbn":     "0-395-19395-8",
						"price":    22.99,
					},
				},
			},
		},
		{
			name:     "filter all books cheapier than 10",
			selector: "$..book[?(@.price<10)]",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"category": "reference",
						"author":   "Nigel Rees",
						"title":    "Sayings of the Century",
						"price":    8.95,
					},
					map[string]interface{}{
						"category": "fiction",
						"author":   "Herman Melville",
						"title":    "Moby Dick",
						"isbn":     "0-553-21311-3",
						"price":    8.99,
					},
				},
			},
		},
		{
			name:     "filter all books that are not expensive",
			selector: "$..book[?(@.price<$.expensive)]",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"category": "reference",
						"author":   "Nigel Rees",
						"title":    "Sayings of the Century",
						"price":    8.95,
					},
					map[string]interface{}{
						"category": "fiction",
						"author":   "Herman Melville",
						"title":    "Moby Dick",
						"isbn":     "0-553-21311-3",
						"price":    8.99,
					},
				},
			},
		},
		{
			name:     "All members of JSON structure.",
			selector: "$..*",
			expected: expected{
				target: []interface{}{
					map[string]interface{}{
						"book": []interface{}{
							map[string]interface{}{
								"category": "reference",
								"author":   "Nigel Rees",
								"title":    "Sayings of the Century",
								"price":    8.95,
							},
							map[string]interface{}{
								"category": "fiction",
								"author":   "Evelyn Waugh",
								"title":    "Sword of Honour",
								"price":    12.99,
							},
							map[string]interface{}{
								"category": "fiction",
								"author":   "Herman Melville",
								"title":    "Moby Dick",
								"isbn":     "0-553-21311-3",
								"price":    8.99,
							},
							map[string]interface{}{
								"category": "fiction",
								"author":   "J. R. R. Tolkien",
								"title":    "The Lord of the Rings",
								"isbn":     "0-395-19395-8",
								"price":    22.99,
							},
						},
						"bicycle": map[string]interface{}{
							"color": "red",
							"price": 19.95,
						},
					},
					float64(10),
					[]interface{}{
						map[string]interface{}{
							"category": "reference",
							"author":   "Nigel Rees",
							"title":    "Sayings of the Century",
							"price":    8.95,
						},
						map[string]interface{}{
							"category": "fiction",
							"author":   "Evelyn Waugh",
							"title":    "Sword of Honour",
							"price":    12.99,
						},
						map[string]interface{}{
							"category": "fiction",
							"author":   "Herman Melville",
							"title":    "Moby Dick",
							"isbn":     "0-553-21311-3",
							"price":    8.99,
						},
						map[string]interface{}{
							"category": "fiction",
							"author":   "J. R. R. Tolkien",
							"title":    "The Lord of the Rings",
							"isbn":     "0-395-19395-8",
							"price":    22.99,
						},
					},
					map[string]interface{}{
						"color": "red",
						"price": 19.95,
					},
					map[string]interface{}{
						"category": "reference",
						"author":   "Nigel Rees",
						"title":    "Sayings of the Century",
						"price":    8.95,
					},
					map[string]interface{}{
						"category": "fiction",
						"author":   "Evelyn Waugh",
						"title":    "Sword of Honour",
						"price":    12.99,
					},
					map[string]interface{}{
						"category": "fiction",
						"author":   "Herman Melville",
						"title":    "Moby Dick",
						"isbn":     "0-553-21311-3",
						"price":    8.99,
					},
					map[string]interface{}{
						"category": "fiction",
						"author":   "J. R. R. Tolkien",
						"title":    "The Lord of the Rings",
						"isbn":     "0-395-19395-8",
						"price":    22.99,
					},
					"reference",
					"Nigel Rees",
					"Sayings of the Century",
					8.95,
					"fiction",
					"Evelyn Waugh",
					"Sword of Honour",
					12.99,
					"fiction",
					"Herman Melville",
					"Moby Dick",
					"0-553-21311-3",
					8.99,
					"fiction",
					"J. R. R. Tolkien",
					"The Lord of the Rings",
					"0-395-19395-8",
					22.99,
					"red",
					19.95,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, actualErr := QueryString(test.selector, sampleDataString)
			assert.ElementsMatch(t, test.expected.target, actual, fmt.Sprintf("'%s' invalid result", test.selector))
			assert.Equal(t, test.expected.err, actualErr, fmt.Sprintf("'%s' invalid error", test.selector))
		})
	}
}

func Test_Compile(t *testing.T) {

	type input struct {
		selector string
	}

	type expected struct {
		err    string
		tokens int
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				selector: "",
			},
			expected: expected{
				err: "invalid JSONPath selector '' unexpected token '' at index 0",
			},
		},
		{
			input: input{
				selector: "@.[1, 2]",
			},
			expected: expected{
				tokens: 2,
			},
		},
		{
			input: input{
				selector: "@.length<1",
			},
			expected: expected{
				tokens: 2,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			selector, err := Compile(test.input.selector)
			if test.expected.err != "" {
				assert.Nil(t, selector)
				assert.EqualError(t, err, test.expected.err)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, selector)

			assert.Equal(t, test.input.selector, selector.selector)
			assert.Len(t, selector.tokens, test.expected.tokens)
		})
	}

}

func Test_QueryString(t *testing.T) {

	type input struct {
		selector string
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
				selector: "$.expensive",
			},
			expected: expected{
				err: "invalid data. unexpected type or nil",
			},
		},
		{
			input: input{
				selector: "invalid",
				jsonData: "{}",
			},
			expected: expected{
				err: "invalid JSONPath selector 'invalid' unexpected token 'i' at index 0",
			},
		},
		{
			input: input{

				selector: "$.expensive",
				jsonData: "{}",
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				selector: "$.expensive",
				jsonData: `{"expensive": "test"}`,
			},
			expected: expected{
				value: "test",
			},
		},
		{
			input: input{
				selector: "$.expensive",
				jsonData: sampleDataString,
			},
			expected: expected{
				value: int64(10),
			},
		},
		{
			input: input{
				selector: "$..author",
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
				selector: "$.store.book.length",
				jsonData: sampleDataString,
			},
			expected: expected{
				value: int64(4),
			},
		},
		{
			input: input{
				selector: "$..book.length",
				jsonData: sampleDataString,
			},
			expected: expected{
				value: []interface{}{
					int64(4),
				},
			},
		},
		{
			input: input{
				selector: "$.length",
				jsonData: `[1,2,3]`,
			},
			expected: expected{
				value: int64(3),
			},
		},
		{
			input: input{
				selector: "$.length",
				jsonData: `[1,2,]`,
			},
			expected: expected{
				err: "invalid data. invalid character ']' looking for beginning of value",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			value, err := QueryString(test.input.selector, test.input.jsonData)

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

func Test_Query(t *testing.T) {

	type input struct {
		selector string
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
				selector: "$.expensive",
			},
			expected: expected{
				err: "key: invalid token target. expected [map] got [nil]",
			},
		},
		{
			input: input{
				selector: "invalid",
				jsonData: &sampleData{},
			},
			expected: expected{
				err: "invalid JSONPath selector 'invalid' unexpected token 'i' at index 0",
			},
		},
		{
			input: input{

				selector: "$.expensive",
				jsonData: &storeData{},
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				selector: "$.expensive",
				jsonData: &sampleData{Expensive: 15},
			},
			expected: expected{
				value: float64(15),
			},
		},
		{
			input: input{
				selector: "$.expensive",
				jsonData: sampleDataObject,
			},
			expected: expected{
				value: float64(10),
			},
		},
		{
			input: input{
				selector: "$..author",
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
		{
			input: input{
				selector: "$.store.book.length",
				jsonData: sampleDataObject,
			},
			expected: expected{
				value: int64(4),
			},
		},
		{
			input: input{
				selector: "$..book.length",
				jsonData: sampleDataObject,
			},
			expected: expected{
				value: []interface{}{
					int64(4),
				},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			value, err := Query(test.input.selector, test.input.jsonData)

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

func Test_Selector_compile(t *testing.T) {

	type input struct {
		selector string
		engine   script.Engine
	}

	type expected struct {
		err    string
		tokens int
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				selector: "",
			},
			expected: expected{
				err: "unexpected token '' at index 0",
			},
		},
		{
			input: input{
				selector: "@.[1,(]",
			},
			expected: expected{
				err: "invalid token. '[1,(]' does not match any token format",
			},
		},
		{
			input: input{
				selector: "@.length<1",
			},
			expected: expected{
				tokens: 2,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			selector := &Selector{}
			actual := selector.compile(test.input.selector, test.input.engine)
			if test.expected.err == "" {
				assert.Nil(t, actual)
			} else {
				assert.EqualError(t, actual, test.expected.err)
			}

			assert.Len(t, selector.tokens, test.expected.tokens)
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
