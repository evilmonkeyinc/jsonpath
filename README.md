[![codecov](https://codecov.io/gh/evilmonkeyinc/jsonpath/branch/main/graph/badge.svg?token=4PU85I7J2R)](https://codecov.io/gh/evilmonkeyinc/jsonpath)
[![main](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml)
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

### Options

Part of the JSONPath object, Options allows you to specify what additional functionality, if any, that you want to enable while querying data.

You are able to enable index referencing support for maps for all tokens using `AllowMapReferenceByIndex` or use enable it for each token type individually.

You are able to enable index referencing support for strings for all tokens using `AllowStringReferenceByIndex` or use enable it for each token type individually.

## Supported Syntax

| syntax | name  | example |
| --- | ---  | --- |
| `$` | root |  `$` | 
| `.` | child |  `$.store`  |
| `..` | recursive | `$..book`  |
| `*` | wildcard | `$.store.book.*` |
| `[]` | subscript |  `$.store.book[1]` | 
| `[,]` | union | `$.store.book[0,1]` | 
| `[start:end:step]` | range |  `$.store.book[0:3:1)]` | 
| `[?()]` | filter |  `$.store.book[?(@.price > 10)]` |
| `[()]` | script |  `$.store.book[(@.length-1)]` |
| `@` | current |  `(@.length-1)`| 

### Root

`$`

represents the data object being queried 

this should always be the first token in a query. It is also possible to use the root symbol in scripts and filters, for example `$.store.book[?(@.category == $.onSaleCategory)]` would allow you to filter the elements i the book array based on its `category` value compared to the `onSaleCategory` value on the root object.

### Child

`.key` or `['key']`

The child operator allows you to specify that you want the child element of a map or struct based on the elements key/name.

If the key, or field name, includes special characters including spaces then it is required to use the subscript with single quotes syntax. If the required key has a single quote in them then it can be escaped using `\`, for example `['key\'s']`.

### Recursive

`..key`

A recursive check through the data structure for the specified child element.

### Wildcard

`*` or `[*]`

a wildcard operator used to denote that you want all the child members of the previous object

can also be used with the subscript syntax `$.store.book[*]`

### Subscript

`[0]` or `['key']`

allows for additional operators to be applied to the current object to retrieve a child element.

A negative value for an index is supported, resulting in the elements being counted in reverse, `-1` would represent the last item in the collection, `-2` the second last, and so on.

### Union

`[0,1]` or `['first','second']`

allows for a comma separated list of indices or keys to denote the elements to return

It is possible to use script expressions to define the union keys i.e. `$.store.book[0,(@.length-1)]` returns the first and last elements of the book collection.

A negative value for an index is supported, resulting in the elements being counted in reverse, `-1` would represent the last item in the collection, `-2` the second last, and so on.

### Range

`[start:]` or `[:end]` or `[start:end]` or `[start:end:step]` 

Allows to define a range of elements in an array to return. Starting the first keys `start` up to, but not including, the second keys `end`. the the third keys `step` allows you to skip alternating elements.

It is possible to use script expressions to define the range keys i.e. `$.store.book[1:(@.length-1)]:1` returns the elements of the book array excluding the first and last element.

An empty keys are treated as:
1. `start` as `0`
2. `end` as the collection length
3. `step` as `1` 

A negative value for `start` or `end` is supported, resulting in the elements being counted in reverse, `-1` would represent the last item in the collection, `-2` the second last, and so on.

A negative value for `step` will return the results in the opposite order, but the range is still determined in the original order then it is reversed.

### Filter

`[?(expression)]`

Evaluates the filters expression to return if the element should be returned as part of the new array.

A filter expression should return a boolean, but if a non-nil value is returned it will also be treated as true, expect for an empty string which is considered false. This allows for filters such as `[?(@.isbn)]` where only the elements that have an `isbn` value would be included.

### Script

`[(expression)]`

Evaluates the scripts expression to return the key or index for the target element.

A script expression must return either an integer index or, if the preceding object was a map or struct, a string key or field name.

### Current

`@`

Only usable in scripts and filters, and will represent different things depending where it is used. 

In a script it will represent the object referenced by the previous token, allowing you to get the length of the array to determine an end index.

In a filter it will represent the child elements of the object referenced by the previous token, allowing you to determine if it should be included in the filtered array by referring to the child elements values.

### Length

`.length`

the length token will allow you to return the length of an array, map, slice, or string. 

If used with a map that has a key `length` it will return the corresponding value instead of the length of the map.

### Subscript, Union, and Range with maps and strings

Using the Compile() function, and modifying the JSONPath Options, it is possible to use a map or a string in place of an array with the subscript `[1]` union `[1,2,3]` and range `[0:3]` operations. 

For maps, the keys will be sorted into alphabetical order and they will be used to determine the index order. For example, if you had a map with strings `a` and `b`, regardless of the order, `a` would be the `0` index, and `b` the `1` index.

For strings instead of returning an array of characters instead will return a substring. For example if you applied `[0:3]` to the string `string` it would return `str`.

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

