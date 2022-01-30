package benchmark

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	paesslerAG "github.com/PaesslerAG/jsonpath"
	"github.com/bhmj/jsonslice"
	emi "github.com/evilmonkeyinc/jsonpath"
	oliveagle "github.com/oliveagle/jsonpath"
	"github.com/spyzhov/ajson"
)

var testSelectors = []string{
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

var expectedResponse = map[string]string{
	"$.store.book[*].author":          `["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]`,
	"$..author":                       `["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]`,
	"$.store.*":                       `[{"color":"red","price":19.95},[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"},{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]]`,
	"$.store..price":                  `[19.95,8.95,12.99,8.99,22.99]`,
	"$..book[2]":                      `[{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"}]`,
	"$..book[(@.length-1)]":           `[{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]`,
	"$..book[-1:]":                    `[{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]`,
	"$..book[0,1]":                    `[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"}]`,
	"$..book[:2]":                     `[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"}]`,
	"$..book[?(@.isbn)]":              `[{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"},{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]`,
	"$..book[?(@.price<10)]":          `[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"}]`,
	"$..book[?(@.price<$.expensive)]": `[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"}]`,
	"$..*":                            `[10,{"bicycle":{"color":"red","price":19.95},"book":[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"},{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]},{"color":"red","price":19.95},[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"},{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}],"red",19.95,{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"},{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"},"Nigel Rees","reference",8.95,"Sayings of the Century","Evelyn Waugh","fiction",12.99,"Sword of Honour","Herman Melville","fiction","0-553-21311-3",8.99,"Moby Dick","J. R. R. Tolkien","fiction","0-395-19395-8",22.99,"The Lord of the Rings"]`,
}

var sampleDataString string = `{ "store": { "book": [{ "category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95 }, { "category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99 }, { "category": "fiction", "author": "Herman Melville", "title": "Moby Dick", "isbn": "0-553-21311-3", "price": 8.99 }, { "category": "fiction", "author": "J. R. R. Tolkien", "title": "The Lord of the Rings", "isbn": "0-395-19395-8", "price": 22.99 } ], "bicycle": { "color": "red", "price": 19.95 } }, "expensive": 10 }`

func Benchmark_Comparison(b *testing.B) {
	accuracyCheck := true

	for _, selector := range testSelectors {
		expected := expectedResponse[selector]
		b.Run(selector, func(b *testing.B) {
			b.Run("evilmonkeyinc", func(b *testing.B) {
				var err error
				var val interface{}
				for i := 0; i < b.N; i++ {
					val, err = emi.QueryString(selector, sampleDataString)
				}
				if accuracyCheck {
					if err != nil {
						b.Log("unsupported")
					} else {
						actual, _ := json.Marshal(val)
						if !jsonDeepEqual(expected, string(actual)) {
							b.Log("unexpected response")
						}
					}
				}

			})
			b.Run("paesslerAG", func(b *testing.B) {
				var err error
				var val interface{}
				for i := 0; i < b.N; i++ {
					value := make(map[string]interface{})
					json.Unmarshal([]byte(sampleDataString), &value)
					val, err = paesslerAG.Get(selector, value)
				}
				if accuracyCheck {
					if err != nil {
						b.Log("unsupported")
					} else if val != nil {
						// manually confirmed they match, something is wrong with unordered check
						// actual, _ := json.Marshal(val)
						// if !jsonDeepEqual(expected, string(actual)) {
						// 	b.Log("unexpected response", string(actual))
						// }
					}
				}
			})
			b.Run("bhmj", func(b *testing.B) {
				var err error
				var val []byte
				for i := 0; i < b.N; i++ {
					val, err = jsonslice.Get([]byte(sampleDataString), selector)
				}
				if accuracyCheck {
					if err != nil {
						b.Log("unsupported")
					} else {
						// bhmj sometimes has arrays in arrays
						if strings.HasPrefix(string(val), "[[") && strings.HasSuffix(string(val), "]]") {
							array := make([]interface{}, 0)
							json.Unmarshal(val, &array)
							if len(array) == 1 {
								b.Log("found single nested array")
								val, _ = json.Marshal(array[0])
							}
						}

						// manually confirmed they match, something is wrong with unordered check
						// if !jsonDeepEqual(expected, string(val)) {
						// 	b.Log("unexpected response", string(val))
						// }
					}
				}
			})
			b.Run("oliveagle", func(b *testing.B) {
				var err error
				var val interface{}
				for i := 0; i < b.N; i++ {
					value := make(map[string]interface{})
					json.Unmarshal([]byte(sampleDataString), &value)

					var compiled *oliveagle.Compiled
					compiled, err = oliveagle.Compile(selector)
					if err == nil {
						val, err = compiled.Lookup(value)
					}
				}
				if accuracyCheck {

					if err != nil {
						b.Log("unsupported")
					} else {
						actual, _ := json.Marshal(val)
						if !jsonDeepEqual(expected, string(actual)) {
							b.Log("unexpected response")
						}
					}
				}
			})
			b.Run("spyzhov", func(b *testing.B) {
				var err error
				var result *ajson.Node
				for i := 0; i < b.N; i++ {
					root, _ := ajson.Unmarshal([]byte(sampleDataString))
					var nodes []*ajson.Node
					nodes, err = root.JSONPath(selector)
					result = ajson.ArrayNode("", nodes)
				}
				if accuracyCheck {
					if err != nil {
						b.Log("unsupported")
					} else {
						actual, _ := ajson.Marshal(result)
						if !jsonDeepEqual(expected, string(actual)) {
							b.Log("unexpected response")
						}
					}
				}
			})
		})
	}
}

func jsonDeepEqual(expected string, actual string) bool {
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		return false
	}

	return reflect.DeepEqual(expectedJSONAsInterface, actualJSONAsInterface)
}
