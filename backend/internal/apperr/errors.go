package apperr

import "errors"

type Error struct {
	Status  int
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Message }

func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

func As(err error) (*Error, bool) {
	var target *Error
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

var (
	ErrUnauthorized = New(401, "unauthorized", "Sign in to continue")
	ErrNotFound     = New(404, "not_found", "Not found")
)
