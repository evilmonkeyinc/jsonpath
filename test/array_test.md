|query|data|consensus|actual|match|
|---|---|---|---|---|
|$[1:3]|["first", "second", "third", "forth", "fifth"]|[second third]|[second third]|:white_check_mark:|
|$[0:5]|["first", "second", "third", "forth", "fifth"]|[first second third forth fifth]|[first second third forth fifth]|:white_check_mark:|
|$[7:10]|["first", "second", "third"]|[]|[]|:white_check_mark:|
|$[1:3]|{":": 42, "more": "string", "a": 1, "b": 2, "c": 3, "1:3": "nice"}|[]|[42 1]|:no_entry:|
|$[1:10]|["first", "second", "third"]|[second third]|[second third]|:white_check_mark:|
|$[2:113667776004]|["first", "second", "third", "forth", "fifth"]|[third forth fifth]|[third forth fifth]|:white_check_mark:|
|$[2:-113667776004:-1]|["first", "second", "third", "forth", "fifth"]|none|[]|:question:|
|$[-113667776004:2]|["first", "second", "third", "forth", "fifth"]|[first second]|[first second]|:white_check_mark:|
|$[113667776004:2:-1]|["first", "second", "third", "forth", "fifth"]|[]|[]|:white_check_mark:|
|$[-4:-5]|[2, "a", 4, 5, 100, "nice"]|[]|[]|:white_check_mark:|
|$[-4:-4]|[2, "a", 4, 5, 100, "nice"]|[]|[]|:white_check_mark:|
|$[-4:-3]|[2, "a", 4, 5, 100, "nice"]|[4]|[4]|:white_check_mark:|
|$[-4:1]|[2, "a", 4, 5, 100, "nice"]|[]|[]|:white_check_mark:|
|$[-4:2]|[2, "a", 4, 5, 100, "nice"]|[]|[]|:white_check_mark:|
|$[-4:3]|[2, "a", 4, 5, 100, "nice"]|[4]|[4]|:white_check_mark:|
|$[3:0:-2]|["first", "second", "third", "forth", "fifth"]|[]|[]|:white_check_mark:|
|$[7:3:-1]|["first", "second", "third", "forth", "fifth"]|[]|[]|:white_check_mark:|
|$[0:3:-2]|["first", "second", "third", "forth", "fifth"]|none|[third first]|:question:|
|$[::-2]|["first", "second", "third", "forth", "fifth"]|none|[fifth third first]|:question:|
|$[1:]|["first", "second", "third", "forth", "fifth"]|[second third forth fifth]|[second third forth fifth]|:white_check_mark:|
|$[3::-1]|["first", "second", "third", "forth", "fifth"]|none|[fifth forth]|:question:|
|$[:2]|["first", "second", "third", "forth", "fifth"]|[first second]|[first second]|:white_check_mark:|
|$[:]|["first","second"]|[first second]|[first second]|:white_check_mark:|
|$[:]|{":": 42, "more": "string"}|[]|[42 string]|:no_entry:|
|$[::]|["first","second"]|[first second]|[first second]|:white_check_mark:|
|$[:2:-1]|["first", "second", "third", "forth", "fifth"]|none|[second first]|:question:|
|$[3:-4]|[2, "a", 4, 5, 100, "nice"]|[]|[]|:white_check_mark:|
|$[3:-3]|[2, "a", 4, 5, 100, "nice"]|[]|[]|:white_check_mark:|
|$[3:-2]|[2, "a", 4, 5, 100, "nice"]|[5]|[5]|:white_check_mark:|
|$[2:1]|["first", "second", "third", "forth"]|[]|[]|:white_check_mark:|
|$[0:0]|["first", "second"]|[]|[]|:white_check_mark:|
|$[0:1]|["first", "second"]|[first]|[first]|:white_check_mark:|
|$[-1:]|["first", "second", "third"]|[third]|[third]|:white_check_mark:|
|$[-2:]|["first", "second", "third"]|[second third]|[second third]|:white_check_mark:|
|$[-4:]|["first", "second", "third"]|[first second third]|[first second third]|:white_check_mark:|
|$[0:3:2]|["first", "second", "third", "forth", "fifth"]|[first third]|[first third]|:white_check_mark:|
|$[0:3:0]|["first", "second", "third", "forth", "fifth"]|none|<nil>|:question:|
|$[0:3:1]|["first", "second", "third", "forth", "fifth"]|[first second third]|[first second third]|:white_check_mark:|
|$[010:024:010]|[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]|[10 20]|[10 20]|:white_check_mark:|
|$[0:4:2]|["first", "second", "third", "forth", "fifth"]|[first third]|[first third]|:white_check_mark:|
|$[1:3:]|["first", "second", "third", "forth", "fifth"]|[second third]|[second third]|:white_check_mark:|
|$[::2]|["first", "second", "third", "forth", "fifth"]|[first third fifth]|[first third fifth]|:white_check_mark:|