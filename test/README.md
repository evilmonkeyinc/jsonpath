# cburgmer tests

This test package has tests detailed by https://cburgmer.github.io/json-path-comparison/ comparison matrix which details the community consensus on the expected response from multiple JSONPath queries
	
This implementation would be closer to the 'Scalar consensus' as it does not always return an array, but instead can return a single item when that is expected.
	

## Array Test

|query|data|consensus|actual|match|
|---|---|---|---|---|
|`$[1:3]`|`["first", "second", "third", "forth", "fifth"]`|`["second","third"]`|`["second","third"]`|:white_check_mark:|
|`$[0:5]`|`["first", "second", "third", "forth", "fifth"]`|`["first","second","third","forth","fifth"]`|`["first","second","third","forth","fifth"]`|:white_check_mark:|
|`$[7:10]`|`["first", "second", "third"]`|`[]`|`[]`|:white_check_mark:|
|`$[1:3]`|`{":": 42, "more": "string", "a": 1, "b": 2, "c": 3, "1:3": "nice"}`|`nil`|`null`|:white_check_mark:|
|`$[1:10]`|`["first", "second", "third"]`|`["second","third"]`|`["second","third"]`|:white_check_mark:|
|`$[2:113667776004]`|`["first", "second", "third", "forth", "fifth"]`|`["third","forth","fifth"]`|`["third","forth","fifth"]`|:white_check_mark:|
|`$[2:-113667776004:-1]`|`["first", "second", "third", "forth", "fifth"]`|none|`[]`|:question:|
|`$[-113667776004:2]`|`["first", "second", "third", "forth", "fifth"]`|`["first","second"]`|`["first","second"]`|:white_check_mark:|
|`$[113667776004:2:-1]`|`["first", "second", "third", "forth", "fifth"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:-5]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:-4]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:-3]`|`[2, "a", 4, 5, 100, "nice"]`|`[4]`|`[4]`|:white_check_mark:|
|`$[-4:1]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:2]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:3]`|`[2, "a", 4, 5, 100, "nice"]`|`[4]`|`[4]`|:white_check_mark:|
|`$[3:0:-2]`|`["first", "second", "third", "forth", "fifth"]`|`[]`|`[]`|:white_check_mark:|
|`$[7:3:-1]`|`["first", "second", "third", "forth", "fifth"]`|`[]`|`[]`|:white_check_mark:|
|`$[0:3:-2]`|`["first", "second", "third", "forth", "fifth"]`|none|`["third","first"]`|:question:|
|`$[::-2]`|`["first", "second", "third", "forth", "fifth"]`|none|`["fifth","third","first"]`|:question:|
|`$[1:]`|`["first", "second", "third", "forth", "fifth"]`|`["second","third","forth","fifth"]`|`["second","third","forth","fifth"]`|:white_check_mark:|
|`$[3::-1]`|`["first", "second", "third", "forth", "fifth"]`|none|`["fifth","forth"]`|:question:|
|`$[:2]`|`["first", "second", "third", "forth", "fifth"]`|`["first","second"]`|`["first","second"]`|:white_check_mark:|
|`$[:]`|`["first","second"]`|`["first","second"]`|`["first","second"]`|:white_check_mark:|
|`$[:]`|`{":": 42, "more": "string"}`|`nil`|`null`|:white_check_mark:|
|`$[::]`|`["first","second"]`|`["first","second"]`|`["first","second"]`|:white_check_mark:|
|`$[:2:-1]`|`["first", "second", "third", "forth", "fifth"]`|none|`["second","first"]`|:question:|
|`$[3:-4]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[3:-3]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[3:-2]`|`[2, "a", 4, 5, 100, "nice"]`|`[5]`|`[5]`|:white_check_mark:|
|`$[2:1]`|`["first", "second", "third", "forth"]`|`[]`|`[]`|:white_check_mark:|
|`$[0:0]`|`["first", "second"]`|`[]`|`[]`|:white_check_mark:|
|`$[0:1]`|`["first", "second"]`|`["first"]`|`["first"]`|:white_check_mark:|
|`$[-1:]`|`["first", "second", "third"]`|`["third"]`|`["third"]`|:white_check_mark:|
|`$[-2:]`|`["first", "second", "third"]`|`["second","third"]`|`["second","third"]`|:white_check_mark:|
|`$[-4:]`|`["first", "second", "third"]`|`["first","second","third"]`|`["first","second","third"]`|:white_check_mark:|
|`$[0:3:2]`|`["first", "second", "third", "forth", "fifth"]`|`["first","third"]`|`["first","third"]`|:white_check_mark:|
|`$[0:3:0]`|`["first", "second", "third", "forth", "fifth"]`|none|`null`|:question:|
|`$[0:3:1]`|`["first", "second", "third", "forth", "fifth"]`|`["first","second","third"]`|`["first","second","third"]`|:white_check_mark:|
|`$[010:024:010]`|`[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]`|`[10,20]`|`[10,20]`|:white_check_mark:|
|`$[0:4:2]`|`["first", "second", "third", "forth", "fifth"]`|`["first","third"]`|`["first","third"]`|:white_check_mark:|
|`$[1:3:]`|`["first", "second", "third", "forth", "fifth"]`|`["second","third"]`|`["second","third"]`|:white_check_mark:|
|`$[::2]`|`["first", "second", "third", "forth", "fifth"]`|`["first","third","fifth"]`|`["first","third","fifth"]`|:white_check_mark:|

