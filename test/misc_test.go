package test

import "testing"

func Test_Misc(t *testing.T) {
	batchTest(t, miscTests)
}

var miscTests []testData = []testData{
	{
		query:         "", // empty
		data:          `{"a": 42, "": 21}`,
		expected:      nil,
		consensus:     nil,
		expectedError: "invalid JSONPath query '' unexpected token '' at index 0",
	},
	{
		query:         `$.data.sum()`,
		data:          `{"data": [1,2,3,4]}`,
		expected:      nil,
		consensus:     consensusNone,
		expectedError: "key: invalid token target. expected [map] got [slice]",
	},
	{
		query:         `$(key,more)`,
		data:          `{"key": 1, "some": 2, "more": 3}`,
		expected:      nil,
		consensus:     nil,
		expectedError: "invalid JSONPath query '$(key,more)' unexpected token '(' at index 1",
	},
	{
		query:         `$..`,
		data:          `[{"a": {"b": "c"}}, [0, 1]]`,
		expected:      []interface{}{[]interface{}{map[string]interface{}{"a": map[string]interface{}{"b": "c"}}, []interface{}{float64(0), float64(1)}}, map[string]interface{}{"a": map[string]interface{}{"b": "c"}}, map[string]interface{}{"b": "c"}, "c", []interface{}{float64(0), float64(1)}, float64(0), float64(1)},
		consensus:     consensusNone,
		expectedError: "",
	},
	{
		query:         `$.key..`,
		data:          `{"some key": "value", "key": {"complex": "string", "primitives": [0, 1]}}`,
		expected:      []interface{}{map[string]interface{}{"complex": "string", "primitives": []interface{}{float64(0), float64(1)}}, []interface{}{float64(0), float64(1)}, float64(0), float64(1), "string"},
		consensus:     consensusNone,
		expectedError: "",
	},
	{
		query:         `$`,
		data:          `{ "key": "value", "another key": { "complex": [ "a", 1 ] } }`,
		expected:      map[string]interface{}{"key": "value", "another key": map[string]interface{}{"complex": []interface{}{"a", float64(1)}}},
		consensus:     map[string]interface{}{"key": "value", "another key": map[string]interface{}{"complex": []interface{}{"a", float64(1)}}},
		expectedError: "",
	},
	{
		query:         `$`,
		data:          `42`,
		expected:      int64(42),
		consensus:     int64(42),
		expectedError: "",
	},
	{
		query:         `$`,
		data:          `false`,
		expected:      false,
		consensus:     false,
		expectedError: "",
	},
	{
		query:         `$`,
		data:          `true`,
		expected:      true,
		consensus:     true,
		expectedError: "",
	},
	{
		query:         `$[(@.length-1)]`,
		data:          `["first", "second", "third", "forth", "fifth"]`,
		expected:      "fifth",
		consensus:     nil,
		expectedError: "",
	},
}
