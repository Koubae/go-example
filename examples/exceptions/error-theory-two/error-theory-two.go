package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ErrType string

const (
	unmapped    ErrType = "unmapped"
	NotFound    ErrType = "not_found"
	WalletEmpty ErrType = "wallet_empty"

	WrapLv1 ErrType = "wrap_lv1"
	WrapLv2 ErrType = "wrap_lv2"
	WrapLv3 ErrType = "wrap_lv3"
	WrapLv4 ErrType = "wrap_lv4"
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
	return e.Type == unmapped
}

func Catch(err error) *AppError {
	if appErr, ok := errors.AsType[*AppError](err); ok {
		return appErr
	}
	return &AppError{Err: err, Msg: "unexpected error", Type: unmapped}
}

var (
	SomeGenericErr = errors.New("some generic error")
	WrappedLevel1  = fmt.Errorf("wrapped: %w", SomeGenericErr)
	WrappedLevel2  = fmt.Errorf("wrapped: %w", WrappedLevel1)
	WrappedLevel3  = fmt.Errorf("wrapped: %w", WrappedLevel2)
	WrappedLevel4  = fmt.Errorf("wrapped: %w", WrappedLevel3)
)

func main() {
	exampleOne()
	exampleTwo()
	exampleThree()
}

func printExample(name string) {
	fmt.Println("\n===============================================================")
	fmt.Println("\t\t\t" + name)
	fmt.Println("===============================================================\n")
}

// Checking if error wrapping works!
func exampleOne() {
	printExample("exampleOne")
	someFunc := func() error {
		return &AppError{Err: SomeGenericErr, Msg: "some error", Type: WalletEmpty}
	}
	fmt.Println("App running...")
	if err := someFunc(); err != nil {
		appErr := Catch(err)
		fmt.Println(appErr)
		fmt.Println(fmt.Sprintf("is SomeGeneric: %t", errors.Is(appErr, SomeGenericErr)))
	}
}

// here I'm checking whether some nested error wrapping still working!
func exampleTwo() {
	printExample("exampleTwo")
	someFunc := func(i int) error {
		switch i {
		case 1:
			return &AppError{Err: WrappedLevel1, Msg: "some error", Type: WrapLv1}
		case 2:
			return &AppError{Err: WrappedLevel2, Msg: "some error", Type: WrapLv2}
		case 3:
			return &AppError{Err: WrappedLevel3, Msg: "some error", Type: WrapLv3}
		case 4:
			return &AppError{Err: WrappedLevel4, Msg: "some error", Type: WrapLv4}
		default:
			return sql.ErrNoRows
		}
	}

	for _, i := range []int{1, 2, 3, 4, 5} {
		if err := someFunc(i); err != nil {
			appErr := Catch(err)
			fmt.Println(fmt.Sprintf("i: %d, appErr: %v", i, appErr))

			fmt.Println(fmt.Sprintf("SomeGenericErr: %t", errors.Is(appErr, SomeGenericErr)))
			fmt.Println(fmt.Sprintf("WrappedLevel1: %t", errors.Is(appErr, WrappedLevel1)))
			fmt.Println(fmt.Sprintf("WrappedLevel2: %t", errors.Is(appErr, WrappedLevel2)))
			fmt.Println(fmt.Sprintf("WrappedLevel3: %t", errors.Is(appErr, WrappedLevel3)))
			fmt.Println(fmt.Sprintf("WrappedLevel4: %t", errors.Is(appErr, WrappedLevel4)))

			switch appErr.Type {
			case WrapLv1:
				fmt.Println("Type: WrapLv1!!!!")
			case WrapLv2:
				fmt.Println("Type: WrapLv2!!!!")
			case WrapLv3:
				fmt.Println("Type: WrapLv3!!!!")
			case WrapLv4:
				fmt.Println("Type: WrapLv4!!!!")
			case unmapped:
				fmt.Println("Type: unmapped!!!!")
			}

			fmt.Println("***********************************************************************")
		}
	}
}

func exampleThree() {
	printExample("exampleThree")

	// e := ErrType("blabal")
	// peasantString := "someotherthing"
	someFunc := func() error {
		return &AppError{Err: WrappedLevel1, Msg: "some error", Type: "peasantString"}
		// return &AppError{Err: WrappedLevel1, Msg: "some error", Type: peasantString}
	}

	fmt.Println("App running...")
	if err := someFunc(); err != nil {
		appErr := Catch(err)
		fmt.Println(appErr)
		fmt.Println(fmt.Sprintf("is SomeGeneric: %t", errors.Is(appErr, SomeGenericErr)))

		fmt.Println(appErr.Type == "peasantString")
		fmt.Println(appErr.Type == ErrType("peasantString"))
		switch appErr.Type {
		case WrapLv1:
			fmt.Println("Type: WrapLv1!!!!")
		case WrapLv2:
			fmt.Println("Type: WrapLv2!!!!")
		case WrapLv3:
			fmt.Println("Type: WrapLv3!!!!")
		case WrapLv4:
			fmt.Println("Type: WrapLv4!!!!")
		case unmapped:
			fmt.Println("Type: unmapped!!!!")
		case "peasantString":
			fmt.Println("Type: peasantString!!!!")
		default:
			fmt.Println("dunno what to do with this")
		}

	}
}
