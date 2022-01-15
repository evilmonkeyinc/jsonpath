[![codecov](https://codecov.io/gh/evilmonkeyinc/jsonpath/branch/main/graph/badge.svg?token=4PU85I7J2R)](https://codecov.io/gh/evilmonkeyinc/jsonpath)
[![main](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml)
[![develop](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml/badge.svg?branch=develop)](https://github.com/evilmonkeyinc/jsonpath/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/evilmonkeyinc/jsonpath.svg)](https://pkg.go.dev/github.com/evilmonkeyinc/jsonpath)

> This library is on the unstable version v0.X.X, which means there is a chance that any minor update may introduce a breaking change. Where I will endeavor to avoid this, care should be taken updating your dependency on this library until the first stable release v1.0.0 at which point any future breaking changes will result in a new major release.

# JSONPath

Golang JSONPath parser

## Install

`go get github.com/evilmonkeyinc/jsonpath`

## Functions

The following functions are exported to support the functionality

### Compile

Will parse a JSONPath query and return a JSONPath object that can be used to query multiple JSON data objects or strings

### Query

Will compile a JSONPath query and will query the supplied JSON data in `map[string]interface{}` format.

This function is supplied to help support custom JSON parsing other than the standard `encoding/json` package.

### QueryString

Will compile a JSONPath query and will query the supplied JSON data.

The JSON data will be unmarshaled to a queryable format using the standard `encoding/json` package.

### QueryObject

Will compile a JSONPath query and will query the supplied JSON data in `interface{}` format.

The JSON data will be marshaled/unmarshaled to a queryable format using the standard `encoding/json` package.

## Types

### JSONPath

This object is returned by the `Compile` function.

The JSONPath struct represents a reusable compiled JSONPath query which supports the `Query`, `QueryString`, and `QueryObject` functions as detailed above.

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
| > | greater than | | | |
| >= | greater than or equal to | | | |
| < | less than | | | |
| <= | less than or equal to  | | | |
| && | combine and | expression\|bool | x | y |
| \|\| | combine or | expression\|bool | x | y |
| () | sub-expression | expression | x | y |

## History

The [original specification for JSONPath]((https://goessner.net/articles/JsonPath/)
) was proposed in 2007, and was a programing challenge I had not attempted before while being a practical tool.

There are many [implementations](https://cburgmer.github.io/json-path-comparison/) in multiple languages so I will not claim that this library is better in any way but I believe that it is true to the original specification and was an enjoyable challenge.
