package jsonpath

import (
	"encoding/json"
	"fmt"
	"testing"

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
			actual, actualErr := QueryString(test.query, sampleDataString)
			assert.ElementsMatch(t, test.expected.target, actual, fmt.Sprintf("'%s' invalid result", test.query))
			assert.Equal(t, test.expected.err, actualErr, fmt.Sprintf("'%s' invalid error", test.query))
		})
	}
}

func Test_Compile(t *testing.T) {

	type input struct {
		queryPath string
		isStrict  bool
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
				queryPath: "",
				isStrict:  false,
			},
			expected: expected{
				err: "invalid JSONPath query '' unexpected token '' at index 0",
			},
		},
		{
			input: input{
				queryPath: "@.[1, 2]",
				isStrict:  true,
			},
			expected: expected{
				err: "invalid JSONPath query '@.[1, 2]' invalid token. '[1, 2]' does not match any token format",
			},
		},
		{
			input: input{
				queryPath: "@.[1, 2]",
				isStrict:  false,
			},
			expected: expected{
				tokens: 2,
			},
		},
		{
			input: input{
				queryPath: "@.length<1",
			},
			expected: expected{
				tokens: 2,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			jsonPath, err := Compile(test.input.queryPath, test.input.isStrict)
			if test.expected.err != "" {
				assert.Nil(t, jsonPath)
				assert.EqualError(t, err, test.expected.err)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, jsonPath)

			assert.NotNil(t, jsonPath.options)
			assert.Equal(t, test.input.isStrict, jsonPath.options.IsStrict)
			assert.Equal(t, test.input.queryPath, jsonPath.queryString)
			assert.Len(t, jsonPath.tokens, test.expected.tokens)
		})
	}

}