## Bracket Test

|query|data|consensus|actual|match|
|---|---|---|---|---|
|`$['key']`|`{"key": "value"}`|`"value"`|`"value"`|:white_check_mark:|
|`$['missing']`|`{"key": "value"}`|`nil`|`null`|:white_check_mark:|
|`$..[0]`|`[ "first", { "key": [ "first nested", { "more": [ { "nested": ["deepest", "second"] }, ["more", "values"] ] } ] } ]`|`["deepest","first nested","first","more",{"nested":["deepest","second"]}]`|`["deepest","first nested","first","more",{"nested":["deepest","second"]}]`|:white_check_mark:|
|`$['ü']`|`{"ü": 42}`|`nil`|`null`|:white_check_mark:|
|`$['two.some']`|`{ "one": {"key": "value"}, "two": {"some": "more", "key": "other value"}, "two.some": "42" }`|`"42"`|`"42"`|:white_check_mark:|
|`$["key"]`|`{"key": "value"}`|`"value"`|`null`|:no_entry:|
|`$[]`|`{"": 42, "''": 123, "\"\"": 222}`|`nil`|`null`|:white_check_mark:|
|`$['']`|`{"": 42, "''": 123, "\"\"": 222}`|`42`|`42`|:white_check_mark:|
|`$[""]`|`{"": 42, "''": 123, "\"\"": 222}`|`42`|`null`|:no_entry:|
|`$[-2]`|`["one element"]`|`nil`|`null`|:white_check_mark:|
|`$[2]`|`["first", "second", "third", "forth", "fifth"]`|`"third"`|`"third"`|:white_check_mark:|
|`$[0]`|`{ "0": "value" }`|`nil`|`null`|:white_check_mark:|
|`$[1]`|`["one element"]`|`nil`|`null`|:white_check_mark:|
|`$[0]`|`Hello World`|`nil`|`null`|:white_check_mark:|
|`$.*[1]`|`[[1], [2,3]]`|`[3]`|`[3]`|:white_check_mark:|
|`$[-1]`|`["first", "second", "third"]`|`"third"`|`"third"`|:white_check_mark:|
|`$[-1]`|`[]`|`nil`|`null`|:white_check_mark:|
|`$[0]`|`["first", "second", "third", "forth", "fifth"]`|`"first"`|`"first"`|:white_check_mark:|
|`$[':']`|`{ ":": "value", "another": "entry" }`|`"value"`|`"value"`|:white_check_mark:|
|`$[']']`|`{"]": 42}`|`42`|`42`|:white_check_mark:|
|`$['@']`|`{ "@": "value", "another": "entry" }`|`"value"`|`"value"`|:white_check_mark:|
|`$['.']`|`{ ".": "value", "another": "entry" }`|`"value"`|`"value"`|:white_check_mark:|
|`$['.*']`|`{"key": 42, ".*": 1, "": 10}`|`1`|`1`|:white_check_mark:|
|`$['"']`|`{ "\"": "value", "another": "entry" }`|`"value"`|`"value"`|:white_check_mark:|
|`$['\\']`|`{"\\":"value"}`|none|`"value"`|:question:|
|`$['\'']`|`{"'":"value"}`|`"value"`|`"value"`|:white_check_mark:|
|`$['0']`|`{ "0": "value" }`|`"value"`|`"value"`|:white_check_mark:|
|`$['$']`|`{ "$": "value", "another": "entry" }`|`"value"`|`"value"`|:white_check_mark:|
|`$[':@."$,*\'\\']`|`{":@.\"$,*'\\": 42}`|none|`42`|:question:|
|`$['single'quote']`|`{"single'quote":"value"}`|`nil`|`null`|:white_check_mark:|
|`$[',']`|`{ ",": "value", "another": "entry" }`|`"value"`|`"value"`|:white_check_mark:|
|`$['*']`|`{ "*": "value", "another": "entry" }`|`"value"`|`"value"`|:white_check_mark:|
|`$['*']`|`{ "another": "entry" }`|`nil`|`null`|:white_check_mark:|
|`$[ 'a' ]`|`{" a": 1, "a": 2, " a ": 3, "a ": 4, " 'a' ": 5, " 'a": 6, "a' ": 7, " \"a\" ": 8, "\"a\"": 9}`|`2`|`2`|:white_check_mark:|
|`$['ni.*']`|`{"nice": 42, "ni.*": 1, "mice": 100}`|`1`|`1`|:white_check_mark:|
|`$['two'.'some']`|`{ "one": {"key": "value"}, "two": {"some": "more", "key": "other value"}, "two.some": "42", "two'.'some": "43" }`|`nil`|`null`|:white_check_mark:|
|`$[two.some]`|`{ "one": {"key": "value"}, "two": {"some": "more", "key": "other value"}, "two.some": "42" }`|`nil`|`null`|:white_check_mark:|
|`$[*]`|`[ "string", 42, { "key": "value" }, [0, 1] ]`|`["string",42,{"key":"value"},[0,1]]`|`["string",42,{"key":"value"},[0,1]]`|:white_check_mark:|
|`$[*]`|`[]`|`[]`|`[]`|:white_check_mark:|
|`$[*]`|`{}`|`[]`|`[]`|:white_check_mark:|
|`$[*]`|`[ 40, null, 42 ]`|`[40,null,42]`|`[40,null,42]`|:white_check_mark:|
|`$[*]`|`{ "some": "string", "int": 42, "object": { "key": "value" }, "array": [0, 1] }`|`["string",42,{"key":"value"},[0,1]]`|`["string",42,{"key":"value"},[0,1]]`|:white_check_mark:|
|`$[0:2][*]`|`[[1, 2], ["a", "b"], [0, 0]]`|`[1,2,"a","b"]`|`[[1,2],["a","b"]]`|:no_entry:|
|`$[*].bar[*]`|`[{"bar": [42]}]`|`[42]`|`[[42]]`|:no_entry:|
|`$..[*]`|`{ "key": "value", "another key": { "complex": "string", "primitives": [0, 1] } }`|`["value",{"complex":"string","primitives":[0,1]},[0,1],"string",0,1]`|`["value",{"complex":"string","primitives":[0,1]},[0,1],"string",0,1]`|:white_check_mark:|
|`$[key]`|`{ "key": "value" }`|`nil`|`null`|:white_check_mark:|

