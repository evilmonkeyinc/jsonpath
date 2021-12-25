package jsonpath

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type tokenOperation string

// TODO: NOTE
// SCRIPT @ is current object/array
// FILTER @ is child element

const (
	tokenOperationRoot      tokenOperation = "$"
	tokenOperationCurrent   tokenOperation = "@"
	tokenOperationWildcard  tokenOperation = "*"
	tokenOperationRecursive tokenOperation = ".."
	tokenOperationKey       tokenOperation = "KEY"
	tokenOperationFilter    tokenOperation = "FILTER"
	tokenOperationRange     tokenOperation = "RANGE"
	tokenOperationScript    tokenOperation = "SCRIPT"
	tokenOperationIndex     tokenOperation = "INDEX"
	tokenOperationUnion     tokenOperation = "UNION"
)

type token struct {
	operation tokenOperation
	arguments []interface{}
}

// TODO: either here or in functions that call tokenize
// 1. validate scripts and filters are valid syntax (@. etc)
// 2. validate that subset commands [] are valid (must be scripts or numerical operations or single quotes)
func tokenize(query string) ([]string, error) {
	if query == "" {
		return nil, errors.ErrQueryNotSpecified
	}

	tokens := []string{}
	tokenString := ""

	for idx, rne := range query {
		tokenString += string(rne)

		if idx == 0 {
			if tokenString != "$" && tokenString != "@" {
				return nil, errors.ErrInvalidInitialToken
			}

			if next := query[1]; next != '.' && next != '[' {
				return nil, errors.ErrInvalidInitialToken
			}

			tokens = append(tokens, tokenString[:])
			tokenString = ""
			continue
		}

		if tokenString == "." {
			continue
		}

		if tokenString == ".." {
			// recursive operator
			tokens = append(tokens, tokenString[:])
			tokenString = ""
			continue
		}

		if rne == '[' {
			if tokenString == "[" || tokenString == ".[" {
				// open bracket and at start of token
				continue
			}
			// open bracket in middle of token, new subscript
			if strings.Count(tokenString, "[") > 1 {
				// this is not the only opening bracket, subscript in subscript
				continue
			} else {
				// subscript should be own token
				if tokenString[0] == '.' {
					tokenString = tokenString[1 : len(tokenString)-1]
				} else {
					tokenString = tokenString[:len(tokenString)-1]
				}

				tokens = append(tokens, tokenString[:])

				tokenString = "["
				continue
			}
		}

		if strings.Contains(tokenString, "[") {
			startCount := strings.Count(tokenString, "[")
			endCount := strings.Count(tokenString, "]")
			// if end bracket, as long as it's not been escaped
			if rne == ']' && startCount == endCount {
				if tokenString[0] == '.' {
					tokenString = tokenString[1:]
				} else {
					tokenString = tokenString[:]
				}

				tokens = append(tokens, tokenString[:])

				tokenString = ""
				continue
			}
		} else if rne == '.' {
			if tokenString[0] == '.' {
				tokenString = tokenString[1 : len(tokenString)-1]
			} else {
				tokenString = tokenString[:len(tokenString)-1]
			}

			tokens = append(tokens, tokenString[:])

			tokenString = "."
			continue
		}
	}

	// parse the last token
	if len(tokenString) > 0 {
		if tokenString[0] == '.' {
			tokenString = tokenString[1:]
		}

		tokens = append(tokens, tokenString[:])
	}

	return tokens, nil
}

