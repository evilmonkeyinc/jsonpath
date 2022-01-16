# cburgmer tests

This test package has tests detailed by https://cburgmer.github.io/json-path-comparison/ comparison matrix which details the community consensus on the expected response from multiple JSONPath queries

This implementation would be closer to the `Scalar consensus` as it does not always return an array, but instead can return a single item when that is expected.

## Array Test

|query|data|consensus|actual|match|
|---|---|---|---|---|
|`$[1:3]`|`["first", "second", "third", "forth", "fifth"]`|`[second third]`|`[second third]`|:white_check_mark:|
|`$[0:5]`|`["first", "second", "third", "forth", "fifth"]`|`[first second third forth fifth]`|`[first second third forth fifth]`|:white_check_mark:|
|`$[7:10]`|`["first", "second", "third"]`|`[]`|`[]`|:white_check_mark:|
|`$[1:3]`|`{":": 42, "more": "string", "a": 1, "b": 2, "c": 3, "1:3": "nice"}`|`nil`|`nil`|:white_check_mark:|
|`$[1:10]`|`["first", "second", "third"]`|`[second third]`|`[second third]`|:white_check_mark:|
|`$[2:113667776004]`|`["first", "second", "third", "forth", "fifth"]`|`[third forth fifth]`|`[third forth fifth]`|:white_check_mark:|
|`$[2:-113667776004:-1]`|`["first", "second", "third", "forth", "fifth"]`|`none`|[]|:question:|
|`$[-113667776004:2]`|`["first", "second", "third", "forth", "fifth"]`|`[first second]`|`[first second]`|:white_check_mark:|
|`$[113667776004:2:-1]`|`["first", "second", "third", "forth", "fifth"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:-5]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:-4]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:-3]`|`[2, "a", 4, 5, 100, "nice"]`|`[4]`|`[4]`|:white_check_mark:|
|`$[-4:1]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:2]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[-4:3]`|`[2, "a", 4, 5, 100, "nice"]`|`[4]`|`[4]`|:white_check_mark:|
|`$[3:0:-2]`|`["first", "second", "third", "forth", "fifth"]`|`[]`|`[]`|:white_check_mark:|
|`$[7:3:-1]`|`["first", "second", "third", "forth", "fifth"]`|`[]`|`[]`|:white_check_mark:|
|`$[0:3:-2]`|`["first", "second", "third", "forth", "fifth"]`|`none`|[third first]|:question:|
|`$[::-2]`|`["first", "second", "third", "forth", "fifth"]`|`none`|[fifth third first]|:question:|
|`$[1:]`|`["first", "second", "third", "forth", "fifth"]`|`[second third forth fifth]`|`[second third forth fifth]`|:white_check_mark:|
|`$[3::-1]`|`["first", "second", "third", "forth", "fifth"]`|`none`|[fifth forth]|:question:|
|`$[:2]`|`["first", "second", "third", "forth", "fifth"]`|`[first second]`|`[first second]`|:white_check_mark:|
|`$[:]`|`["first","second"]`|`[first second]`|`[first second]`|:white_check_mark:|
|`$[:]`|`{":": 42, "more": "string"}`|`nil`|`nil`|:white_check_mark:|
|`$[::]`|`["first","second"]`|`[first second]`|`[first second]`|:white_check_mark:|
|`$[:2:-1]`|`["first", "second", "third", "forth", "fifth"]`|`none`|[second first]|:question:|
|`$[3:-4]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[3:-3]`|`[2, "a", 4, 5, 100, "nice"]`|`[]`|`[]`|:white_check_mark:|
|`$[3:-2]`|`[2, "a", 4, 5, 100, "nice"]`|`[5]`|`[5]`|:white_check_mark:|
|`$[2:1]`|`["first", "second", "third", "forth"]`|`[]`|`[]`|:white_check_mark:|
|`$[0:0]`|`["first", "second"]`|`[]`|`[]`|:white_check_mark:|
|`$[0:1]`|`["first", "second"]`|`[first]`|`[first]`|:white_check_mark:|
|`$[-1:]`|`["first", "second", "third"]`|`[third]`|`[third]`|:white_check_mark:|
|`$[-2:]`|`["first", "second", "third"]`|`[second third]`|`[second third]`|:white_check_mark:|
|`$[-4:]`|`["first", "second", "third"]`|`[first second third]`|`[first second third]`|:white_check_mark:|
|`$[0:3:2]`|`["first", "second", "third", "forth", "fifth"]`|`[first third]`|`[first third]`|:white_check_mark:|
|`$[0:3:0]`|`["first", "second", "third", "forth", "fifth"]`|`none`|nil|:question:|
|`$[0:3:1]`|`["first", "second", "third", "forth", "fifth"]`|`[first second third]`|`[first second third]`|:white_check_mark:|
|`$[010:024:010]`|`[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]`|`[10 20]`|`[10 20]`|:white_check_mark:|
|`$[0:4:2]`|`["first", "second", "third", "forth", "fifth"]`|`[first third]`|`[first third]`|:white_check_mark:|
|`$[1:3:]`|`["first", "second", "third", "forth", "fifth"]`|`[second third]`|`[second third]`|:white_check_mark:|
|`$[::2]`|`["first", "second", "third", "forth", "fifth"]`|`[first third fifth]`|`[first third fifth]`|:white_check_mark:|

## Bracket Tests

