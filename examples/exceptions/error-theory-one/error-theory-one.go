package main

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type AppError struct {
	Err  error
	Msg  string
	Type string
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Msg, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Msg)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrTypeOne = errors.New("err type one")
)

func main() {

	err1 := &AppError{Err: ErrTypeOne, Msg: "error one", Type: "type one"}
	err2 := &AppError{Err: pgx.ErrNoRows, Msg: "error two", Type: "no-record-found"}

	fmt.Println(err1)
	fmt.Println(err2)

	if appError, ok := errors.AsType[*AppError](err2); ok {
		fmt.Println("app error type is: ", appError.Type)
	} else {
		fmt.Println("unexpected server error")
	}

	fmt.Println(fmt.Sprintf("err1 is ErrTypeOne: %t", errors.Is(err1, ErrTypeOne)))
	fmt.Println(fmt.Sprintf("err2 is pgx.ErrNoRows: %t", errors.Is(err2, pgx.ErrNoRows)))

	fmt.Println("***********************************************************************")

	err := someFuncThatDoesStuffAndMayReturnErr()
	if err != nil {
		appErr := wrapErr(err)
		fmt.Println(appErr)
		fmt.Println(appErr.Type)
		fmt.Println(fmt.Sprintf("appErr is pgx.ErrNoRows: %t", errors.Is(appErr, pgx.ErrNoRows)))
	}
}

func wrapErr(err error) *AppError {
	if appErr, ok := errors.AsType[*AppError](err); ok {
		return appErr
	}
	return &AppError{Err: err, Msg: "unexpected error"}
}

func someFuncThatDoesStuffAndMayReturnErr() error {
	return &AppError{pgx.ErrNoRows, "opsie opsie", "no-record-found"}
}