## Dot Test

|query|data|consensus|actual|match|
|---|---|---|---|---|
|`@.a`|`{"a": 1}`|`nil`|`1`|:no_entry:|
|`$.['key']`|`{ "key": "value", "other": {"key": [{"key": 42}]} }`|`"value"`|`"value"`|:white_check_mark:|
|`$.["key"]`|`{ "key": "value", "other": {"key": [{"key": 42}]} }`|none|`null`|:question:|
|`$.[key]`|`{ "key": "value", "other": {"key": [{"key": 42}]} }`|none|`null`|:question:|
|`$.key`|`{ "key": "value" }`|`"value"`|`"value"`|:white_check_mark:|
|`$.key`|`[0, 1]`|`nil`|`null`|:white_check_mark:|
|`$.key`|`{ "key": ["first", "second"] }`|`["first","second"]`|`["first","second"]`|:white_check_mark:|
|`$.id`|`[{"id": 2}]`|`nil`|`null`|:white_check_mark:|
|`$.key`|`{ "key": {} }`|`{}`|`{}`|:white_check_mark:|
|`$.key`|`{ "key": null }`|`nil`|`null`|:white_check_mark:|
|`$.missing`|`{"key": "value"}`|`nil`|`null`|:white_check_mark:|
|`$[0:2].key`|`[{"key": "ey"}, {"key": "bee"}, {"key": "see"}]`|`["ey","bee"]`|`["ey","bee"]`|:white_check_mark:|
|`$..[1].key`|`{ "k": [{"key": "some value"}, {"key": 42}], "kk": [[{"key": 100}, {"key": 200}, {"key": 300}], [{"key": 400}, {"key": 500}, {"key": 600}]], "key": [0, 1] }`|`[200,42,500]`|`[200,42,500]`|:white_check_mark:|
|`$[*].a`|`[{"a": 1},{"a": 1}]`|`[1,1]`|`[1,1]`|:white_check_mark:|
|`$[*].a`|`[{"a": 1}]`|`[1]`|`[1]`|:white_check_mark:|
|`$[*].a`|`[{"a": 1},{"b": 1}]`|`[1]`|`[1]`|:white_check_mark:|
|`$[?(@.id==42)].name`|`[{"id": 42, "name": "forty-two"}, {"id": 1, "name": "one"}]`|`["forty-two"]`|`["forty-two"]`|:white_check_mark:|
|`$..key`|`{ "object": { "key": "value", "array": [ {"key": "something"}, {"key": {"key": "russian dolls"}} ] }, "key": "top" }`|`["russian dolls","something","top","value",{"key":"russian dolls"}]`|`["russian dolls","something","top","value",{"key":"russian dolls"}]`|:white_check_mark:|
|`$.store..price`|`{ "store": { "book": [ { "category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95 }, { "category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99 }, { "category": "fiction", "author": "Herman Melville", "title": "Moby Dick", "isbn": "0-553-21311-3", "price": 8.99 }, { "category": "fiction", "author": "J. R. R. Tolkien", "title": "The Lord of the Rings", "isbn": "0-395-19395-8", "price": 22.99 } ], "bicycle": { "color": "red", "price": 19.95 } } }`|`[12.99,19.95,22.99,8.95,8.99]`|`[12.99,19.95,22.99,8.95,8.99]`|:white_check_mark:|
|`$...key`|`{ "object": { "key": "value", "array": [ {"key": "something"}, {"key": {"key": "russian dolls"}} ] }, "key": "top" }`|`["russian dolls","something","top","value"]`|`["russian dolls","something","top","value",{"key":"russian dolls"}]`|:no_entry:|
|`$[0,2].key`|`[{"key": "ey"}, {"key": "bee"}, {"key": "see"}]`|`["ey","see"]`|`["ey","see"]`|:white_check_mark:|
|`$['one','three'].key`|`{ "one": {"key": "value"}, "two": {"k": "v"}, "three": {"some": "more", "key": "other value"} }`|`["value","other value"]`|`["value","other value"]`|:white_check_mark:|
|`$.key-dash`|`{ "key": 42, "key-": 43, "-": 44, "dash": 45, "-dash": 46, "": 47, "key-dash": "value", "something": "else" }`|`"value"`|`"value"`|:white_check_mark:|
|`$."key"`|`{ "key": "value", "\"key\"": 42 }`|none|`42`|:question:|
|`$.."key"`|`{ "object": { "key": "value", "\"key\"": 100, "array": [ {"key": "something", "\"key\"": 0}, {"key": {"key": "russian dolls"}, "\"key\"": {"\"key\"": 99}} ] }, "key": "top", "\"key\"": 42 }`|none|`[0,42,99,100,{"\"key\"":99}]`|:question:|
|`$.`|`{"key": 42, "": 9001, "''": "nice"}`|none|`null`|:question:|
|`$.in`|`{ "in": "value" }`|`"value"`|`"value"`|:white_check_mark:|
|`$.length`|`{ "length": "value" }`|`"value"`|`"value"`|:white_check_mark:|
|`$.length`|`[4, 5, 6]`|`nil`|`3`|:no_entry:|
|`$.null`|`{ "null": "value" }`|`"value"`|`"value"`|:white_check_mark:|
|`$.true`|`{ "true": "value" }`|`"value"`|`"value"`|:white_check_mark:|
|`$.$`|`{ "$": "value" }`|none|`{"$":"value"}`|:question:|
|`$.屬性`|`{ "屬性": "value" }`|`"value"`|`"value"`|:white_check_mark:|
|`$.2`|`["first", "second", "third", "forth", "fifth"]`|none|`null`|:question:|
|`$.2`|`{"a": "first", "2": "second", "b": "third"}`|`"second"`|`"second"`|:white_check_mark:|
|`$.-1`|`["first", "second", "third", "forth", "fifth"]`|`nil`|`null`|:white_check_mark:|
|`$.'key'`|`{ "key": "value", "'key'": 42 }`|none|`42`|:question:|
|`$..'key'`|`{ "object": { "key": "value", "'key'": 100, "array": [ {"key": "something", "'key'": 0}, {"key": {"key": "russian dolls"}, "'key'": {"'key'": 99}} ] }, "key": "top", "'key'": 42 }`|none|`[42,100,0,{"'key'":99},99]`|:question:|
|`$.'some.key'`|`{"some.key": 42, "some": {"key": "value"}, "'some.key'": 43}`|none|`43`|:question:|
|`$. a`|`{" a": 1, "a": 2, " a ": 3, "": 4}`|none|`2`|:question:|
|`$.*`|`[ "string", 42, { "key": "value" }, [0, 1] ]`|`["string",42,{"key":"value"},[0,1]]`|`["string",42,{"key":"value"},[0,1]]`|:white_check_mark:|
|`$.*`|`[]`|`[]`|`[]`|:white_check_mark:|
|`$.*`|`{}`|`[]`|`[]`|:white_check_mark:|
|`$.*`|`{ "some": "string", "int": 42, "object": { "key": "value" }, "array": [0, 1] }`|`["string",42,[0,1],{"key":"value"}]`|`["string",42,[0,1],{"key":"value"}]`|:white_check_mark:|
|`$.*.bar.*`|`[{"bar": [42]}]`|`[42]`|`[[42]]`|:no_entry:|
|`$.*.*`|`[[1, 2, 3], [4, 5, 6]]`|`[1,2,3,4,5,6]`|`[[1,2,3],[4,5,6]]`|:no_entry:|
|`$..*`|`{ "key": "value", "another key": { "complex": "string", "primitives": [0, 1] } }`|`["string","value",0,1,[0,1],{"complex":"string","primitives":[0,1]}]`|`["string","value",0,1,[0,1],{"complex":"string","primitives":[0,1]}]`|:white_check_mark:|
|`$..*`|`[ 40, null, 42 ]`|`[40,42,null]`|`[40,42,null]`|:white_check_mark:|
|`$..*`|`42`|`[]`|`[]`|:white_check_mark:|
|`$a`|`{"a": 1, "$a": 2}`|`nil`|`null`|:white_check_mark:|
|`.key`|`{ "key": "value" }`|`nil`|`null`|:white_check_mark:|
|`key`|`{ "key": "value" }`|`nil`|`null`|:white_check_mark:|

