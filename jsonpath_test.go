package jsonpath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleData string = `
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
		query    string
		expected expected
	}{
		{
			name:  "the authors of all books in the store",
			query: "$.store.book[*].author",
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
			name:  "all authors",
			query: "$..author",
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
			name:  "all things in store, which are some books and a red bicycle.",
			query: "$.store.*",
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
			name:  "the price of everything in the store.",
			query: "$.store..price",
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
			name:  "the third book",
			query: "$..book[2]",
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
			name:  "the last book in order.",
			query: "$..book[(@.length-1)]",
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
			name:  "the last book in order alt.",
			query: "$..book[-1:]",
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
			name:  "the first two books",
			query: "$..book[0,1]",
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
			name:  "the first two books alt",
			query: "$..book[:2]",
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
			name:  "filter all books with isbn number",
			query: "$..book[?(@.isbn)]",
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
			name:  "filter all books cheapier than 10",
			query: "$..book[?(@.price<10)]",
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
			name:  "filter all books that are not expensive",
			query: "$..book[?(@.price<$.expensive)]",
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
			name:  "All members of JSON structure.",
			query: "$..*",
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
			actual, actualErr := FindFromJSONString(test.query, sampleData)
			assert.ElementsMatch(t, test.expected.target, actual, fmt.Sprintf("'%s' invalid result", test.query))
			assert.Equal(t, test.expected.err, actualErr, fmt.Sprintf("'%s' invalid error", test.query))
		})
	}
}

func Test_JSONPath_Find(t *testing.T) {

	type input struct {
		query string
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
				query: "$.expensive",
			},
			expected: expected{
				obj: float64(10),
			},
		},
		{
			input: input{
				query: "$.store.book[0].title",
			},
			expected: expected{
				obj: "Sayings of the Century",
			},
		},
		{
			input: input{
				query: "$.store.book[*].author",
			},
			expected: expected{
				obj: []interface{}{
					"Nigel Rees",
					"Evelyn Waugh",
					"Herman Melville",
					"J. R. R. Tolkien",
				},
			},
		},
		{
			input: input{
				query: "$..author",
			},
			expected: expected{
				obj: []interface{}{
					"Nigel Rees",
					"Evelyn Waugh",
					"Herman Melville",
					"J. R. R. Tolkien",
				},
			},
		},
		{
			input: input{
				query: "$.store.*",
			},
			expected: expected{
				obj: []interface{}{
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
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			jsonPath := &JSONPath{}
			err := jsonPath.compile(test.input.query)

			assert.Nil(t, err)

			obj, err := jsonPath.FindFromJSONString(sampleData)
			if test.expected.err != "" {
				assert.EqualError(t, err, test.expected.err)
			} else {
				assert.Nil(t, err)
			}

			if expectArray, ok := test.expected.obj.([]interface{}); ok {
				assert.ElementsMatch(t, expectArray, obj)
			} else {
				assert.EqualValues(t, test.expected.obj, obj)
			}
		})
	}
}