func Test_Query(t *testing.T) {

	jsonData := make(map[string]interface{})
	json.Unmarshal([]byte(sampleDataString), &jsonData)

	type input struct {
		queryString string
		jsonData    map[string]interface{}
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
				queryString: "invalid",
			},
			expected: expected{
				err: "invalid JSONPath query 'invalid' unexpected token 'i' at index 0",
			},
		},
		{
			input: input{
				queryString: "$.expensive",
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				queryString: "$.expensive",
				jsonData: map[string]interface{}{
					"expensive": "test",
				},
			},
			expected: expected{
				value: "test",
			},
		},
		{
			input: input{
				queryString: "$.expensive",
				jsonData:    jsonData,
			},
			expected: expected{
				value: int64(10),
			},
		},
		{
			input: input{
				queryString: "$..author",
				jsonData:    jsonData,
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
			value, err := Query(test.input.queryString, test.input.jsonData)

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

func Test_QueryString(t *testing.T) {

	type input struct {
		queryString string
		jsonData    string
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
				queryString: "$.expensive",
			},
			expected: expected{
				err: "invalid data. unexpected end of JSON input",
			},
		},
		{
			input: input{
				queryString: "invalid",
				jsonData:    "{}",
			},
			expected: expected{
				err: "invalid JSONPath query 'invalid' unexpected token 'i' at index 0",
			},
		},
		{
			input: input{

				queryString: "$.expensive",
				jsonData:    "{}",
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				queryString: "$.expensive",
				jsonData:    `{"expensive": "test"}`,
			},
			expected: expected{
				value: "test",
			},
		},
		{
			input: input{
				queryString: "$.expensive",
				jsonData:    sampleDataString,
			},
			expected: expected{
				value: int64(10),
			},
		},
		{
			input: input{
				queryString: "$..author",
				jsonData:    sampleDataString,
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
				queryString: "$.store.book.length",
				jsonData:    sampleDataString,
			},
			expected: expected{
				value: int64(4),
			},
		},
		{
			input: input{
				queryString: "$..book.length",
				jsonData:    sampleDataString,
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
			value, err := QueryString(test.input.queryString, test.input.jsonData)

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

func Test_QueryObject(t *testing.T) {

	type input struct {
		queryString string
		jsonData    interface{}
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
				queryString: "$.expensive",
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				queryString: "invalid",
				jsonData:    &sampleData{},
			},
			expected: expected{
				err: "invalid JSONPath query 'invalid' unexpected token 'i' at index 0",
			},
		},
		{
			input: input{

				queryString: "$.expensive",
				jsonData:    &storeData{},
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				queryString: "$.expensive",
				jsonData:    &sampleData{Expensive: 15},
			},
			expected: expected{
				value: float64(15),
			},
		},
		{
			input: input{
				queryString: "$.expensive",
				jsonData:    sampleDataObject,
			},
			expected: expected{
				value: float64(10),
			},
		},
		{
			input: input{
				queryString: "$..author",
				jsonData:    sampleDataObject,
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
				queryString: "$.store.book.length",
				jsonData:    sampleDataObject,
			},
			expected: expected{
				value: int64(4),
			},
		},
		{
			input: input{
				queryString: "$..book.length",
				jsonData:    sampleDataObject,
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
			value, err := QueryObject(test.input.queryString, test.input.jsonData)

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

func Test_JSONPath_compile(t *testing.T) {

	type expected struct {
		err    string
		tokens int
	}

	tests := []struct {
		input    string
		expected expected
	}{
		{
			input: "",
			expected: expected{
				err: "unexpected token '' at index 0",
			},
		},
		{
			input: "@.[1, 2]",
			expected: expected{
				err: "invalid token. '[1, 2]' does not match any token format",
			},
		},
		{
			input: "@.length<1",
			expected: expected{
				tokens: 2,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			jsonPath := &JSONPath{}
			actual := jsonPath.compile(test.input)
			if test.expected.err == "" {
				assert.Nil(t, actual)
			} else {
				assert.EqualError(t, actual, test.expected.err)
			}

			assert.Len(t, jsonPath.tokens, test.expected.tokens)
		})
	}
}

func Test_JSONPath_Query(t *testing.T) {

	sampleQuery, _ := Compile("$.expensive", false)
	altSampleQuery, _ := Compile("$..author", false)

	jsonData := make(map[string]interface{})
	json.Unmarshal([]byte(sampleDataString), &jsonData)

	type input struct {
		jsonPath *JSONPath
		jsonData map[string]interface{}
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
				jsonPath: &JSONPath{
					queryString: "invalid",
				},
			},
			expected: expected{
				err: "invalid JSONPath query 'invalid'",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
				jsonData: map[string]interface{}{
					"expensive": "test",
				},
			},
			expected: expected{
				value: "test",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
				jsonData: jsonData,
			},
			expected: expected{
				value: int64(10),
			},
		},
		{
			input: input{
				jsonPath: altSampleQuery,
				jsonData: jsonData,
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
			value, err := test.input.jsonPath.Query(test.input.jsonData)

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

func Test_JSONPath_QueryString(t *testing.T) {

	sampleQuery, _ := Compile("$.expensive", false)
	altSampleQuery, _ := Compile("$..author", false)

	type input struct {
		jsonPath *JSONPath
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
				jsonPath: sampleQuery,
			},
			expected: expected{
				err: "invalid data. unexpected end of JSON input",
			},
		},
		{
			input: input{
				jsonPath: &JSONPath{
					queryString: "invalid",
				},
				jsonData: "{}",
			},
			expected: expected{
				err: "invalid JSONPath query 'invalid'",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
				jsonData: "{}",
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
				jsonData: `{"expensive": "test"}`,
			},
			expected: expected{
				value: "test",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
				jsonData: sampleDataString,
			},
			expected: expected{
				value: int64(10),
			},
		},
		{
			input: input{
				jsonPath: altSampleQuery,
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
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			value, err := test.input.jsonPath.QueryString(test.input.jsonData)

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

func Test_JSONPath_QueryObject(t *testing.T) {

	sampleQuery, _ := Compile("$.expensive", false)
	altSampleQuery, _ := Compile("$..author", false)

	type input struct {
		jsonPath *JSONPath
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
				jsonPath: sampleQuery,
				jsonData: make(chan bool, 1),
			},
			expected: expected{
				err: "invalid data. json: unsupported type: chan bool",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
				jsonData: "not something that can be marshalled",
			},
			expected: expected{
				err: "invalid data. json: cannot unmarshal string into Go value of type map[string]interface {}",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				jsonPath: &JSONPath{
					queryString: "invalid",
				},
				jsonData: &sampleData{},
			},
			expected: expected{
				err: "invalid JSONPath query 'invalid'",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
				jsonData: &bookData{},
			},
			expected: expected{
				err: "key: invalid token key 'expensive' not found",
			},
		},
		{
			input: input{
				jsonPath: sampleQuery,
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
				jsonPath: sampleQuery,
				jsonData: sampleDataObject,
			},
			expected: expected{
				value: float32(10),
			},
		},
		{
			input: input{
				jsonPath: altSampleQuery,
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
			value, err := test.input.jsonPath.QueryObject(test.input.jsonData)

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
