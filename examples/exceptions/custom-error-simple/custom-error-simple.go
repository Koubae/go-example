package main

import (
	"errors"
	"fmt"
)

type customErr struct {
	Err  error
	code int
}

func (e *customErr) Error() string {
	return fmt.Sprintf("custom error, code %d, error: %s", e.code, e.Err.Error())
}

type specificCustomErr struct {
	// customErr
	Err  error
	code int
}

func (e *specificCustomErr) Error() string {
	return fmt.Sprintf("custom error, code %d, error: %s", e.code, e.Err.Error())
}
func (e *specificCustomErr) Unwrap() error {
	return e.Err
}

type ValidationError struct {
	Msg string
	Err error
}

func (e *ValidationError) Error() string {
	return e.Msg + ": " + e.Err.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

var dbError = errors.New("db-mock-error")

func main() {
	myErr := &customErr{Err: dbError, code: 404}
	fmt.Println(myErr.Error())

	if errors.Is(dbError, myErr) {
		fmt.Println("dbError is myErr")
	} else if errors.Is(myErr, &customErr{}) {
		fmt.Println("myErr is customErr")
	} else {
		fmt.Println("myErr is not customErr")
	}

	if e, ok := errors.AsType[*customErr](myErr); ok {
		fmt.Println("asType WORKS!!!! ", e)
	} else {
		fmt.Println("wtf, is still not working??")
	}

	nestedErr := fmt.Errorf("nested error: %w", myErr)
	fmt.Println(fmt.Sprintf("errors.Is: %t", errors.Is(nestedErr, myErr)))

	specificErr := &specificCustomErr{Err: dbError, code: 404}
	var sErr *specificCustomErr
	fmt.Println(fmt.Sprintf("specificErr: errors.Is: %t", errors.Is(specificErr, dbError)))
	fmt.Println(fmt.Sprintf("specificErr: errors.Is: %t", errors.Is(specificErr, &specificCustomErr{})))
	fmt.Println(fmt.Sprintf("specificErr: errors.Is: %t", errors.As(specificErr, &sErr)))
}