## Filter Test

|query|data|consensus|actual|match|
|---|---|---|---|---|
|`$[?(@.key)]`|`{"key": 42, "another": {"key": 1}}`|none|`[{"key":1}]`|:question:|
|`$..*[?(@.id>2)]`|`[ { "complext": { "one": [ { "name": "first", "id": 1 }, { "name": "next", "id": 2 }, { "name": "another", "id": 3 }, { "name": "more", "id": 4 } ], "more": { "name": "next to last", "id": 5 } } }, { "name": "last", "id": 6 } ]`|none|`[[],[],[{"id":5,"name":"next to last"}],[{"id":3,"name":"another"},{"id":4,"name":"more"}],[],[],[],[],[]]`|:question:|
|`$..[?(@.id==2)]`|`{"id": 2, "more": [{"id": 2}, {"more": {"id": 2}}, {"id": {"id": 2}}, [{"id": 2}]]}`|none|`[{"id":2},{"id":2},{"id":2},{"id":2}]`|:question:|
|`$[?(@.key+50==100)]`|`[{"key": 60}, {"key": 50}, {"key": 10}, {"key": -50}, {"key+50": 100}]`|none|`[{"key":50}]`|:question:|
|`$[?(@.key>42 && @.key<44)]`|`[ {"key": 42}, {"key": 43}, {"key": 44} ]`|`[{"key":43}]`|`[{"key":43}]`|:white_check_mark:|
|`$[?(@.key>0 && false)]`|`[ {"key": 1}, {"key": 3}, {"key": "nice"}, {"key": true}, {"key": null}, {"key": false}, {"key": {}}, {"key": []}, {"key": -1}, {"key": 0}, {"key": ""} ]`|none|`[]`|:question:|
|`$[?(@.key>0 && true)]`|`[ {"key": 1}, {"key": 3}, {"key": "nice"}, {"key": true}, {"key": null}, {"key": false}, {"key": {}}, {"key": []}, {"key": -1}, {"key": 0}, {"key": ""} ]`|none|`[{"key":1},{"key":3}]`|:question:|
|`$[?(@.key>43 \|\| @.key<43)]`|`[ {"key": 42}, {"key": 43}, {"key": 44} ]`|`[{"key":42},{"key":44}]`|`[{"key":42},{"key":44}]`|:white_check_mark:|
|`$[?(@.key>0 \|\| false)]`|`[ {"key": 1}, {"key": 3}, {"key": "nice"}, {"key": true}, {"key": null}, {"key": false}, {"key": {}}, {"key": []}, {"key": -1}, {"key": 0}, {"key": ""} ]`|none|`[{"key":1},{"key":3}]`|:question:|
|`$[?(@.key>0 \|\| true)]`|`[ {"key": 1}, {"key": 3}, {"key": "nice"}, {"key": true}, {"key": null}, {"key": false}, {"key": {}}, {"key": []}, {"key": -1}, {"key": 0}, {"key": ""} ]`|none|`[{"key":1},{"key":3},{"key":-1},{"key":0}]`|:question:|
|`$[?(@['key']==42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"some": "value"} ]`|`[{"key":42}]`|`[{"key":42}]`|:white_check_mark:|
|`$[?(@['@key']==42)]`|`[ {"@key": 0}, {"@key": 42}, {"key": 42}, {"@key": 43}, {"some": "value"} ]`|`[{"@key":42}]`|`[{"@key":42}]`|:white_check_mark:|
|`$[?(@[-1]==2)]`|`[[2, 3], ["a"], [0, 2], [2]]`|none|`[[0,2],[2]]`|:question:|
|`$[?(@[1]=='b')]`|`[["a", "b"], ["x", "y"]]`|`[["a","b"]]`|`[["a","b"]]`|:white_check_mark:|
|`$[?(@[1]=='b')]`|`{"1": ["a", "b"], "2": ["x", "y"]}`|none|`[["a","b"]]`|:question:|
|`$[?(@)]`|`[ "some value", null, "value", 0, 1, -1, "", [], {}, false, true ]`|none|`["some value","value",0,1,-1,true]`|:question:|
|`$[?(@.a && (@.b \|\| @.c))]`|`[ { "a": true }, { "a": true, "b": true }, { "a": true, "b": true, "c": true }, { "b": true, "c": true }, { "a": true, "c": true }, { "c": true }, { "b": true } ]`|none|`[]`|:question:|
|`[?(@.a && @.b \|\| @.c)]`|`[ { "a": true, "b": true }, { "a": true, "b": true, "c": true }, { "b": true, "c": true }, { "a": true, "c": true }, { "a": true }, { "b": true }, { "c": true }, { "d": true }, {} ]`|none|`null`|:question:|
|`$[?(@.key/10==5)]`|`[{"key": 60}, {"key": 50}, {"key": 10}, {"key": -50}, {"key/10": 5}]`|none|`[{"key":50}]`|:question:|
|`$[?(@.key-dash == 'value')]`|`[ { "key-dash": "value" } ]`|none|`[]`|:question:|
|`$[?(@.2 == 'second')]`|`[{"a": "first", "2": "second", "b": "third"}]`|none|`[{"2":"second","a":"first","b":"third"}]`|:question:|
|`$[?(@.2 == 'third')]`|`[["first", "second", "third", "forth", "fifth"]] `|none|`[]`|:question:|
|`$[?()]`|`[1, {"key": 42}, "value", null]`|`nil`|`null`|:white_check_mark:|
|`$[?(@.key==42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`|none|`[{"key":42}]`|:question:|
|`$[?(@==42)]`|`[ 0, 42, -1, 41, 43, 42.0001, 41.9999, null, 100 ]`|`[42]`|`[]`|:no_entry:|
|`$[?(@.key==43)]`|`[{"key": 42}]`|`[]`|`[]`|:white_check_mark:|
|`$[?(@.key==42)]`|`{ "a": {"key": 0}, "b": {"key": 42}, "c": {"key": -1}, "d": {"key": 41}, "e": {"key": 43}, "f": {"key": 42.0001}, "g": {"key": 41.9999}, "h": {"key": 100}, "i": {"some": "value"} }`|none|`[{"key":42}]`|:question:|
|`$[?(@.id==2)]`|`{"id": 2}`|none|`[]`|:question:|
|`$[?(@.d==["v1","v2"])]`|`[ { "d": [ "v1", "v2" ] }, { "d": [ "a", "b" ] }, { "d": "v1" }, { "d": "v2" }, { "d": {} }, { "d": [] }, { "d": null }, { "d": -1 }, { "d": 0 }, { "d": 1 }, { "d": "['v1','v2']" }, { "d": "['v1', 'v2']" }, { "d": "v1,v2" }, { "d": "[\"v1\", \"v2\"]" }, { "d": "[\"v1\",\"v2\"]" } ]`|none|`[]`|:question:|
|`$[?(@[0:1]==[1])]`|`[[1, 2, 3], [1], [2, 3], 1, 2]`|none|`[]`|:question:|
|`$[?(@.*==[1,2])]`|`[[1,2], [2,3], [1], [2], [1, 2, 3], 1, 2, 3]`|none|`[]`|:question:|
|`$[?(@.d==['v1','v2'])]`|`[ { "d": [ "v1", "v2" ] }, { "d": [ "a", "b" ] }, { "d": "v1" }, { "d": "v2" }, { "d": {} }, { "d": [] }, { "d": null }, { "d": -1 }, { "d": 0 }, { "d": 1 }, { "d": "['v1','v2']" }, { "d": "['v1', 'v2']" }, { "d": "v1,v2" }, { "d": "[\"v1\", \"v2\"]" }, { "d": "[\"v1\",\"v2\"]" } ]`|none|`[]`|:question:|
|`$[?((@.key<44)==false)]`|`[{"key": 42}, {"key": 43}, {"key": 44}]`|none|`[{"key":44}]`|:question:|
|`$[?(@.key==false)]`|`[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`|none|`[{"key":false}]`|:question:|
|`$[?(@.key==null)]`|`[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`|none|`[]`|:question:|
|`$[?(@[0:1]==1)]`|`[[1, 2, 3], [1], [2, 3], 1, 2]`|none|`[]`|:question:|
|`$[?(@[*]==2)]`|`[[1,2], [2,3], [1], [2], [1, 2, 3], 1, 2, 3]`|none|`[]`|:question:|
|`$[?(@.*==2)]`|`[[1,2], [2,3], [1], [2], [1, 2, 3], 1, 2, 3]`|none|`[]`|:question:|
|`$[?(@.key==-0.123e2)]`|`[{"key": -12.3}, {"key": -0.123}, {"key": -12}, {"key": 12.3}, {"key": 2}, {"key": "-0.123e2"}]`|none|`[{"key":-12.3}]`|:question:|
|`$[?(@.key==010)]`|`[{"key": "010"}, {"key": "10"}, {"key": 10}, {"key": 0}, {"key": 8}]`|none|`[{"key":8}]`|:question:|
|`$[?(@.d=={"k":"v"})]`|`[ { "d": { "k": "v" } }, { "d": { "a": "b" } }, { "d": "k" }, { "d": "v" }, { "d": {} }, { "d": [] }, { "d": null }, { "d": -1 }, { "d": 0 }, { "d": 1 }, { "d": "[object Object]" }, { "d": "{\"k\": \"v\"}" }, { "d": "{\"k\":\"v\"}" }, "v" ]`|none|`[]`|:question:|
|`$[?(@.key=="value")]`|`[ {"key": "some"}, {"key": "value"}, {"key": null}, {"key": 0}, {"key": 1}, {"key": -1}, {"key": ""}, {"key": {}}, {"key": []}, {"key": "valuemore"}, {"key": "morevalue"}, {"key": ["value"]}, {"key": {"some": "value"}}, {"key": {"key": "value"}}, {"some": "value"} ]`|`[{"key":"value"}]`|`[{"key":"value"}]`|:white_check_mark:|
|`$[?(@.key=="Motörhead")]`|`[ {"key": "something"}, {"key": "Mot\u00f6rhead"}, {"key": "mot\u00f6rhead"}, {"key": "Motorhead"}, {"key": "Motoo\u0308rhead"}, {"key": "motoo\u0308rhead"} ]`|`[{"key":"Motörhead"}]`|`[{"key":"Motörhead"}]`|:white_check_mark:|
|`$[?(@.key=="hi@example.com")]`|`[ {"key": "some"}, {"key": "value"}, {"key": "hi@example.com"} ]`|`[{"key":"hi@example.com"}]`|`[{"key":"hi@example.com"}]`|:white_check_mark:|
|`$[?(@.key=="some.value")]`|`[ {"key": "some"}, {"key": "value"}, {"key": "some.value"} ]`|`[{"key":"some.value"}]`|`[{"key":"some.value"}]`|:white_check_mark:|
|`$[?(@.key=='value')]`|`[ {"key": "some"}, {"key": "value"} ]`|`[{"key":"value"}]`|`[{"key":"value"}]`|:white_check_mark:|
|`$[?(@.key=="Mot\u00f6rhead")]`|`[ {"key": "something"}, {"key": "Mot\u00f6rhead"}, {"key": "mot\u00f6rhead"}, {"key": "Motorhead"}, {"key": "Motoo\u0308rhead"}, {"key": "motoo\u0308rhead"} ]`|none|`[{"key":"Motörhead"}]`|:question:|
|`$[?(@.key==true)]`|`[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`|none|`[{"key":true}]`|:question:|
|`$[?(@.key1==@.key2)]`|`[ {"key1": 10, "key2": 10}, {"key1": 42, "key2": 50}, {"key1": 10}, {"key2": 10}, {}, {"key1": null, "key2": null}, {"key1": null}, {"key2": null}, {"key1": 0, "key2": 0}, {"key1": 0}, {"key2": 0}, {"key1": -1, "key2": -1}, {"key1": "", "key2": ""}, {"key1": false, "key2": false}, {"key1": false}, {"key2": false}, {"key1": true, "key2": true}, {"key1": [], "key2": []}, {"key1": {}, "key2": {}}, {"key1": {"a": 1, "b": 2}, "key2": {"b": 2, "a": 1}} ]`|none|`[{"key1":10,"key2":10},{"key1":0,"key2":0},{"key1":-1,"key2":-1},{"key1":"","key2":""},{"key1":false,"key2":false},{"key1":true,"key2":true}]`|:question:|
|`$.items[?(@.key==$.value)]`|`{"value": 42, "items": [{"key": 10}, {"key": 42}, {"key": 50}]}`|none|`[{"key":42}]`|:question:|
|`$[?(@.key>42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`|none|`[{"key":43},{"key":42.0001},{"key":100}]`|:question:|
|`$[?(@.key>=42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`|none|`[{"key":42},{"key":43},{"key":42.0001},{"key":100}]`|:question:|
|`$[?(@.d in [2, 3])]`|`[{"d": 1}, {"d": 2}, {"d": 1}, {"d": 3}, {"d": 4}]`|`nil`|`[]`|:no_entry:|
|`$[?(2 in @.d)]`|`[{"d": [1, 2, 3]}, {"d": [2]}, {"d": [1]}, {"d": [3, 4]}, {"d": [4, 2]}]`|`nil`|`[]`|:no_entry:|
|`$[?(@.key<42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`|none|`[{"key":0},{"key":-1},{"key":41},{"key":41.9999}]`|:question:|
|`$[?(@.key<=42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`|none|`[{"key":0},{"key":42},{"key":-1},{"key":41},{"key":41.9999}]`|:question:|
|`$[?(@.key*2==100)]`|`[{"key": 60}, {"key": 50}, {"key": 10}, {"key": -50}, {"key*2": 100}]`|none|`[{"key":50}]`|:question:|
|`$[?(!(@.key==42))]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`|none|`[{"key":0},{"key":-1},{"key":41},{"key":43},{"key":42.0001},{"key":41.9999},{"key":100}]`|:question:|
|`$[?(!(@.key<42))]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`|none|`[{"key":42},{"key":43},{"key":42.0001},{"key":100}]`|:question:|
|`$[?(!@.key)]`|`[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`|none|`[{"key":false}]`|:question:|
|`$[?(@.key!=42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`|none|`[{"key":0},{"key":-1},{"key":1},{"key":41},{"key":43},{"key":42.0001},{"key":41.9999},{"key":100},{"key":420}]`|:question:|
|`$[*].bookmarks[?(@.page == 45)]^^^`|`[ { "title": "Sayings of the Century", "bookmarks": [{ "page": 40 }] }, { "title": "Sword of Honour", "bookmarks": [ { "page": 35 }, { "page": 45 } ] }, { "title": "Moby Dick", "bookmarks": [ { "page": 3035 }, { "page": 45 } ] } ]`|`nil`|`[[],[],[]]`|:no_entry:|
|`$[?(@.name=~/hello.*/)]`|`[ {"name": "hullo world"}, {"name": "hello world"}, {"name": "yes hello world"}, {"name": "HELLO WORLD"}, {"name": "good bye"} ]`|none|`[]`|:question:|
|`$[?(@.name=~/@.pattern/)]`|`[ {"name": "hullo world"}, {"name": "hello world"}, {"name": "yes hello world"}, {"name": "HELLO WORLD"}, {"name": "good bye"}, {"pattern": "hello.*"} ]`|none|`[]`|:question:|
|`$[?(@[*]>=4)]`|`[[1,2],[3,4],[5,6]]`|none|`[]`|:question:|
|`$.x[?(@[*]>=$.y[*])]`|`{"x":[[1,2],[3,4],[5,6]],"y":[3,4,5]}`|none|`[]`|:question:|
|`$[?(@.key=42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`|`nil`|`[]`|:no_entry:|
|`$[?(@.a[?(@.price>10)])]`|`[ { "a": [{"price": 1}, {"price": 3}] }, { "a": [{"price": 11}] }, { "a": [{"price": 8}, {"price": 12}, {"price": 3}] }, { "a": [] } ]`|none|`[]`|:question:|
|`$[?(@.address.city=='Berlin')]`|`[ { "address": { "city": "Berlin" } }, { "address": { "city": "London" } } ]`|`[{"address":{"city":"Berlin"}}]`|`[{"address":{"city":"Berlin"}}]`|:white_check_mark:|
|`$[?(@.key-50==-100)]`|`[{"key": 60}, {"key": 50}, {"key": 10}, {"key": -50}, {"key-50": -100}]`|none|`[{"key":-50}]`|:question:|
|`$[?(1==1)]`|`[1, 3, "nice", true, null, false, {}, [], -1, 0, ""]`|none|`[1,3,"nice",true,null,false,{},[],-1,0,""]`|:question:|
|`$[?(@.key===42)]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`|none|`[]`|:question:|
|`$[?(@.key)]`|`[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`|none|`[{"key":true},{"key":"value"},{"key":0},{"key":1},{"key":-1},{"key":42}]`|:question:|
|`$.*[?(@.key)]`|`[ { "some": "some value" }, { "key": "value" } ]`|none|`[[],[]]`|:question:|
|`$..[?(@.id)]`|`{"id": 2, "more": [{"id": 2}, {"more": {"id": 2}}, {"id": {"id": 2}}, [{"id": 2}]]}`|none|`[{"id":2},{"id":2},{"id":2},{"id":2}]`|:question:|
|`$[?(false)]`|`[1, 3, "nice", true, null, false, {}, [], -1, 0, ""]`|none|`[]`|:question:|
|`$[?(@..child)]`|`[{"key": [{"child": 1}, {"child": 2}]}, {"key": [{"child": 2}]}, {"key": [{}]}, {"key": [{"something": 42}]}, {}]`|none|`[]`|:question:|
|`$[?(null)]`|`[1, 3, "nice", true, null, false, {}, [], -1, 0, ""]`|none|`[]`|:question:|
|`$[?(true)]`|`[1, 3, "nice", true, null, false, {}, [], -1, 0, ""]`|none|`[1,3,"nice",true,null,false,{},[],-1,0,""]`|:question:|
|`$[?@.key==42]`|`[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`|none|`null`|:question:|
|`$[?(@.key)]`|`[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`|none|`[{"key":true},{"key":"value"},{"key":0},{"key":1},{"key":-1},{"key":42}]`|:question:|

