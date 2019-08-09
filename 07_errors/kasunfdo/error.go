package main

import "fmt"

type ErrCode uint32

type Error struct {
	Message string
	Code    ErrCode
	Cause   error
}

const (
	ErrInvalid  ErrCode = 400
	ErrNotFound ErrCode = 404
	ErrInternal ErrCode = 500
)

func (e ErrCode) String() string {
	switch e {
	case ErrInvalid:
		return "invalid input: %v"
	case ErrNotFound:
		return "not found: %v"
	case ErrInternal:
		return "internal error"
	default:
		return "error occurred"
	}
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%s\n\t%s", e.Message, e.Cause.Error())
}

func NewError(code ErrCode, cause error, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(code.String(), args...),
		Code:    code,
		Cause:   cause,
	}
}
