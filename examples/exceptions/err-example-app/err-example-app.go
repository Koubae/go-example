package main

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/koubae/go-example/examples/exceptions/err-example-app/errs"
)

func main() {
	handlerGetAccount()
	handlerGetAccountType2()
}

func handlerGetAccount() {

	err := serviceGetAccount()
	if err != nil {
		appErro := errs.Catch(err)
		switch appErro.Type {
		case errs.RecordDoesNotExist:
			fmt.Println("404")
		case errs.Unmapped:
			fmt.Println("500")
		default:
			fmt.Println("500")
		}
	}
}

func serviceGetAccount() error {
	// GO to a repo calling SQL
	// err := accountRepoGet()
	return &errs.AppError{
		Err:  fmt.Errorf("service err repo error %w", errs.ErrAccountExists),
		Msg:  "error-getting-account",
		Type: errs.RecordDoesNotExist,
	}
}

func handlerGetAccountType2() {
	fmt.Println("==================================")
	fmt.Println("handlerGetAccountType2")
	err := serviceGetAccountType2()
	if err != nil {
		appErro := errs.Catch2(err)
		switch {
		case errors.Is(appErro, errs.ErrNotFound):
			fmt.Println("404")
			fmt.Println(fmt.Sprintf("pgx.ErrNoRows: %t", errors.Is(appErro, pgx.ErrNoRows)))

		case errors.Is(appErro, errs.ErrUnmapped):
			fmt.Println("500")
		default:
			fmt.Println("500")
		}
	}
}

func serviceGetAccountType2() error {
	// GO to a repo calling SQL
	// err := accountRepoGet()
	dbErrCustom := fmt.Errorf("service err repo error %w", pgx.ErrNoRows)
	return &errs.AppError{
		Err:  errors.Join(dbErrCustom, errs.ErrNotFound),
		Msg:  "error-getting-account",
		Type: errs.RecordDoesNotExist,
	}
}

// func accountRepoGet() error {
// 	// SELECT ... FROM account WHERE id = ...
// 	// 404
// 	return &errs.AppError{Err: sql.ErrNoRows, Msg: "no-record-found", Type: errs.RecordDoesNotExist}
//
// }

// func makeHTTPHandler(h apiHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
//
// 		if err := h(w, r); err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 	} // Generic user-facing error } } }
