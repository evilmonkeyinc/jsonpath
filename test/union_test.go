package test

import "testing"

func Test_Union(t *testing.T) {
	batchTest(t, unionTests)
}

func Benchmark_Union(b *testing.B) {
	batchBenchmark(b, unionTests)
}

var unionTests []testData = []testData{
	{
		selector:      `$[0,1]`,
		data:          `["first", "second", "third"]`,
		expected:      []interface{}{"first", "second"},
		consensus:     []interface{}{"first", "second"},
		expectedError: "",
	},
	{
		selector:      `$[0,0]`,
		data:          `["a"]`,
		expected:      []interface{}{"a", "a"},
		consensus:     []interface{}{"a", "a"},
		expectedError: "",
	},
	{
		selector:      `$['a','a']`,
		data:          `{"a":1}`,
		expected:      []interface{}{float64(1), float64(1)},
		consensus:     []interface{}{float64(1), float64(1)},
		expectedError: "",
	},
	{
		selector:      `$[?(@.key<3),?(@.key>6)]`,
		data:          `[{"key": 1}, {"key": 8}, {"key": 3}, {"key": 10}, {"key": 7}, {"key": 2}, {"key": 6}, {"key": 4}]`,
		expected:      []interface{}{},
		consensus:     consensusNone,
		expectedError: "",
	},
	{
		selector:      `$['key','another']`,
		data:          `{ "key": "value", "another": "entry" }`,
		expected:      []interface{}{"value", "entry"},
		consensus:     []interface{}{"value", "entry"},
		expectedError: "",
	},
	{
		selector:      `$['missing','key']`,
		data:          `{ "key": "value", "another": "entry" }`,
		expected:      []interface{}{"value"},
		consensus:     []interface{}{"value"},
		expectedError: "",
	},
	{
		selector:      `$[:]['c','d']`,
		data:          `[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`,
		expected:      []interface{}{[]interface{}{"cc1", "dd1"}, []interface{}{"cc2", "dd2"}},
		consensus:     []interface{}{"cc1", "dd1", "cc2", "dd2"},
		expectedError: "",
	},
	{
		selector:      `$[0]['c','d']`,
		data:          `[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`,
		expected:      []interface{}{"cc1", "dd1"},
		consensus:     []interface{}{"cc1", "dd1"},
		expectedError: "",
	},
	{
		selector:      `$.*['c','d']`,
		data:          `[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`,
		expected:      []interface{}{[]interface{}{"cc1", "dd1"}, []interface{}{"cc2", "dd2"}},
		consensus:     []interface{}{"cc1", "dd1", "cc2", "dd2"},
		expectedError: "",
	},
	{
		selector:      `$..['c','d']`,
		data:          `[{"c":"cc1","d":"dd1","e":"ee1"}, {"c": "cc2", "child": {"d": "dd2"}}, {"c": "cc3"}, {"d": "dd4"}, {"child": {"c": "cc5"}}]`,
		expected:      []interface{}{"cc1", "dd1", "cc2", "dd2", "cc3", "dd4", "cc5"},
		consensus:     consensusNone,
		expectedError: "",
	},
	{
		selector:      `$[4,1]`,
		data:          `[1,2,3,4,5]`,
		expected:      []interface{}{float64(5), float64(2)},
		consensus:     []interface{}{float64(5), float64(2)},
		expectedError: "",
	},
	{
		selector:      `$.*[0,:5]`,
		data:          `{ "a": [ "string", null, true ], "b": [ false, "string", 5.4 ] }`,
		expected:      nil,
		consensus:     consensusNone,
		expectedError: "invalid JSONPath selector '$.*[0,:5]' invalid token. '[0,:5]' does not match any token format",
	},
	{
		selector:      `$[1:3,4]`,
		data:          `[1,2,3,4,5]`,
		expected:      nil,
		consensus:     consensusNone,
		expectedError: "invalid JSONPath selector '$[1:3,4]' invalid token. '[1:3,4]' does not match any token format",
	},
	{
		selector:      `$[ 0 , 1 ]`,
		data:          `["first", "second", "third"]`,
		expected:      []interface{}{"first", "second"},
		consensus:     []interface{}{"first", "second"},
		expectedError: "",
	},
	{
		selector:      `$[*,1]`,
		data:          `["first", "second", "third", "forth", "fifth"]`,
		expected:      nil,
		consensus:     nil,
		expectedError: "invalid JSONPath selector '$[*,1]' invalid token. '[*,1]' does not match any token format",
	},
}
