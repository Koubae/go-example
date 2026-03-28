package errs

import (
	"errors"
	"fmt"
)

const Thing = "thing"

type ErrType string

const (
	Unmapped           ErrType = "unmapped"
	NotFound           ErrType = "not_found"
	RecordDoesNotExist ErrType = "not_exist"
	WalletEmpty        ErrType = "wallet_empty"
	AccountExists              = "account_exists"
	ServerErr          ErrType = "server_err"
)

type AppError struct {
	Err  error
	Msg  string
	Type ErrType
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s, err: %v", e.Type, e.Msg, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) IsUnmapped() bool {
	return e.Type == Unmapped
}

func Catch(err error) *AppError {
	if appErr, ok := errors.AsType[*AppError](err); ok {
		return appErr
	}
	return &AppError{Err: err, Msg: "unexpected error", Type: Unmapped}
}
