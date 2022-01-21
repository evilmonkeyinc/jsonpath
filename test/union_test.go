package test

import "testing"

func Test_Union(t *testing.T) {

	tests := []testData{
		{
			query:         `$[0,1]`,
			data:          `["first", "second", "third"]`,
			expected:      []interface{}{"first", "second"},
			consensus:     []interface{}{"first", "second"},
			expectedError: "",
		},
		{
			query:         `$[0,0]`,
			data:          `["a"]`,
			expected:      []interface{}{"a", "a"},
			consensus:     []interface{}{"a", "a"},
			expectedError: "",
		},
		{
			query:         `$['a','a']`,
			data:          `{"a":1}`,
			expected:      []interface{}{float64(1), float64(1)},
			consensus:     []interface{}{float64(1), float64(1)},
			expectedError: "",
		},
		{
			query:         `$[?(@.key<3),?(@.key>6)]`,
			data:          `[{"key": 1}, {"key": 8}, {"key": 3}, {"key": 10}, {"key": 7}, {"key": 2}, {"key": 6}, {"key": 4}]`,
			expected:      []interface{}{},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$['key','another']`,
			data:          `{ "key": "value", "another": "entry" }`,
			expected:      []interface{}{"value", "entry"},
			consensus:     []interface{}{"value", "entry"},
			expectedError: "",
		},
		{
			query:         `$['missing','key']`,
			data:          `{ "key": "value", "another": "entry" }`,
			expected:      []interface{}{"value"},
			consensus:     []interface{}{"value"},
			expectedError: "",
		},
		{
			query:         `$[:]['c','d']`,
			data:          `[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`,
			expected:      []interface{}{[]interface{}{"cc1", "dd1"}, []interface{}{"cc2", "dd2"}},
			consensus:     []interface{}{"cc1", "dd1", "cc2", "dd2"},
			expectedError: "",
		},
		{
			query:         `$[0]['c','d']`,
			data:          `[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`,
			expected:      []interface{}{"cc1", "dd1"},
			consensus:     []interface{}{"cc1", "dd1"},
			expectedError: "",
		},
		{
			query:         `$.*['c','d']`,
			data:          `[{"c":"cc1","d":"dd1","e":"ee1"},{"c":"cc2","d":"dd2","e":"ee2"}]`,
			expected:      []interface{}{[]interface{}{"cc1", "dd1"}, []interface{}{"cc2", "dd2"}},
			consensus:     []interface{}{"cc1", "dd1", "cc2", "dd2"},
			expectedError: "",
		},
		{
			query:         `$..['c','d']`,
			data:          `[{"c":"cc1","d":"dd1","e":"ee1"}, {"c": "cc2", "child": {"d": "dd2"}}, {"c": "cc3"}, {"d": "dd4"}, {"child": {"c": "cc5"}}]`,
			expected:      []interface{}{"cc1", "dd1", "cc2", "dd2", "cc3", "dd4", "cc5"},
			consensus:     consensusNone,
			expectedError: "",
		},
		{
			query:         `$[4,1]`,
			data:          `[1,2,3,4,5]`,
			expected:      []interface{}{float64(5), float64(2)},
			consensus:     []interface{}{float64(5), float64(2)},
			expectedError: "",
		},
		{
			query:         `$.*[0,:5]`,
			data:          `{ "a": [ "string", null, true ], "b": [ false, "string", 5.4 ] }`,
			expected:      nil,
			consensus:     consensusNone,
			expectedError: "invalid JSONPath query '$.*[0,:5]' invalid token. '[0,:5]' does not match any token format",
		},
		{
			query:         `$[1:3,4]`,
			data:          `[1,2,3,4,5]`,
			expected:      nil,
			consensus:     consensusNone,
			expectedError: "invalid JSONPath query '$[1:3,4]' invalid token. '[1:3,4]' does not match any token format",
		},
		{
			query:         `$[ 0 , 1 ]`,
			data:          `["first", "second", "third"]`,
			expected:      []interface{}{"first", "second"},
			consensus:     []interface{}{"first", "second"},
			expectedError: "",
		},
		{
			query:         `$[*,1]`,
			data:          `["first", "second", "third", "forth", "fifth"]`,
			expected:      nil,
			consensus:     nil,
			expectedError: "invalid JSONPath query '$[*,1]' invalid token. '[*,1]' does not match any token format",
		},
	}

	batchTest(t, tests)
	//printConsensusMatrix(tests)
}