## Misc Test

|query|data|consensus|actual|match|
|---|---|---|---|---|
|``|`{"a": 42, "": 21}`|`nil`|`null`|:white_check_mark:|
|`$.data.sum()`|`{"data": [1,2,3,4]}`|none|`null`|:question:|
|`$(key,more)`|`{"key": 1, "some": 2, "more": 3}`|`nil`|`null`|:white_check_mark:|
|`$..`|`[{"a": {"b": "c"}}, [0, 1]]`|none|`[[{"a":{"b":"c"}},[0,1]],{"a":{"b":"c"}},{"b":"c"},"c",[0,1],0,1]`|:question:|
|`$.key..`|`{"some key": "value", "key": {"complex": "string", "primitives": [0, 1]}}`|none|`[{"complex":"string","primitives":[0,1]},[0,1],0,1,"string"]`|:question:|
|`$`|`{ "key": "value", "another key": { "complex": [ "a", 1 ] } }`|`{"another key":{"complex":["a",1]},"key":"value"}`|`{"another key":{"complex":["a",1]},"key":"value"}`|:white_check_mark:|
|`$`|`42`|`42`|`42`|:white_check_mark:|
|`$`|`false`|`false`|`false`|:white_check_mark:|
|`$`|`true`|`true`|`true`|:white_check_mark:|
|`$[(@.length-1)]`|`["first", "second", "third", "forth", "fifth"]`|`nil`|`"fifth"`|:no_entry:|

