package errs

import (
	"errors"
	"fmt"
)

var (
	ErrUnmapped           = errors.New("unmapped error")
	ErrNotFound           = errors.New("not_found error")
	ErrRecordDoesNotExist = errors.New("not_exist")
	ErrAccountExists      = errors.New("account_exists")
)

type AppError2 struct {
	Err error
	Msg string
}

func (e *AppError2) Error() string {
	return fmt.Sprintf("%s, err: %v", e.Msg, e.Err)
}

func (e *AppError2) Unwrap() error {
	return e.Err
}

func Catch2(err error) *AppError2 {
	if appErr, ok := errors.AsType[*AppError2](err); ok {
		return appErr
	}
	return &AppError2{Err: err, Msg: "unexpected error"}
}
