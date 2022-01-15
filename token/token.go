package token

import (
	"strconv"
	"strings"
)

// Token represents a component of a JSON Path query
type Token interface {
	Apply(root, current interface{}, next []Token) (interface{}, error)
	Type() string
	String() string
}

// Tokenize converts a JSON Path query to a collection of parsable tokens
func Tokenize(query string) ([]string, string, error) {
	if query == "" {
		return nil, query, getUnexpectedTokenError("", 0)
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
				return nil, "", getUnexpectedTokenError(string(rne), idx)
			}

			if len(query) > 1 {
				if next := query[1]; next != '.' && next != '[' {
					return nil, "", getUnexpectedTokenError(string(next), idx+1)
				}
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

// ParseOptions represents the options for the parse function
type ParseOptions struct {
	// IsString if true will allow Parse() to error
	// if there is additional whitespace in tokens
	// or other minor format issues that could be ignored
	IsStrict bool
}

// Parse will parse a single token string and return an actionable token
func Parse(tokenString string, options *ParseOptions) (Token, error) {
	if options == nil {
		options = &ParseOptions{
			IsStrict: true,
		}
	}

	isScript := func(token string) bool {
		return len(token) > 2 && strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")")
	}

	isKey := func(token string) bool {
		return len(token) > 2 && strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'")
	}

	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, getInvalidTokenEmpty()
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
			return nil, getInvalidTokenFormatError(tokenString)
		}

		if tokenString == "length" {
			return &lengthToken{}, nil
		}

		return &keyToken{key: tokenString}, nil

	}

	if !strings.HasSuffix(tokenString, "]") {
		return nil, getInvalidTokenFormatError(tokenString)
	}
	// subscript, or child operator

	subscript := strings.TrimSpace(tokenString[1 : len(tokenString)-1])
	if subscript == "" {
		return nil, getInvalidTokenFormatError(tokenString)
	}

	if subscript == "*" {
		// range all
		return &wildcardToken{}, nil
	} else if strings.HasPrefix(subscript, "?") {
		// filter
		if !strings.HasPrefix(subscript, "?(") || !strings.HasSuffix(subscript, ")") {
			return nil, getInvalidTokenFormatError(tokenString)
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

	bufferString := ""
	for idx, rne := range subscript {
		bufferString += string(rne)
		switch rne {
		case ' ':

			if !openQuote && openBracketCount == closeBracketCount {
				if options.IsStrict {
					// do not allow spaces outside of quotes keys or scripts
					return nil, getInvalidTokenFormatError(tokenString)
				}

				// remove whitespace
				bufferString = strings.TrimSpace(bufferString)
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
				script := bufferString[:]
				if !isScript(script) {
					return nil, getInvalidExpressionFormatError(script)
				}
				args = append(args, script)
				bufferString = ""
			}
			break
		case '\'':
			if openBracketCount != closeBracketCount {
				continue
			}
			openQuote = !openQuote

			if openQuote {
				// open quote
				if bufferString != "'" {
					return nil, getInvalidTokenFormatError(tokenString)
				}
			} else {
				// close quote
				if !isKey(bufferString) {
					return nil, getInvalidTokenFormatError(tokenString)
				}
				args = append(args, bufferString[:])
				bufferString = ""
			}
			break
		case ':':
			if openQuote || (openBracketCount != closeBracketCount) {
				continue
			}
			if arg := bufferString[:len(bufferString)-1]; arg != "" {
				if num, err := strconv.ParseInt(arg, 10, 64); err == nil {
					args = append(args, num)
				} else {
					return nil, getInvalidTokenFormatError(tokenString)
				}
			} else if idx == 0 {
				// if the token starts with :
				args = append(args, nil)
			}
			args = append(args, ":")

			bufferString = ""
			break
		case ',':
			if openQuote || (openBracketCount != closeBracketCount) {
				continue
			}

			if arg := bufferString[:len(bufferString)-1]; arg != "" {
				if num, err := strconv.ParseInt(arg, 10, 64); err == nil {
					args = append(args, num)
				} else {
					args = append(args, arg)
				}
			}
			args = append(args, ",")

			bufferString = ""
			break
		}
	}

	if bufferString != "" {
		if num, err := strconv.ParseInt(bufferString, 10, 64); err == nil {
			args = append(args, num)
		} else {
			args = append(args, bufferString[:])
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
			return nil, getInvalidTokenFormatError(tokenString)
		} else if intArg, ok := isInteger(arg); ok {
			return &indexToken{index: intArg}, nil
		}
		return nil, getInvalidTokenFormatError(tokenString)
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
				// cannot have two colons in a row
				return nil, getInvalidTokenFormatError(tokenString)
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
		return nil, getInvalidTokenFormatError(tokenString)
	} else if commaCount > 0 {
		// Union

		// we should always have one less comma than args
		if commaCount >= len(args) {
			return nil, getInvalidTokenFormatError(tokenString)
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
					return nil, getInvalidTokenFormatError(tokenString)
				}
			} else if _, ok := isInteger(arg); !ok {
				return nil, getInvalidTokenFormatError(tokenString)
			}
		}

		return &unionToken{arguments: args}, nil
	} else if colonCount > 0 {
		// Range
		if colonCount > 2 {
			return nil, getInvalidTokenFormatError(tokenString)
		}
		if colonCount == 1 && len(args) == 1 {
			// to help support [x:] tokens
			args = append(args, nil)
		}

		var from, to, step interface{} = args[0], args[1], nil
		if len(args) > 2 {
			step = args[2]
		}

		if from == nil {
			// This could be a slice token if step is not set
			if len(args) == 2 && args[1] != nil {

				number := args[1]
				if strValue, ok := number.(string); ok {
					if !isScript(strValue) {
						return nil, getInvalidExpressionFormatError(strValue)
					}
					number = &expressionToken{
						expression: strValue[1 : len(strValue)-1],
					}
				}

				return &sliceToken{
					number: number,
				}, nil
			}

			return nil, getInvalidTokenFormatError(tokenString)
		}

		if strFrom, ok := from.(string); ok {
			if !isScript(strFrom) {
				return nil, getInvalidExpressionFormatError(strFrom)
			}
			from = &expressionToken{
				expression: strFrom[1 : len(strFrom)-1],
			}
		}
		if strTo, ok := to.(string); ok {
			if !isScript(strTo) {
				return nil, getInvalidExpressionFormatError(strTo)
			}
			to = &expressionToken{
				expression: strTo[1 : len(strTo)-1],
			}
		}
		if strStep, ok := step.(string); ok {
			if !isScript(strStep) {
				return nil, getInvalidExpressionFormatError(strStep)
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

	return nil, getInvalidTokenFormatError(tokenString)
}