## Union Test

|query|data|consensus|actual|match|
|---|---|---|---|---|
|`$[0,1]`|`["first", "second", "third"]`|`["first","second"]`|`["first","second"]`|:white_check_mark:|
|`$[0,0]`|`["a"]`|`["a","a"]`|`["a","a"]`|:white_check_mark:|
|`$['a','a']`|`{"a":1}`|`[1,1]`|`[1,1]`|:white_check_mark:|
|`$[?(@.key<3),?(@.key>6)]`|`[{"key": 1}, {"key": 8}, {"key": 3}, {"key": 10}, {"key": 7}, {"key": 2}, {"key": 6}, {"key": 4}]`|none|`[]`|:question:|
|`$['key','another']`|`{ "key": "value", "another": "entry" }`|`["value","entry"]`|`["value","entry"]`|:white_check_mark:|
|`$['missing','key']`|`{ "key": "value", "another": "entry" }`|`["value"]`|`["value"]`|:white_check_mark:|
|`$[:]['c','d']`|`[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`|`["cc1","dd1","cc2","dd2"]`|`[["cc1","dd1"],["cc2","dd2"]]`|:no_entry:|
|`$[0]['c','d']`|`[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`|`["cc1","dd1"]`|`["cc1","dd1"]`|:white_check_mark:|
|`$.*['c','d']`|`[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`|`["cc1","dd1","cc2","dd2"]`|`[["cc1","dd1"],["cc2","dd2"]]`|:no_entry:|
|`$..['c','d']`|`[{"c":"cc1","d":"dd1","e":"ee1"}, {"c": "cc2", "child": {"d": "dd2"}}, {"c": "cc3"}, {"d": "dd4"}, {"child": {"c": "cc5"}}]`|none|`["cc1","dd1","cc2","dd2","cc3","dd4","cc5"]`|:question:|
|`$[4,1]`|`[1,2,3,4,5]`|`[5,2]`|`[5,2]`|:white_check_mark:|
|`$.*[0,:5]`|`{ "a": [ "string", null, true ], "b": [ false, "string", 5.4 ] }`|none|`null`|:question:|
|`$[1:3,4]`|`[1,2,3,4,5]`|none|`null`|:question:|
|`$[ 0 , 1 ]`|`["first", "second", "third"]`|`["first","second"]`|`["first","second"]`|:white_check_mark:|
|`$[*,1]`|`["first", "second", "third", "forth", "fifth"]`|`nil`|`null`|:white_check_mark:|
