package jsonpath

import (
	goErr "errors"
	"fmt"

	"github.com/evilmonkeyinc/jsonpath/errors"
)

var (
	errDataIsUnexpectedTypeOrNil error = fmt.Errorf("unexpected type or nil")
)

func getInvalidJSONData(reason error) error {
	return fmt.Errorf("%w. %s", errors.ErrInvalidJSONData, reason.Error())
}

func getInvalidJSONPathQuery(query string) error {
	return fmt.Errorf("%w '%s'", errors.ErrInvalidJSONPathQuery, query)
}

func getInvalidJSONPathQueryWithReason(query string, reason error) error {
	if goErr.Is(reason, errors.ErrInvalidJSONPathQuery) {
		return reason
	}
	return fmt.Errorf("%w '%s' %s", errors.ErrInvalidJSONPathQuery, query, reason.Error())
}