|query|data|consensus|actual|match|
|---|---|---|---|---|
|`$['key']`|`{"key": "value"}`|`value`|`value`|:white_check_mark:|
|`$['missing']`|`{"key": "value"}`|`nil`|`nil`|:white_check_mark:|
|`$..[0]`|`[ "first", { "key": [ "first nested", { "more": [ { "nested": ["deepest", "second"] }, ["more", "values"] ] } ] } ]`|`[deepest first nested first more map[nested:[deepest second]]]`|`[deepest first nested first more map[nested:[deepest second]]]`|:white_check_mark:|
|`$['ü']`|`{"ü": 42}`|`nil`|`nil`|:white_check_mark:|
|`$['two.some']`|`{ "one": {"key": "value"}, "two": {"some": "more", "key": "other value"}, "two.some": "42" }`|`42`|`42`|:white_check_mark:|
|`$["key"]`|`{"key": "value"}`|`value`|`nil`|:no_entry:|
|`$[]`|`{"": 42, "''": 123, "\"\"": 222}`|`nil`|`nil`|:white_check_mark:|
|`$['']`|`{"": 42, "''": 123, "\"\"": 222}`|`42`|`42`|:white_check_mark:|
|`$[""]`|`{"": 42, "''": 123, "\"\"": 222}`|`42`|`nil`|:no_entry:|
|`$[-2]`|`["one element"]`|`nil`|`nil`|:white_check_mark:|
|`$[2]`|`["first", "second", "third", "forth", "fifth"]`|`third`|`third`|:white_check_mark:|
|`$[0]`|`{ "0": "value" }`|`nil`|`nil`|:white_check_mark:|
|`$[1]`|`["one element"]`|`nil`|`nil`|:white_check_mark:|
|`$[0]`|`Hello World`|`nil`|`nil`|:white_check_mark:|
|`$.*[1]`|`[[1], [2,3]]`|`[3]`|`[3]`|:white_check_mark:|
|`$[-1]`|`["first", "second", "third"]`|`third`|`third`|:white_check_mark:|
|`$[-1]`|`[]`|`nil`|`nil`|:white_check_mark:|
|`$[0]`|`["first", "second", "third", "forth", "fifth"]`|`first`|`first`|:white_check_mark:|
|`$[':']`|`{ ":": "value", "another": "entry" }`|`value`|`value`|:white_check_mark:|
|`$[']']`|`{"]": 42}`|`42`|`42`|:white_check_mark:|
|`$['@']`|`{ "@": "value", "another": "entry" }`|`value`|`value`|:white_check_mark:|
|`$['.']`|`{ ".": "value", "another": "entry" }`|`value`|`value`|:white_check_mark:|
|`$['.*']`|`{"key": 42, ".*": 1, "": 10}`|`1`|`1`|:white_check_mark:|
|`$['"']`|`{ "\"": "value", "another": "entry" }`|`value`|`value`|:white_check_mark:|
|`$['\\']`|`{"\\":"value"}`|`none`|value|:question:|
|`$['\'']`|`{"'":"value"}`|`value`|`value`|:white_check_mark:|
|`$['0']`|`{ "0": "value" }`|`value`|`value`|:white_check_mark:|
|`$['$']`|`{ "$": "value", "another": "entry" }`|`value`|`value`|:white_check_mark:|
|`$[':@."$,*\'\\']`|`{":@.\"$,*'\\": 42}`|`none`|42|:question:|
|`$['single'quote']`|`{"single'quote":"value"}`|`nil`|`nil`|:white_check_mark:|
|`$[',']`|`{ ",": "value", "another": "entry" }`|`value`|`value`|:white_check_mark:|
|`$['*']`|`{ "*": "value", "another": "entry" }`|`value`|`value`|:white_check_mark:|
|`$['*']`|`{ "another": "entry" }`|`nil`|`nil`|:white_check_mark:|
|`$[ 'a' ]`|`{" a": 1, "a": 2, " a ": 3, "a ": 4, " 'a' ": 5, " 'a": 6, "a' ": 7, " \"a\" ": 8, "\"a\"": 9}`|`2`|`2`|:white_check_mark:|
|`$['ni.*']`|`{"nice": 42, "ni.*": 1, "mice": 100}`|`1`|`1`|:white_check_mark:|
|`$['two'.'some']`|`{ "one": {"key": "value"}, "two": {"some": "more", "key": "other value"}, "two.some": "42", "two'.'some": "43" }`|`nil`|`nil`|:white_check_mark:|
|`$[two.some]`|`{ "one": {"key": "value"}, "two": {"some": "more", "key": "other value"}, "two.some": "42" }`|`nil`|`nil`|:white_check_mark:|
|`$[*]`|`[ "string", 42, { "key": "value" }, [0, 1] ]`|`[string 42 map[key:value] [0 1]]`|`[string 42 map[key:value] [0 1]]`|:white_check_mark:|
|`$[*]`|`[]`|`[]`|`[]`|:white_check_mark:|
|`$[*]`|`{}`|`[]`|`[]`|:white_check_mark:|
|`$[*]`|`[ 40, null, 42 ]`|`[40 <nil> 42]`|`[40 <nil> 42]`|:white_check_mark:|
|`$[*]`|`{ "some": "string", "int": 42, "object": { "key": "value" }, "array": [0, 1] }`|`[string 42 map[key:value] [0 1]]`|`[string 42 map[key:value] [0 1]]`|:white_check_mark:|
|`$[0:2][*]`|`[[1, 2], ["a", "b"], [0, 0]]`|`[1 2 a b]`|`[[1 2] [a b]]`|:no_entry:|
|`$[*].bar[*]`|`[{"bar": [42]}]`|`[42]`|`[[42]]`|:no_entry:|
|`$..[*]`|`{ "key": "value", "another key": { "complex": "string", "primitives": [0, 1] } }`|`[value map[complex:string primitives:[0 1]] [0 1] string 0 1]`|`[value map[complex:string primitives:[0 1]] [0 1] string 0 1]`|:white_check_mark:|
|`$[key]`|`{ "key": "value" }`|`nil`|`nil`|:white_check_mark:|
