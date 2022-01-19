package test

import "testing"

func Test_Filter(t *testing.T) {

	tests := []testData{
		{
			query:         `$[?(@.key)]`,
			data:          `{"key": 42, "another": {"key": 1}}`,
			expected:      []interface{}{map[string]interface{}{"key": float64(1)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$..*[?(@.id>2)]`,
			data:  `[ { "complext": { "one": [ { "name": "first", "id": 1 }, { "name": "next", "id": 2 }, { "name": "another", "id": 3 }, { "name": "more", "id": 4 } ], "more": { "name": "next to last", "id": 5 } } }, { "name": "last", "id": 6 } ]`,
			expected: []interface{}{
				[]interface{}{},
				[]interface{}{},
				[]interface{}{map[string]interface{}{"id": float64(5), "name": "next to last"}},
				[]interface{}{map[string]interface{}{"id": float64(3), "name": "another"}, map[string]interface{}{"id": float64(4), "name": "more"}},
				[]interface{}{},
				[]interface{}{},
				[]interface{}{},
				[]interface{}{},
				[]interface{}{},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$..[?(@.id==2)]`,
			data:  `{"id": 2, "more": [{"id": 2}, {"more": {"id": 2}}, {"id": {"id": 2}}, [{"id": 2}]]}`,
			expected: []interface{}{
				map[string]interface{}{"id": float64(2)},
				map[string]interface{}{"id": float64(2)},
				map[string]interface{}{"id": float64(2)},
				map[string]interface{}{"id": float64(2)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(@.key+50==100)]`,
			data:  `[{"key": 60}, {"key": 50}, {"key": 10}, {"key": -50}, {"key+50": 100}]`,
			expected: []interface{}{
				map[string]interface{}{"key": float64(50)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key>42 && @.key<44)]`,
			data:          `[ {"key": 42}, {"key": 43}, {"key": 44} ]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(43)}},
			consensus:     []interface{}{map[string]interface{}{"key": float64(43)}},
			expectedError: "",
		},
		{
			query:         `$[?(@.key>0 && false)]`,
			data:          `[ {"key": 1}, {"key": 3}, {"key": "nice"}, {"key": true}, {"key": null}, {"key": false}, {"key": {}}, {"key": []}, {"key": -1}, {"key": 0}, {"key": ""} ]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key>0 && true)]`,
			data:          `[ {"key": 1}, {"key": 3}, {"key": "nice"}, {"key": true}, {"key": null}, {"key": false}, {"key": {}}, {"key": []}, {"key": -1}, {"key": 0}, {"key": ""} ]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(1)}, map[string]interface{}{"key": float64(3)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key>43 || @.key<43)]`,
			data:          `[ {"key": 42}, {"key": 43}, {"key": 44} ]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(42)}, map[string]interface{}{"key": float64(44)}},
			consensus:     []interface{}{map[string]interface{}{"key": float64(42)}, map[string]interface{}{"key": float64(44)}},
			expectedError: "",
		},
		{
			query:         `$[?(@.key>0 || false)]`,
			data:          `[ {"key": 1}, {"key": 3}, {"key": "nice"}, {"key": true}, {"key": null}, {"key": false}, {"key": {}}, {"key": []}, {"key": -1}, {"key": 0}, {"key": ""} ]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(1)}, map[string]interface{}{"key": float64(3)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key>0 || true)]`,
			data:          `[ {"key": 1}, {"key": 3}, {"key": "nice"}, {"key": true}, {"key": null}, {"key": false}, {"key": {}}, {"key": []}, {"key": -1}, {"key": 0}, {"key": ""} ]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(1)}, map[string]interface{}{"key": float64(3)}, map[string]interface{}{"key": float64(-1)}, map[string]interface{}{"key": float64(0)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@['key']==42)]`,
			data:          `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"some": "value"} ]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(42)}},
			consensus:     []interface{}{map[string]interface{}{"key": float64(42)}},
			expectedError: "",
		},
		{
			query:         `$[?(@['@key']==42)]`,
			data:          `[ {"@key": 0}, {"@key": 42}, {"key": 42}, {"@key": 43}, {"some": "value"} ]`,
			expected:      []interface{}{map[string]interface{}{"@key": float64(42)}},
			consensus:     []interface{}{map[string]interface{}{"@key": float64(42)}},
			expectedError: "",
		},
		{
			query:         `$[?(@[-1]==2)]`,
			data:          `[[2, 3], ["a"], [0, 2], [2]]`,
			expected:      []interface{}{[]interface{}{float64(0), float64(2)}, []interface{}{float64(2)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@[1]=='b')]`,
			data:          `[["a", "b"], ["x", "y"]]`,
			expected:      []interface{}{[]interface{}{"a", "b"}},
			consensus:     []interface{}{[]interface{}{"a", "b"}},
			expectedError: "",
		},
		{
			query:         `$[?(@[1]=='b')]`,
			data:          `{"1": ["a", "b"], "2": ["x", "y"]}`,
			expected:      []interface{}{[]interface{}{"a", "b"}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@)]`,
			data:          `[ "some value", null, "value", 0, 1, -1, "", [], {}, false, true ]`,
			expected:      []interface{}{"some value", "value", float64(0), float64(1), float64(-1), true},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.a && (@.b || @.c))]`,
			data:          `[ { "a": true }, { "a": true, "b": true }, { "a": true, "b": true, "c": true }, { "b": true, "c": true }, { "a": true, "c": true }, { "c": true }, { "b": true } ]`,
			expected:      []interface{}{}, // TODO : this should be different
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `[?(@.a && @.b || @.c)]`,
			data:          `[ { "a": true, "b": true }, { "a": true, "b": true, "c": true }, { "b": true, "c": true }, { "a": true, "c": true }, { "a": true }, { "b": true }, { "c": true }, { "d": true }, {} ]`,
			expected:      nil,
			consensus:     consensusNone,
			expectedError: "invalid JSONPath query '[?(@.a && @.b || @.c)]' unexpected token '[' at index 0",
		},
		{
			query:         `$[?(@.key/10==5)]`,
			data:          `[{"key": 60}, {"key": 50}, {"key": 10}, {"key": -50}, {"key/10": 5}]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(50)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key-dash == 'value')]`,
			data:          `[ { "key-dash": "value" } ]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.2 == 'second')]`,
			data:          `[{"a": "first", "2": "second", "b": "third"}]`,
			expected:      []interface{}{map[string]interface{}{"2": "second", "a": "first", "b": "third"}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.2 == 'third')]`,
			data:          `[["first", "second", "third", "forth", "fifth"]] `,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?()]`,
			data:          `[1, {"key": 42}, "value", null]`,
			expected:      nil,
			consensus:     nil,
			expectedError: "invalid expression. is empty",
		},
		{
			query:         `$[?(@.key==42)]`,
			data:          `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(42)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@==42)]`,
			data:          `[ 0, 42, -1, 41, 43, 42.0001, 41.9999, null, 100 ]`,
			expected:      []interface{}{}, // TODO : doesn't match as float == int, needs fixed
			consensus:     []interface{}{float64(42)},
			expectedError: "",
		},
		{
			query:         `$[?(@.key==43)]`,
			data:          `[{"key": 42}]`,
			expected:      []interface{}{},
			consensus:     []interface{}{},
			expectedError: "",
		},
		{
			query:         `$[?(@.key==42)]`,
			data:          `{ "a": {"key": 0}, "b": {"key": 42}, "c": {"key": -1}, "d": {"key": 41}, "e": {"key": 43}, "f": {"key": 42.0001}, "g": {"key": 41.9999}, "h": {"key": 100}, "i": {"some": "value"} }`,
			expected:      []interface{}{map[string]interface{}{"key": float64(42)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.id==2)]`,
			data:          `{"id": 2}`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.d==["v1","v2"])]`,
			data:          `[ { "d": [ "v1", "v2" ] }, { "d": [ "a", "b" ] }, { "d": "v1" }, { "d": "v2" }, { "d": {} }, { "d": [] }, { "d": null }, { "d": -1 }, { "d": 0 }, { "d": 1 }, { "d": "['v1','v2']" }, { "d": "['v1', 'v2']" }, { "d": "v1,v2" }, { "d": "[\"v1\", \"v2\"]" }, { "d": "[\"v1\",\"v2\"]" } ]`,
			expected:      []interface{}{}, // TODO : doesn't know [] is array
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@[0:1]==[1])]`,
			data:          `[[1, 2, 3], [1], [2, 3], 1, 2]`,
			expected:      []interface{}{}, // TODO : doesn't know [] is array
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.*==[1,2])]`,
			data:          `[[1,2], [2,3], [1], [2], [1, 2, 3], 1, 2, 3]`,
			expected:      []interface{}{}, // TODO : need to parse [1,2] as array
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.d==['v1','v2'])]`,
			data:          `[ { "d": [ "v1", "v2" ] }, { "d": [ "a", "b" ] }, { "d": "v1" }, { "d": "v2" }, { "d": {} }, { "d": [] }, { "d": null }, { "d": -1 }, { "d": 0 }, { "d": 1 }, { "d": "['v1','v2']" }, { "d": "['v1', 'v2']" }, { "d": "v1,v2" }, { "d": "[\"v1\", \"v2\"]" }, { "d": "[\"v1\",\"v2\"]" } ]`,
			expected:      []interface{}{}, // TODO : need to parse ['v1','v2'] as array
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?((@.key<44)==false)]`,
			data:          `[{"key": 42}, {"key": 43}, {"key": 44}]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(44)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key==false)]`,
			data:          `[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`,
			expected:      []interface{}{map[string]interface{}{"key": false}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key==null)]`, // TODO : nil/null support
			data:          `[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@[0:1]==1)]`, // TODO array in expression support
			data:          `[[1, 2, 3], [1], [2, 3], 1, 2]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@[*]==2)]`, // TODO array in expression support
			data:          `[[1,2], [2,3], [1], [2], [1, 2, 3], 1, 2, 3]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.*==2)]`,
			data:          `[[1,2], [2,3], [1], [2], [1, 2, 3], 1, 2, 3]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key==-0.123e2)]`,
			data:          `[{"key": -12.3}, {"key": -0.123}, {"key": -12}, {"key": 12.3}, {"key": 2}, {"key": "-0.123e2"}]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(-12.3)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key==010)]`, // TODO : is this expected
			data:          `[{"key": "010"}, {"key": "10"}, {"key": 10}, {"key": 0}, {"key": 8}]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(8)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.d=={"k":"v"})]`, // TODO: will not parse as map
			data:          `[ { "d": { "k": "v" } }, { "d": { "a": "b" } }, { "d": "k" }, { "d": "v" }, { "d": {} }, { "d": [] }, { "d": null }, { "d": -1 }, { "d": 0 }, { "d": 1 }, { "d": "[object Object]" }, { "d": "{\"k\": \"v\"}" }, { "d": "{\"k\":\"v\"}" }, "v" ]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key=="value")]`,
			data:          `[ {"key": "some"}, {"key": "value"}, {"key": null}, {"key": 0}, {"key": 1}, {"key": -1}, {"key": ""}, {"key": {}}, {"key": []}, {"key": "valuemore"}, {"key": "morevalue"}, {"key": ["value"]}, {"key": {"some": "value"}}, {"key": {"key": "value"}}, {"some": "value"} ]`,
			expected:      []interface{}{map[string]interface{}{"key": "value"}},
			consensus:     []interface{}{map[string]interface{}{"key": "value"}},
			expectedError: "",
		},
		{
			query:         `$[?(@.key=="Motörhead")]`,
			data:          `[ {"key": "something"}, {"key": "Mot\u00f6rhead"}, {"key": "mot\u00f6rhead"}, {"key": "Motorhead"}, {"key": "Motoo\u0308rhead"}, {"key": "motoo\u0308rhead"} ]`,
			expected:      []interface{}{map[string]interface{}{"key": "Mot\u00f6rhead"}},
			consensus:     []interface{}{map[string]interface{}{"key": "Mot\u00f6rhead"}},
			expectedError: "",
		},
		{
			query:         `$[?(@.key=="hi@example.com")]`, // TODO: double quotes in filters
			data:          `[ {"key": "some"}, {"key": "value"}, {"key": "hi@example.com"} ]`,
			expected:      []interface{}{},
			consensus:     []interface{}{map[string]interface{}{"key": "hi@example.com"}},
			expectedError: "",
		},
		{
			query:         `$[?(@.key=="some.value")]`,
			data:          `[ {"key": "some"}, {"key": "value"}, {"key": "some.value"} ]`,
			expected:      []interface{}{map[string]interface{}{"key": "some.value"}},
			consensus:     []interface{}{map[string]interface{}{"key": "some.value"}},
			expectedError: "",
		},
		{
			query:         `$[?(@.key=='value')]`,
			data:          `[ {"key": "some"}, {"key": "value"} ]`,
			expected:      []interface{}{map[string]interface{}{"key": "value"}},
			consensus:     []interface{}{map[string]interface{}{"key": "value"}},
			expectedError: "",
		},
		{
			query:         `$[?(@.key=="Mot\u00f6rhead")]`,
			data:          `[ {"key": "something"}, {"key": "Mot\u00f6rhead"}, {"key": "mot\u00f6rhead"}, {"key": "Motorhead"}, {"key": "Motoo\u0308rhead"}, {"key": "motoo\u0308rhead"} ]`,
			expected:      []interface{}{map[string]interface{}{"key": "Motörhead"}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key==true)]`,
			data:          `[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`,
			expected:      []interface{}{map[string]interface{}{"key": true}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(@.key1==@.key2)]`,
			data:  `[ {"key1": 10, "key2": 10}, {"key1": 42, "key2": 50}, {"key1": 10}, {"key2": 10}, {}, {"key1": null, "key2": null}, {"key1": null}, {"key2": null}, {"key1": 0, "key2": 0}, {"key1": 0}, {"key2": 0}, {"key1": -1, "key2": -1}, {"key1": "", "key2": ""}, {"key1": false, "key2": false}, {"key1": false}, {"key2": false}, {"key1": true, "key2": true}, {"key1": [], "key2": []}, {"key1": {}, "key2": {}}, {"key1": {"a": 1, "b": 2}, "key2": {"b": 2, "a": 1}} ]`,
			expected: []interface{}{
				map[string]interface{}{"key1": float64(10), "key2": float64(10)},
				map[string]interface{}{"key1": float64(0), "key2": float64(0)},
				map[string]interface{}{"key1": float64(-1), "key2": float64(-1)},
				map[string]interface{}{"key1": "", "key2": ""},
				map[string]interface{}{"key1": false, "key2": false},
				map[string]interface{}{"key1": true, "key2": true},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$.items[?(@.key==$.value)]`,
			data:          `{"value": 42, "items": [{"key": 10}, {"key": 42}, {"key": 50}]}`,
			expected:      []interface{}{map[string]interface{}{"key": float64(42)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(@.key>42)]`,
			data:  `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`,
			expected: []interface{}{
				map[string]interface{}{"key": float64(43)},
				map[string]interface{}{"key": float64(42.0001)},
				map[string]interface{}{"key": float64(100)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(@.key>=42)]`,
			data:  `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`,
			expected: []interface{}{
				map[string]interface{}{"key": float64(42)},
				map[string]interface{}{"key": float64(43)},
				map[string]interface{}{"key": float64(42.0001)},
				map[string]interface{}{"key": float64(100)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.d in [2, 3])]`, // TODO : in keywork will not work
			data:          `[{"d": 1}, {"d": 2}, {"d": 1}, {"d": 3}, {"d": 4}]`,
			expected:      []interface{}{},
			consensus:     nil,
			expectedError: "",
		},
		{
			query:         `$[?(2 in @.d)]`, // TODO : in keywork will not work
			data:          `[{"d": [1, 2, 3]}, {"d": [2]}, {"d": [1]}, {"d": [3, 4]}, {"d": [4, 2]}]`,
			expected:      []interface{}{},
			consensus:     nil,
			expectedError: "",
		},
		{
			query: `$[?(@.key<42)]`,
			data:  `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`,
			expected: []interface{}{
				map[string]interface{}{"key": float64(0)},
				map[string]interface{}{"key": float64(-1)},
				map[string]interface{}{"key": float64(41)},
				map[string]interface{}{"key": float64(41.9999)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(@.key<=42)]`,
			data:  `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`,
			expected: []interface{}{
				map[string]interface{}{"key": float64(0)},
				map[string]interface{}{"key": float64(42)},
				map[string]interface{}{"key": float64(-1)},
				map[string]interface{}{"key": float64(41)},
				map[string]interface{}{"key": float64(41.9999)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key*2==100)]`,
			data:          `[{"key": 60}, {"key": 50}, {"key": 10}, {"key": -50}, {"key*2": 100}]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(50)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(!(@.key==42))]`,
			data:  `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`,
			expected: []interface{}{
				map[string]interface{}{"key": float64(0)},
				map[string]interface{}{"key": float64(-1)},
				map[string]interface{}{"key": float64(41)},
				map[string]interface{}{"key": float64(43)},
				map[string]interface{}{"key": float64(42.0001)},
				map[string]interface{}{"key": float64(41.9999)},
				map[string]interface{}{"key": float64(100)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(!(@.key<42))]`,
			data:  `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "43"}, {"key": "42"}, {"key": "41"}, {"key": "value"}, {"some": "value"} ]`,
			expected: []interface{}{
				map[string]interface{}{"key": float64(42)},
				map[string]interface{}{"key": float64(43)},
				map[string]interface{}{"key": float64(42.0001)},
				map[string]interface{}{"key": float64(100)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(!@.key)]`,
			data:          `[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`,
			expected:      []interface{}{map[string]interface{}{"key": false}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(@.key!=42)]`,
			data:  `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`,
			expected: []interface{}{
				map[string]interface{}{"key": float64(0)},
				map[string]interface{}{"key": float64(-1)},
				map[string]interface{}{"key": float64(1)},
				map[string]interface{}{"key": float64(41)},
				map[string]interface{}{"key": float64(43)},
				map[string]interface{}{"key": float64(42.0001)},
				map[string]interface{}{"key": float64(41.9999)},
				map[string]interface{}{"key": float64(100)},
				map[string]interface{}{"key": float64(420)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[*].bookmarks[?(@.page == 45)]^^^`,
			data:          `[ { "title": "Sayings of the Century", "bookmarks": [{ "page": 40 }] }, { "title": "Sword of Honour", "bookmarks": [ { "page": 35 }, { "page": 45 } ] }, { "title": "Moby Dick", "bookmarks": [ { "page": 3035 }, { "page": 45 } ] } ]`,
			expected:      []interface{}{[]interface{}{}, []interface{}{}, []interface{}{}},
			consensus:     nil,
			expectedError: "",
		},
		{
			query:         `$[?(@.name=~/hello.*/)]`, // TODO : need regex support
			data:          `[ {"name": "hullo world"}, {"name": "hello world"}, {"name": "yes hello world"}, {"name": "HELLO WORLD"}, {"name": "good bye"} ]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.name=~/@.pattern/)]`,
			data:          `[ {"name": "hullo world"}, {"name": "hello world"}, {"name": "yes hello world"}, {"name": "HELLO WORLD"}, {"name": "good bye"}, {"pattern": "hello.*"} ]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@[*]>=4)]`,
			data:          `[[1,2],[3,4],[5,6]]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$.x[?(@[*]>=$.y[*])]`,
			data:          `{"x":[[1,2],[3,4],[5,6]],"y":[3,4,5]}`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key=42)]`,
			data:          `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`,
			expected:      []interface{}{},
			consensus:     nil,
			expectedError: "",
		},
		{
			query:         `$[?(@.a[?(@.price>10)])]`,
			data:          `[ { "a": [{"price": 1}, {"price": 3}] }, { "a": [{"price": 11}] }, { "a": [{"price": 8}, {"price": 12}, {"price": 3}] }, { "a": [] } ]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.address.city=='Berlin')]`,
			data:          `[ { "address": { "city": "Berlin" } }, { "address": { "city": "London" } } ]`,
			expected:      []interface{}{map[string]interface{}{"address": map[string]interface{}{"city": "Berlin"}}},
			consensus:     []interface{}{map[string]interface{}{"address": map[string]interface{}{"city": "Berlin"}}},
			expectedError: "",
		},
		{
			query:         `$[?(@.key-50==-100)]`,
			data:          `[{"key": 60}, {"key": 50}, {"key": 10}, {"key": -50}, {"key-50": -100}]`,
			expected:      []interface{}{map[string]interface{}{"key": float64(-50)}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(1==1)]`,
			data:          `[1, 3, "nice", true, null, false, {}, [], -1, 0, ""]`,
			expected:      []interface{}{float64(1), float64(3), "nice", true, interface{}(nil), false, map[string]interface{}{}, []interface{}{}, float64(-1), float64(0), ""},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@.key===42)]`,
			data:          `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$[?(@.key)]`,
			data:  `[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`,
			expected: []interface{}{
				map[string]interface{}{"key": true},
				map[string]interface{}{"key": "value"},
				map[string]interface{}{"key": float64(0)},
				map[string]interface{}{"key": float64(1)},
				map[string]interface{}{"key": float64(-1)},
				map[string]interface{}{"key": float64(42)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$.*[?(@.key)]`,
			data:          `[ { "some": "some value" }, { "key": "value" } ]`,
			expected:      []interface{}{[]interface{}{}, []interface{}{}},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query: `$..[?(@.id)]`,
			data:  `{"id": 2, "more": [{"id": 2}, {"more": {"id": 2}}, {"id": {"id": 2}}, [{"id": 2}]]}`,
			expected: []interface{}{
				map[string]interface{}{"id": float64(2)},
				map[string]interface{}{"id": float64(2)},
				map[string]interface{}{"id": float64(2)},
				map[string]interface{}{"id": float64(2)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(false)]`,
			data:          `[1, 3, "nice", true, null, false, {}, [], -1, 0, ""]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(@..child)]`,
			data:          `[{"key": [{"child": 1}, {"child": 2}]}, {"key": [{"child": 2}]}, {"key": [{}]}, {"key": [{"something": 42}]}, {}]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(null)]`, // TODO: null should be parsed to nil
			data:          `[1, 3, "nice", true, null, false, {}, [], -1, 0, ""]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?(true)]`,
			data:          `[1, 3, "nice", true, null, false, {}, [], -1, 0, ""]`,
			expected:      []interface{}{float64(1), float64(3), "nice", true, interface{}(nil), false, map[string]interface{}{}, []interface{}{}, float64(-1), float64(0), ""},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[?@.key==42]`,
			data:          `[ {"key": 0}, {"key": 42}, {"key": -1}, {"key": 1}, {"key": 41}, {"key": 43}, {"key": 42.0001}, {"key": 41.9999}, {"key": 100}, {"key": "some"}, {"key": "42"}, {"key": null}, {"key": 420}, {"key": ""}, {"key": {}}, {"key": []}, {"key": [42]}, {"key": {"key": 42}}, {"key": {"some": 42}}, {"some": "value"} ]`,
			expected:      nil,
			consensus:     consensusNone,
			expectedError: "invalid JSONPath query '$[?@.key==42]' invalid token. '[?@.key==42]' does not match any token format",
		},
		{
			query: `$[?(@.key)]`,
			data:  `[ { "some": "some value" }, { "key": true }, { "key": false }, { "key": null }, { "key": "value" }, { "key": "" }, { "key": 0 }, { "key": 1 }, { "key": -1 }, { "key": 42 }, { "key": {} }, { "key": [] } ]`,
			expected: []interface{}{
				map[string]interface{}{"key": true},
				map[string]interface{}{"key": "value"},
				map[string]interface{}{"key": float64(0)},
				map[string]interface{}{"key": float64(1)},
				map[string]interface{}{"key": float64(-1)},
				map[string]interface{}{"key": float64(42)},
			},
			consensus:     consensusNone,
			expectedError: "",
		},
	}

	batchTest(t, tests)
	// printConsensusMatrix(tests)
}