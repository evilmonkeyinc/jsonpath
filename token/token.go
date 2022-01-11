package token

import (
	"strconv"
	"strings"

	"github.com/evilmokeyinc/jsonpath/errors"
)

// TODO : something like [:2] is not the same as [0:2] it is more like [0:1]

// Token represents a component of a JSON Path query
type Token interface {
	Apply(root, current interface{}, next []Token) (interface{}, error)
}

// TODO : need fix and unit tests for just giving $ or @

// Tokenize converts a JSON Path query to a collection of parsable tokens
func Tokenize(query string) ([]string, string, error) {
	if query == "" {
		return nil, query, errors.ErrQueryNotSpecified
	}

	tokens := []string{}
	tokenString := ""
	remainder := query

tokenize:
	for idx, rne := range query {
		tokenString += string(rne)
		remainder = remainder[1:]

		if idx == 0 {
			if tokenString != "$" && tokenString != "@" {
				return nil, "", errors.ErrInvalidInitialToken
			}

			if next := query[1]; next != '.' && next != '[' {
				return nil, "", errors.ErrInvalidInitialToken
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
		} else {
			// check for script operators outside of subscript
			switch rne {
			case '*':
				// '*' is an operator if it is part of a larger token
				// '*' is a wildcard if by self or with proceding '.'
				if tokenString == ".*" || tokenString == "*" {
					continue
				}
				fallthrough
			case '-', '+', '/', '%', '>', '<', '=', '!':
				// strip operator and break tokenize loop
				tokenString = tokenString[0 : len(tokenString)-1]
				remainder = query[idx:]

				break tokenize
			default:
				// not a script operator
				continue
			}
		}
	}

	// parse the last token
	if len(tokenString) > 0 {
		if tokenString[0] == '.' {
			tokenString = tokenString[1:]
		}

		tokens = append(tokens, tokenString[:])
	}

	return tokens, remainder, nil
}

// Parse will parse a single token string and return an actionable token
func Parse(tokenString string) (Token, error) {

	isScript := func(token string) bool {
		return len(token) > 2 && strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")")
	}

	isKey := func(token string) bool {
		return len(token) > 2 && strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'")
	}

	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, errors.ErrInvalidTokenEmpty
	}

	if tokenString == "$" {
		return &rootToken{}, nil
	}
	if tokenString == "@" {
		return &currentToken{}, nil
	}
	if tokenString == "*" {
		return &wildcardToken{}, nil
	}
	if tokenString == ".." {
		return &recursiveToken{}, nil
	}

	if !strings.HasPrefix(tokenString, "[") {

		if _, err := strconv.ParseInt(tokenString, 10, 64); err == nil {
			return nil, errors.ErrInvalidTokenUnexpectedIndex
		}

		if tokenString == "length" {
			return &lengthToken{}, nil
		}

		return &keyToken{key: tokenString}, nil

	}

	if !strings.HasSuffix(tokenString, "]") {
		return nil, errors.ErrInvalidTokenMissingSubscriptClose
	}
	// subscript, or child operator

	subscript := strings.TrimSpace(tokenString[1 : len(tokenString)-1])
	if subscript == "" {
		return nil, errors.ErrInvalidTokenEmptySubscript
	}

	if subscript == "*" {
		// range all
		return &wildcardToken{}, nil
	} else if strings.HasPrefix(subscript, "?") {
		// filter
		if !strings.HasPrefix(subscript, "?(") || !strings.HasSuffix(subscript, ")") {
			return nil, errors.ErrInvalidTokenInvalidFilterFormat
		}
		return &filterToken{
			expression: strings.TrimSpace(subscript[2 : len(subscript)-1]),
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
				return nil, errors.ErrInvalidTokenUnexpectedSpace
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
					return nil, errors.ErrInvalidTokenInvalidScriptFormat
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
					return nil, errors.ErrInvalidTokenUnexpectedQuote
				}
			} else {
				// close quote
				if !isKey(remainder) {
					return nil, errors.ErrInvalidTokenInvalidKeyFormat
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
				if num, err := strconv.ParseInt(arg, 10, 64); err == nil {
					args = append(args, num)
				} else {
					return nil, errors.ErrInvalidTokenInvalidRangeArguments
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
				if num, err := strconv.ParseInt(arg, 10, 64); err == nil {
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
		if num, err := strconv.ParseInt(remainder, 10, 64); err == nil {
			args = append(args, num)
		} else {
			args = append(args, remainder[:])
		}
	}

	if len(args) == 1 {
		// key, index, or script
		arg := args[0]
		if strArg, ok := arg.(string); ok {
			if isKey(strArg) {
				return &keyToken{
					key: strArg[1 : len(strArg)-1],
				}, nil
			} else if isScript(strArg) {
				return &scriptToken{
					expression: strArg[1 : len(strArg)-1],
				}, nil
			}
			return nil, errors.ErrInvalidTokenUnexpectedString
		} else if intArg, ok := isInteger(arg); ok {
			return &indexToken{index: intArg}, nil
		}
		return nil, errors.ErrInvalidTokenInvalidIndex
	}

	// range or union
	colonCount := 0
	lastWasColon := false
	commaCount := 0

	// includesKeys := false
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
			break
		}
		lastWasColon = false
	}

	args = justArgs

	if colonCount > 0 && commaCount > 0 {
		return nil, errors.ErrInvalidTokenNoRangeInUnion
	} else if commaCount > 0 {
		// Union

		// we should always have one more comma than arg
		if commaCount >= len(args) {
			return nil, errors.ErrInvalidTokenEmptyUnionArguments
		}
		for idx, arg := range args {
			if strArg, ok := arg.(string); ok {
				if isScript(strArg) {
					arg = &expressionToken{
						expression: strArg[1 : len(strArg)-1],
					}
					args[idx] = arg
				} else if isKey(strArg) {
					args[idx] = strArg[1 : len(strArg)-1]
				} else {
					return nil, errors.ErrInvalidTokenUnexpectedUnionArguments
				}
			} else if _, ok := isInteger(arg); !ok {
				return nil, errors.ErrInvalidTokenUnexpectedUnionArguments
			}
		}

		return &unionToken{arguments: args}, nil
	} else if colonCount > 0 {
		// Range
		if colonCount > 2 {
			return nil, errors.ErrInvalidTokenIncorrectNumberOfRangeArguments
		}
		if colonCount == 1 && colonCount == len(args) {
			args = append(args, nil)
		}

		var from, to, step interface{} = args[0], args[1], int64(1)
		if len(args) > 2 {
			step = args[2]
		}

		if strFrom, ok := from.(string); ok {
			if !isScript(strFrom) {
				return nil, errors.ErrInvalidTokenInvalidRangeArguments
			}
			from = &expressionToken{
				expression: strFrom[1 : len(strFrom)-1],
			}
		}
		if strTo, ok := to.(string); ok {
			if !isScript(strTo) {
				return nil, errors.ErrInvalidTokenInvalidRangeArguments
			}
			to = &expressionToken{
				expression: strTo[1 : len(strTo)-1],
			}
		}
		if strStep, ok := step.(string); ok {
			if !isScript(strStep) {
				return nil, errors.ErrInvalidTokenInvalidRangeArguments
			}
			step = &expressionToken{
				expression: strStep[1 : len(strStep)-1],
			}
		}

		return &rangeToken{
			from: from,
			to:   to,
			step: step,
		}, nil
	}

	return nil, errors.ErrInvalidQueryUnexpectedTokens
}
