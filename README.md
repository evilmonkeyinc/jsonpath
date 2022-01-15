[![codecov](https://codecov.io/gh/evilmonkeyinc/jsonpath/branch/main/graph/badge.svg?token=4PU85I7J2R)](https://codecov.io/gh/evilmonkeyinc/jsonpath)
[![main](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml)
[![develop](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml/badge.svg?branch=develop)](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/evilmonkeyinc/jsonpath.svg)](https://pkg.go.dev/github.com/evilmonkeyinc/jsonpath)

> This library is on the unstable version v0.X.X, which means there is a chance that any minor update may introduce a breaking change. Where I will endeavor to avoid this, care should be taken updating your dependency on this library until the first stable release v1.0.0 at which point any future breaking changes will result in a new major release.

# JSONPath

Golang JSONPath parser

## Install

`go get github.com/evilmonkeyinc/jsonpath`

## Usage

```golang
package main

import (
	"fmt"
	"os"

	"github.com/evilmonkeyinc/jsonpath"
)

func main() {
	query := os.Args[1]
	data := os.Args[2]

	result, err := jsonpath.QueryString(query, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result)
	os.Exit(0)
}
```

## Functions

The following functions are exported to support the functionality

### Compile

Will parse a JSONPath query and return a JSONPath object that can be used to query multiple JSON data objects or strings

### Query

Will compile a JSONPath query and will query the supplied JSON data in any various formats.

The parser can support querying struct types, and will use the `json` tags for struct fields if they are present, if not it will use the names as they appear in the golang code.

### QueryString

Will compile a JSONPath query and will query the supplied JSON data. 

QueryString can support a JSON array or object strings, and will unmarshal them to `[]interface{}` or `map[string]interface{}` using the standard `encoding/json` package unmarshal functions.

## Types

### JSONPath

This object is returned by the `Compile` function.

The JSONPath struct represents a reusable compiled JSONPath query which supports the `Query`, and `QueryString` functions as detailed above.

## Supported Syntax

| syntax | name | description | example | notes |
| --- | --- | --- | --- | --- |
| `$` | root | represents the data object being queried  | `$` would return the root object | this should always be the first token in a query |
| `.` | child | used as a separator in the query, signaling that the next token is child of the preceding token | `$.store` would return the object with the key `store` in the root object | the subscript operator can be used instead of the child operator to support child tokens with special characters in them i.e. `$['child key']` |
| `..` | recursive | recursive child token.  | `$..book` would return every entry that is defined by the key `book` regardless of where it is in the data structure |  |
| `*` | any/all | a wildcard operator used to denote that you want all the child members | `$.store.book.*` returns all the members of the book array or map | can also be denoted with the subscript syntax `$.store.book[*]` |
| `[]` | subscript | allows for additional operators to be applied to the current object | `$.store.book[1]` returns the second entry in the book array | it is possible to use indexes to reference elements in a map, the order is determined by the keys in alphabetical order |
| `[,]` | union | allows for a comma separated list of indices or keys to denote the elements to return | `$.store.book[0,1]` returns the first two entries in the book array | it is possible to use script expressions to define the union keys i.e. `$.store.book[0,(@.length-1)]` returns the first and last elements of the book array |
| `[start:end:step]` | range | allows to define a range of elements in an array to return. the step operand allows you to skip alternating elements | `$.store.book[0:3:1)]` returns elements `0`, `1`, and `2` from the book array where `$.store.book[0:3:2)]` would return elements `0`, and `2` | it is possible to use script expressions to define the range keys i.e. `$.store.book[1:(@.length-1)]:1` returns the elements of the book array excluding the first and last element |
| `[?()]` | filter | evaluates the filters expression to return if the element should be returned | `$.store.book[?(@.price > 10)]` will return only the elements in the book array that have a `price` greater than 10 | a filter should return a boolean, but if a non-boolean value is returned  |
| `[()]` | script | evaluates the scripts expression to return the key or index for the target element | `$.store.book[(@.length-1)]` returns the last element of the book array | a script must return either an integer index or, if the preceding object was a map, a string key |
| `@` | current | represents the current object | `(@.length-1)`| only used in scripts and filters, and will represent different things depending where it is used. in a script it will represent the object that preceded it (the array or object), in a filter it will represent the elements of the preceding object (the elements in the array or map) |

### Range Variations

The range operation can be performed with various arguments.

`[start:end:step]` - the standard range operation, start is inclusive, end is exclusive, and step must be non-zero
`[start:end]` - range operation with step set to 1, this would be the same as `[start:end:1]`
`[start:]` - range operation with step set to 1 and the end set to the end of the collection, this would be the same as `[start:(@.length)]`
`[:end]` - range operation when start is treated as 0 and step is 1, this would be the same as `[0:end]`
`[start:end:-1]` - special range operation. the required elements are identified as with the standard range but will be returned in the reverse order, for example if you requested `[0:2:-1]` it would return elements `0`, `1`, and `2` but in the order `2,1,0`. you can specify values lower than -1 to also step over elements, for example `[0:2:-2]` would return elements `2` and `0` of the array, skipping `1`.

### Special Syntax

`.length` this child operator will allow you to return the length of an array, map, slice, or string. if used with a map that has a key `length` it will return the corresponding value instead of the length of the map

`[-1]` any time you can specify an index, either for a subscript, union, or range, you can also specify a negative value. this is used to retrieve the elements at the end of the collection instead of the start, `-1` would represent the last item in the array, `-2` the second last, and so on.

`string[0]` it is possible to get a character or substring from a string value using the subscript index, union, or range operations on a string. whereas these operations would normally return an array, they will instead return the modified string. for example if you applied `[0:2]` to a string `string` it would return `str`


## Supported standard evaluation operations

| symbol | name | supported types | example | notes |
| --- | --- | --- | --- | --- |
| == | equals | any | 1 == 1 returns true | |
| != | not equals | any | 1 != 2 returns true | |
| * | multiplication | int\|float | 2*2 returns 4 | |
| / | division | int\|float | 10/5 returns 2 | if you supply two whole numbers you will only get a whole number response, even if there is a remainder i.e. 10/4 would return 2, not 2.5. to include remainders you would need to have the numerator as a float i.e. 10.0/4 would return 2.5 |
| + | addition | int\|float | 2+2 returns 4 | |
| - | subtraction | int\|float | 2-2 returns 0 | |
| % | remainder | int\|float | 5 % 2 returns 1 | this operator will divide the numerator by the denominator and then return the remainder |
| > | greater than | int\|float | 1 > 0 returns true | |
| >= | greater than or equal to | int\|float | 1 >= 1 returns true | |
| < | less than | int\|float | 1 < 2 returns true | |
| <= | less than or equal to  | int\|float | 1 <= 1 returns true | |
| && | combine and | expression\|bool | true&&false returns false | evaluate two expressions that return true or false, and return true if both are true |
| \|\| | combine or | expression\|bool | true\|\|false returns true | evaluate two expressions that return true or false, and return true if either are true |
| (...) | sub-expression | expression | (1+2)*3 returns 9 | allows you to isolate a sub-expression so it will be evaluated first separate from the rest of the expression |

## History

The [original specification for JSONPath](https://goessner.net/articles/JsonPath/) was proposed in 2007, and was a programing challenge I had not attempted before while being a practical tool.

There are many [implementations](https://cburgmer.github.io/json-path-comparison/) in multiple languages so I will not claim that this library is better in any way but I believe that it is true to the original specification and was an enjoyable challenge.
