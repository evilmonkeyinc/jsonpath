package jsonpath

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/token"
)

// Compile compile the JSON path query
func Compile(queryPath string) (*JSONPath, error) {
	jsonPath := &JSONPath{}
	if err := jsonPath.compile(queryPath); err != nil {
		return nil, getInvalidJSONPathQueryWithReason(queryPath, err)
	}

	return jsonPath, nil
}

// Query will return the result of the JSONPath query applied against the specified JSON data.
func Query(queryPath string, jsonData interface{}) (interface{}, error) {
	jsonPath, err := Compile(queryPath)
	if err != nil {
		return nil, getInvalidJSONPathQueryWithReason(queryPath, err)
	}
	return jsonPath.Query(jsonData)
}

// QueryString will return the result of the JSONPath query applied against the specified JSON data.
func QueryString(queryPath string, jsonData string) (interface{}, error) {
	jsonPath, err := Compile(queryPath)
	if err != nil {
		return nil, getInvalidJSONPathQueryWithReason(queryPath, err)
	}
	return jsonPath.QueryString(jsonData)
}

// JSONPath represents a compiled JSONPath query
// and exposes functions to query JSON data and objects.
type JSONPath struct {
	Options     *option.QueryOptions
	queryString string
	tokens      []token.Token
}

// String returns the compiled query string representation
func (query *JSONPath) String() string {
	jsonPath := ""
	for _, token := range query.tokens {
		jsonPath += fmt.Sprintf("%s", token)
	}
	return jsonPath
}

func (query *JSONPath) compile(queryString string) error {
	query.queryString = queryString

	tokenStrings, err := token.Tokenize(queryString)
	if err != nil {
		return err
	}

	tokens := make([]token.Token, len(tokenStrings))
	for idx, tokenString := range tokenStrings {
		token, err := token.Parse(tokenString, query.Options)
		if err != nil {
			return err
		}
		tokens[idx] = token
	}
	query.tokens = tokens

	return nil
}

// Query will return the result of the JSONPath query applied against the specified JSON data.
func (query *JSONPath) Query(root interface{}) (interface{}, error) {
	if len(query.tokens) == 0 {
		return nil, getInvalidJSONPathQuery(query.queryString)
	}

	tokens := make([]token.Token, 0)
	if len(query.tokens) > 1 {
		tokens = query.tokens[1:]
	}

	found, err := query.tokens[0].Apply(root, root, tokens)
	if err != nil {
		return nil, err
	}
	return found, nil
}

// QueryString will return the result of the JSONPath query applied against the specified JSON data.
func (query *JSONPath) QueryString(jsonData string) (interface{}, error) {
	jsonData = strings.TrimSpace(jsonData)
	if jsonData == "" {
		return nil, getInvalidJSONData(errDataIsUnexpectedTypeOrNil)
	}

	var root interface{}

	if strings.HasPrefix(jsonData, "{") && strings.HasSuffix(jsonData, "}") {
		// object
		root = make(map[string]interface{})
		if err := json.Unmarshal([]byte(jsonData), &root); err != nil {
			return nil, getInvalidJSONData(err)
		}
	} else if strings.HasPrefix(jsonData, "[") && strings.HasSuffix(jsonData, "]") {
		// array
		root = make([]interface{}, 0)
		if err := json.Unmarshal([]byte(jsonData), &root); err != nil {
			return nil, getInvalidJSONData(err)
		}
	} else if len(jsonData) > 2 && strings.HasPrefix(jsonData, "\"") && strings.HasPrefix(jsonData, "\"") {
		// string
		root = jsonData[1 : len(jsonData)-1]
	} else if strings.ToLower(jsonData) == "true" {
		// bool true
		root = true
	} else if strings.ToLower(jsonData) == "false" {
		// bool false
		root = false
	} else if val, err := strconv.ParseInt(jsonData, 10, 64); err == nil {
		// integer
		root = val
	} else if val, err := strconv.ParseFloat(jsonData, 64); err == nil {
		// float
		root = val
	} else {
		return nil, getInvalidJSONData(errDataIsUnexpectedTypeOrNil)
	}

	return query.Query(root)
}
