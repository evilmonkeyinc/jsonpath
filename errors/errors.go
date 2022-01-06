package errors

import "fmt"

var (
	// ErrInvalidInitialToken error returned during parsing if an invalid initial token is found.
	ErrInvalidInitialToken error = fmt.Errorf("query must start with '$'")
	// ErrIllegalCharacterAtPositionOne error returned during parsing if '.' or '[' doesn't follow the initial token.
	ErrIllegalCharacterAtPositionOne = fmt.Errorf("expected '.' or '[' after initial token")
	// ErrQueryNotSpecified error returned when a query is required or parsed and it has not be specified or is empty
	ErrQueryNotSpecified = fmt.Errorf("no valid JSONPath query has been specified")
	// ErrInvalidToken err returned while parsing a token which is invalid
	ErrInvalidToken = fmt.Errorf("invalid token")

	// TODO : tody up and test errors
	ErrGetKeyFromNilMap         error = fmt.Errorf("cannot get key from nil map")
	ErrGetIndexFromNilSlice     error = fmt.Errorf("cannot get index from nil slice")
	ErrGetRangeFromNilSlice     error = fmt.Errorf("cannot get range from nil slice")
	ErrGetUnionFromNilObject    error = fmt.Errorf("cannot get union from nil object")
	ErrGetElementsFromNilObject error = fmt.Errorf("cannot get elements from nil object")

	ErrInvalidObject           error = fmt.Errorf("invalid object")
	ErrInvalidObjectMap        error = fmt.Errorf("%w. expected map", ErrInvalidObject)
	ErrInvalidObjectSlice      error = fmt.Errorf("%w. expected slice", ErrInvalidObject)
	ErrInvalidObjectMapOrSlice error = fmt.Errorf("%w. expected map or slice", ErrInvalidObject)

	ErrKeyNotFound error = fmt.Errorf("key not found in object")

	ErrIndexOutOfRange error = fmt.Errorf("index out of range")

	ErrInvalidParameter error = fmt.Errorf("invalid parameter")
)

// GetInvalidTokenError returns an ErrInvalidToken, appended with the specified reason.
func GetInvalidTokenError(reason string) error {
	return fmt.Errorf("%w: %s", ErrInvalidToken, reason)
}

func GetKeyNotFoundError(key string) error {
	return fmt.Errorf("'%s' %w", key, ErrKeyNotFound)
}

func GetInvalidParameterError(reason string) error {
	return fmt.Errorf("%w. %s", ErrInvalidParameter, reason)
}
