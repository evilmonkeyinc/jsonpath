package jsonpath

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/script"
	"github.com/evilmonkeyinc/jsonpath/script/standard"
	"github.com/evilmonkeyinc/jsonpath/token"
)

// Compile will compile the JSONPath selector
func Compile(selector string) (*Selector, error) {
	engine := new(standard.ScriptEngine)

	if selector == "$[?(@.key<3),?(@.key>6)]" {
		// TODO
		selector = "$[?(@.key<3),?(@.key>6)]"
	}

	jsonPath := &Selector{}
	if err := jsonPath.compile(selector, engine); err != nil {
		return nil, getInvalidJSONPathSelectorWithReason(selector, err)
	}

	return jsonPath, nil
}

// Query will return the result of the JSONPath selector applied against the specified JSON data.
func Query(selector string, jsonData interface{}) (interface{}, error) {
	jsonPath, err := Compile(selector)
	if err != nil {
		return nil, getInvalidJSONPathSelectorWithReason(selector, err)
	}
	return jsonPath.Query(jsonData)
}

// QueryString will return the result of the JSONPath selector applied against the specified JSON data.
func QueryString(selector string, jsonData string) (interface{}, error) {

	jsonPath, err := Compile(selector)
	if err != nil {
		return nil, getInvalidJSONPathSelectorWithReason(selector, err)
	}
	return jsonPath.QueryString(jsonData)
}

// Selector represents a compiled JSONPath selector
// and exposes functions to query JSON data and objects.
type Selector struct {
	Options  *option.QueryOptions
	engine   script.Engine
	tokens   []token.Token
	selector string
}

// String returns the compiled selector string representation
func (query *Selector) String() string {
	jsonPath := ""
	for _, token := range query.tokens {
		jsonPath += fmt.Sprintf("%s", token)
	}
	return jsonPath
}

func (query *Selector) compile(selector string, engine script.Engine) error {
	query.engine = engine
	query.selector = selector

	tokenStrings, err := token.Tokenize(selector)
	if err != nil {
		return err
	}

	tokens := make([]token.Token, len(tokenStrings))
	for idx, tokenString := range tokenStrings {
		token, err := token.Parse(tokenString, query.engine, query.Options)
		if err != nil {
			return err
		}
		tokens[idx] = token
	}
	query.tokens = tokens

	return nil
}

// Query will return the result of the JSONPath query applied against the specified JSON data.
func (query *Selector) Query(root interface{}) (interface{}, error) {
	if len(query.tokens) == 0 {
		return nil, getInvalidJSONPathSelector(query.selector)
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
func (query *Selector) QueryString(jsonData string) (interface{}, error) {
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
