package benchmark

import (
	"encoding/json"
	"testing"

	paesslerAG "github.com/PaesslerAG/jsonpath"
	bhmj "github.com/bhmj/jsonslice"
	emi "github.com/evilmonkeyinc/jsonpath"
	oliveagle "github.com/oliveagle/jsonpath"
	spyzhov "github.com/spyzhov/ajson"
)

var selectors = []string{
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

var sampleDataString string = `{ "store": { "book": [{ "category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95 }, { "category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99 }, { "category": "fiction", "author": "Herman Melville", "title": "Moby Dick", "isbn": "0-553-21311-3", "price": 8.99 }, { "category": "fiction", "author": "J. R. R. Tolkien", "title": "The Lord of the Rings", "isbn": "0-395-19395-8", "price": 22.99 } ], "bicycle": { "color": "red", "price": 19.95 } }, "expensive": 10 }`

func Benchmark_Comparison(b *testing.B) {

	for _, selector := range selectors {
		b.Run(selector, func(b *testing.B) {
			b.Run("evilmonkeyinc", func(b *testing.B) {
				var err error
				for i := 0; i < b.N; i++ {
					_, err = emi.QueryString(selector, sampleDataString)
				}
				if err != nil {
					b.SkipNow()
				}
			})
			b.Run("paesslerAG", func(b *testing.B) {
				var err error
				for i := 0; i < b.N; i++ {
					value := make(map[string]interface{})
					sampleData := json.Unmarshal([]byte(sampleDataString), &value)
					_, err = paesslerAG.Get(selector, sampleData)
				}
				if err != nil {
					b.SkipNow()
				}
			})
			b.Run("bhmj", func(b *testing.B) {
				var err error
				for i := 0; i < b.N; i++ {
					_, err = bhmj.Get([]byte(sampleDataString), selector)
				}
				if err != nil {
					b.SkipNow()
				}
			})
			b.Run("oliveagle", func(b *testing.B) {
				var err error
				for i := 0; i < b.N; i++ {
					value := make(map[string]interface{})
					sampleData := json.Unmarshal([]byte(sampleDataString), &value)

					var compiled *oliveagle.Compiled
					compiled, err = oliveagle.Compile(selector)
					if err == nil {
						_, err = compiled.Lookup(sampleData)
					}
				}
				if err != nil {
					b.SkipNow()
				}
			})
			b.Run("spyzhov", func(b *testing.B) {
				var err error
				for i := 0; i < b.N; i++ {
					root, _ := spyzhov.Unmarshal([]byte(sampleDataString))
					_, err = root.JSONPath(selector)
				}
				if err != nil {
					b.SkipNow()
				}
			})
		})
	}
}

/**
func Test_Comparison(t *testing.T) {

	for _, selector := range selectors {
		t.Run(selector, func(t *testing.T) {

			var response interface{}

			// evilmonkeyinc
			obj, _ := emi.QueryString(selector, sampleDataString)
			bytes, _ := json.Marshal(obj)
			response = string(bytes)
			fmt.Printf("%s %s %v\n", selector, "evilmonkeyinc", response)

			// paesslerAG
			value := interface{}(nil)
			sampleData := json.Unmarshal([]byte(sampleDataString), &value)
			response, _ = paesslerAG.Get(selector, sampleData)
			fmt.Printf("%s %s %v\n", selector, "paesslerAG", response)

			// bhmj
			bytes, _ = bhmj.Get([]byte(sampleDataString), selector)
			response = string(bytes)
			fmt.Printf("%s %s %v\n", selector, "bhmj", response)


			// oliveagle
			value = make(map[string]interface{})
			sampleData = json.Unmarshal([]byte(sampleDataString), &value)

			compiled, err := oliveagle.Compile(selector)
			if err == nil {
				response, _ = compiled.Lookup(sampleData)
				fmt.Printf("%s %s %v\n", selector, "oliveagle", response)
			} else {
				fmt.Printf("%s %s %v\n", selector, "oliveagle", "failed to compile")
			}


			// spyzhov
			root, _ := spyzhov.Unmarshal([]byte(sampleDataString))
			response, _ = root.JSONPath(selector)
			fmt.Printf("%s %s %v\n", selector, "spyzhov", response)
		})
	}
}
**/
