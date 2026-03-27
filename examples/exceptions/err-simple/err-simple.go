// reference: https://gobyexample.com/errors
package main

import (
	"errors"
	"fmt"
)

var ErrBase = errors.New("app-error")

var ErrValidation = fmt.Errorf("%w err validation", ErrBase)
var ErrNotGreaterThan10 = fmt.Errorf("%w err not greater than 10", ErrValidation)

/*
	[ErrBase]
		|

[ErrValidation]

	|

[ErrNotGreaterThan10]
*/
func main() {
	for _, i := range []int{1, 2, 3, 4} {
		if err := validator(i); err != nil {

			if errors.Is(err, ErrNotGreaterThan10) {
				fmt.Println(fmt.Sprintf("i %d, err is not greater than 10, err: %s", i, err.Error()))
			} else if errors.Is(err, ErrValidation) {
				fmt.Println(fmt.Sprintf("i %d, err is validation, err: %s", i, err.Error()))
			} else if errors.Is(err, ErrBase) {
				fmt.Println(fmt.Sprintf("i %d, err is base, err: %s", i, err.Error()))
			} else {
				fmt.Println(fmt.Sprintf("i %d, unknown error, err: %s", i, err.Error()))
			}

		} else {
			fmt.Println("ok")
		}
	}

	fmt.Println("***********************************************************************")

	for _, i := range []int{1, 2, 3, 4} {

		if err := validator(i); err != nil {

			switch {
			case errors.Is(err, ErrNotGreaterThan10):
				fmt.Println(fmt.Sprintf("i %d, err is not greater than 10, err: %s", i, err.Error()))
			case errors.Is(err, ErrValidation):
				fmt.Println(fmt.Sprintf("i %d, err is validation, err: %s", i, err.Error()))
			case errors.Is(err, ErrBase):
				fmt.Println(fmt.Sprintf("i %d, err is base, err: %s", i, err.Error()))
			default:
				fmt.Println(fmt.Sprintf("i %d, unknown error, err: %s", i, err.Error()))
			}

		} else {
			fmt.Println("ok")
		}
	}

}

func validator(val int) error {
	switch val {
	case 1:
		return ErrBase
	case 2:
		return ErrValidation
	case 3:
		return ErrNotGreaterThan10
	default:
		return errors.New("unknown error")
	}

}
