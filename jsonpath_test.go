package jsonpath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Benchmark_Selector(b *testing.B) {

	selectors := []string{
		"$.store.book[*].author",
		"$..author",
		"$.store.*",
		"$.store..price",
		"$..book[2]",
		"$..book[(@.length-1)]",
		"$..book[-1:]",
		"$..book[0,1]",
		"$..book[:2]",
		"$..book[?(@.isbn)]",
		"$..book[?(@.price<10)]",
		"$..book[?(@.price<$.expensive)]",
		"$..*",
	}

	for _, selector := range selectors {
		b.Run(fmt.Sprintf("%s", selector), func(b *testing.B) {
			var err error
			for i := 0; i < b.N; i++ {
				_, err = QueryString(selector, sampleDataString)
				if err != nil {
					b.Error()
				}
			}
		})
	}
}

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
		options  []Option
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
				options: []Option{
					OptionFunction(func(selector *Selector) error {
						assert.Equal(t, "@.[1, 2]", selector.selector)
						return nil
					}),
				},
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
		{
			input: input{
				selector: "@.[1,(]",
			},
			expected: expected{
				err: "invalid JSONPath selector '@.[1,(]' invalid token. '[1,(]' does not match any token format",
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
		{
			input: input{
				selector: "this wont matter",
				options: []Option{
					OptionFunction(func(selector *Selector) error {
						return fmt.Errorf("option error")
					}),
				},
			},
			expected: expected{
				err: "option error",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			selector, err := Compile(test.input.selector, test.input.options...)
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
