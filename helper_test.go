package jsonpath

import (
	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/script"
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

type testScriptEngine struct {
	value interface{}
}

func (engine *testScriptEngine) Compile(expression string, options *option.QueryOptions) (script.CompiledExpression, error) {
	return nil, nil
}

func (engine *testScriptEngine) Evaluate(root, current interface{}, expression string, options *option.QueryOptions) (interface{}, error) {
	return nil, nil
}