func parseToken(tkn string) (*token, error) {

	isScript := func(token string) bool {
		return strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")")
	}

	isKey := func(token string) bool {
		return strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'")
	}

	tkn = strings.TrimSpace(tkn)
	if tkn == "" {
		return nil, errors.GetInvalidTokenError("token can not be empty")
	}

	if tkn == "$" {
		return &token{operation: tokenOperationRoot}, nil
	}
	if tkn == "@" {
		return &token{operation: tokenOperationCurrent}, nil
	}
	if tkn == "*" {
		return &token{operation: tokenOperationWildcard}, nil
	}
	if tkn == ".." {
		return &token{operation: tokenOperationRecursive}, nil
	}

	if !strings.HasPrefix(tkn, "[") {

		if _, err := strconv.Atoi(tkn); err == nil {
			return nil, errors.GetInvalidTokenError("index specified as key")
		}

		return &token{
			operation: tokenOperationKey,
			arguments: []interface{}{fmt.Sprintf("'%s'", tkn)},
		}, nil

	}

	if !strings.HasSuffix(tkn, "]") {
		return nil, errors.GetInvalidTokenError("missing subscript close")
	}
	// subscript, or child operator

	/**
	if in brackets it is one of the following:
	wild: wildcard symbol, states that everything in array or map should be returned
	index: an integer, positive or negative, an index in an array
	key: a string quoted with 'single quotes', the key in a map
	filter: a filter query ?(x), used to filter an array or map based on child attributes
	query: a script (x), used to calculate an index or key based on the current attribute
	range: a range, two or three values seperated by : characters used to determine the range of an array to return. they could be indexes or scripts that return indexes
	union: a union, two or more comma , seperated values. they can be indexes, keys, or a filter that returns an index or key. used to state the elements in an array or map to return
	**/

	subscript := strings.TrimSpace(tkn[1 : len(tkn)-1])
	if subscript == "" {
		return nil, errors.GetInvalidTokenError("empty subscript")
	}

	if subscript == "*" {
		// range all
		return &token{
			operation: tokenOperationWildcard,
		}, nil
	} else if strings.HasPrefix(subscript, "?") {
		// filter
		if !strings.HasPrefix(subscript, "?(") || !strings.HasSuffix(subscript, ")") {
			return nil, errors.GetInvalidTokenError("expected filter '?(' prefix and ')' suffix")
		}
		return &token{
			operation: tokenOperationFilter,
			arguments: []interface{}{
				strings.TrimSpace(subscript[:]),
			},
		}, nil
	}

	// from this point we have the chance of things being nested or wrapped
	// which would result in the parsing being invalid

	openBracketCount, closeBracketCount := 0, 0
	openQuote := false

	args := []interface{}{}

	remainder := ""
	for idx, rne := range subscript {
		remainder += string(rne)
		switch rne {
		case ' ':
			if !openQuote && openBracketCount == closeBracketCount {
				// do not allow spaces outside of quotes keys or scripts
				return nil, errors.GetInvalidTokenError("unexpected space")
			}
			break
		case '(':
			if openQuote {
				continue
			}
			openBracketCount++
			break
		case ')':
			closeBracketCount++

			if openBracketCount == closeBracketCount {
				// if we are closing bracket, add script to args
				script := remainder[:]
				if !isScript(script) {
					return nil, errors.GetInvalidTokenError("invalid script format")
				}
				args = append(args, script)
				remainder = ""
			}
			break
		case '\'':
			if openBracketCount != closeBracketCount {
				continue
			}
			openQuote = !openQuote

			if openQuote {
				// open quote
				if remainder != "'" {
					return nil, errors.GetInvalidTokenError("unexpected single quote")
				}
			} else {
				// close quote
				if !isKey(remainder) {
					return nil, errors.GetInvalidTokenError("invalid key format")
				}
				args = append(args, remainder[:])
				remainder = ""
			}
			break
		case ':':
			if openQuote || (openBracketCount != closeBracketCount) {
				continue
			}
			if arg := remainder[:len(remainder)-1]; arg != "" {
				if num, err := strconv.Atoi(arg); err == nil {
					args = append(args, num)
				} else {
					return nil, errors.GetInvalidTokenError("only integer or scripts allowed in range arguments")
				}
			} else if idx == 0 {
				// if the token starts with :
				args = append(args, nil)
			}
			args = append(args, ":")

			remainder = ""
			break
		case ',':
			if openQuote || (openBracketCount != closeBracketCount) {
				continue
			}

			if arg := remainder[:len(remainder)-1]; arg != "" {
				if num, err := strconv.Atoi(arg); err == nil {
					args = append(args, num)
				} else {
					args = append(args, arg)
				}
			}
			args = append(args, ",")

			remainder = ""
			break
		}
	}

	if remainder != "" {
		if num, err := strconv.Atoi(remainder); err == nil {
			args = append(args, num)
		} else {
			args = append(args, remainder[:])
		}
	}

	var operation tokenOperation

	if len(args) == 1 {
		// key, index, or script
		arg := args[0]
		if strArg, ok := arg.(string); ok {
			if isKey(strArg) {
				operation = tokenOperationKey
			} else if isScript(strArg) {
				operation = tokenOperationScript
			} else {
				return nil, errors.GetInvalidTokenError("unexpected string")
			}
		} else if _, ok := arg.(int); ok {
			operation = tokenOperationIndex
		} else {
			return nil, errors.GetInvalidTokenError("invalid index")
		}
	} else {
		// range or union
		colonCount := 0
		lastWasColon := false
		commaCount := 0

		includesKeys := false
		justArgs := []interface{}{}

		for _, arg := range args {
			switch arg {
			case ":":
				colonCount++
				if lastWasColon {
					justArgs = append(justArgs, nil)
				}
				lastWasColon = true
				continue
			case ",":
				commaCount++
				break
			default:
				justArgs = append(justArgs, arg)
				if strArg, ok := arg.(string); ok {
					if isKey(strArg) {
						includesKeys = true
					}
				}
				break
			}
			lastWasColon = false
		}

		args = justArgs

		if colonCount > 0 && commaCount > 0 {
			return nil, errors.GetInvalidTokenError("cannot specify a range in a union")
		} else if commaCount > 0 {
			operation = tokenOperationUnion

			// we should always have one more comma than arg
			if commaCount >= len(args) {
				return nil, errors.GetInvalidTokenError("empty argument in union")
			}
			for _, arg := range args {
				if strArg, ok := arg.(string); ok {
					if !isScript(strArg) && !isKey(strArg) {
						return nil, errors.GetInvalidTokenError("unexpected union argument")
					}
				} else if _, ok := arg.(int); !ok {
					return nil, errors.GetInvalidTokenError("unexpected union argument")
				}
			}
		} else if colonCount > 0 {
			if colonCount > 2 {
				return nil, errors.GetInvalidTokenError("incorrect number of arguments in range")
			}
			if colonCount == 1 && colonCount == len(args) {
				args = append(args, nil)
			}
			operation = tokenOperationRange
			if includesKeys {
				return nil, errors.GetInvalidTokenError("only integer or scripts allowed in range arguments")
			}
		}
	}

	return &token{
		operation: operation,
		arguments: args,
	}, nil
}
