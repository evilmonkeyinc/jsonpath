# Benchmarks

Benchmarks of this library compared to other golang implementations.

The example selectors from the original specification along with the sample data have been used to create this data.

If any library is not mentioned for a selector, that means that implementation returned an error of some kind

Test of accuracy are based off of the expected response based on the original specification and the consensus from the [json-path-comparison](https://cburgmer.github.io/json-path-comparison/)

## Command

```bash
go test -bench=. -cpu=1 -benchmem -count=1 -benchtime=100x
```

## Libraries

- `github.com/PaesslerAG/jsonpath v0.1.1`  
- `github.com/bhmj/jsonslice v1.1.2`  
- `github.com/evilmonkeyinc/jsonpath v0.7.0`  
- `github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852`  
- `github.com/spyzhov/ajson v0.7.0`  

## TL;DR

This implementation is slower than others, but is only one of two that has a non-error response to all sample selectors, the other being the [spyzhov/ajson](https://github.com/spyzhov/ajson) implementation which is on average twice as fast but relies on its own json marshaller (which is impressive in it's own right)

Generally the accuracy of the implementations that could run are the same, with a minor deviation with how array ranges are handled with one, one implementation ran but did not return a response I suspect the testing method is flawed but without adequate documentation I could not confirm this.

## Selectors

### `$.store.book[*].author`

Expected Response: `["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|43551 ns/op|6496 B/op|188 allocs/op|true|
|paesslerAG|25549 ns/op|6417 B/op|131 allocs/op|false|
|bhmj|6188 ns/op|1188 B/op|14 allocs/op|true|
|spyzhov|17612 ns/op|6608 B/op|127 allocs/op|true|


### `$..author`

Expected Response: `["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|198323 ns/op|16689 B/op|458 allocs/op|true|
|paesslerAG|16293 ns/op|6361 B/op|122 allocs/op|false|
|bhmj|16665 ns/op|1554 B/op|27 allocs/op|true|
|spyzhov|20614 ns/op|7912 B/op|159 allocs/op|true|


### `$.store.*`

Expected Response: `[too large]`
> the expected response is an array with two components, the bike object and and array containing the book

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|25933 ns/op|4928 B/op|130 allocs/op|true|
|paesslerAG|16568 ns/op|6233 B/op|120 allocs/op|false|
|bhmj|8429 ns/op|3708 B/op|9 allocs/op|true|
|spyzhov|13288 ns/op|6376 B/op|117 allocs/op|true|

### `$.store..price`

Expected Response; `[19.95,8.95,12.99,8.99,22.99]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|104648 ns/op|15673 B/op|443 allocs/op|true|
|paesslerAG|99737 ns/op|6297 B/op|125 allocs/op|false|
|bhmj|90572 ns/op|1195 B/op|28 allocs/op|true|
|spyzhov|23793 ns/op|7816 B/op|158 allocs/op|true|

### `$..book[2]`

Expected Response: `[{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"}]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|192611 ns/op|16961 B/op|471 allocs/op|true|
|paesslerAG|25408 ns/op|6545 B/op|130 allocs/op|false|
|bhmj|13719 ns/op|1260 B/op|16 allocs/op|true|
|spyzhov|130744 ns/op|7904 B/op|160 allocs/op|true|

### `$..book[(@.length-1)]`

Expected Response: `[{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|138309 ns/op|18001 B/op|542 allocs/op|true|
|spyzhov|47062 ns/op|8840 B/op|197 allocs/op|true|

### `$..book[-1:]`

Expected Response" `[{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|198634 ns/op|17201 B/op|486 allocs/op|true|
|paesslerAG|64934 ns/op|6801 B/op|137 allocs/op|false|
|bhmj|16392 ns/op|1709 B/op|22 allocs/op|note1|
|spyzhov|17658 ns/op|7968 B/op|164 allocs/op|true|

> note1: returned an array containing the expected response, an array in an array, but the correct object

### `$..book[0,1]`

Expected Response: `[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"}]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|111628 ns/op|17297 B/op|489 allocs/op|true|
|paesslerAG|54361 ns/op|6817 B/op|136 allocs/op|false|
|bhmj|23537 ns/op|2285 B/op|23 allocs/op|note1|
|spyzhov|49976 ns/op|8048 B/op|165 allocs/op|true|

> note1: returned an array containing the expected response, an array in an array, but the correct object

### `$..book[:2]`

Expected Response: `[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"}]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|138072 ns/op|17281 B/op|483 allocs/op|true|
|paesslerAG|28601 ns/op|6801 B/op|137 allocs/op|false|
|bhmj|21478 ns/op|2349 B/op|24 allocs/op|note1|
|spyzhov|77671 ns/op|7984 B/op|164 allocs/op|true|

> note1: returned an array containing the expected response, an array in an array, but the correct object

### `$..book[?(@.isbn)]`

Expected Response: ` [{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"},{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|211344 ns/op|20265 B/op|556 allocs/op|true|
|paesslerAG|138063 ns/op|6937 B/op|143 allocs/op|false|
|bhmj|78538 ns/op|2731 B/op|30 allocs/op|note1|
|spyzhov|71054 ns/op|8864 B/op|217 allocs/op|true|

> note1: returned an array containing the expected response, an array in an array, but the correct object

### `$..book[?(@.price<10)]`

Expected Response: `{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"}]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|282446 ns/op|20153 B/op|564 allocs/op|true|
|bhmj|79741 ns/op|2899 B/op|43 allocs/op|note1|
|spyzhov|79312 ns/op|10160 B/op|263 allocs/op|true|

> note1: returned an array containing the expected response, an array in an array, but the correct object

### `$..book[?(@.price<$.expensive)]`

Expected Response: `{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"}]`

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|305200 ns/op|21449 B/op|628 allocs/op|true|
|bhmj|147911 ns/op|2995 B/op|46 allocs/op|note1|
|spyzhov|232748 ns/op|10088 B/op|285 allocs/op|true|

> note1: returned an array containing the expected response, an array in an array, but the correct object

### `$..*`

Expected Response: `[too large]`
> the expected response is an array that contains every value from the sample data, this will include an array, objects, and then each individual element of those collections

|library|ns/op|B/op|allocs/op|accurate|
|-|-|-|-|-|
|evilmonkeyinc|144373 ns/op|20193 B/op|546 allocs/op|true|
|paesslerAG|32120 ns/op|6216 B/op|117 allocs/op|false|
|bhmj|78242 ns/op|31209 B/op|69 allocs/op|true|
|spyzhov|71936 ns/op|9288 B/op|187 allocs/op|true|
