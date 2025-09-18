/*
Let us assume the following formula for
displacement s as a function of time t, acceleration a, initial velocity vo,
and initial displacement so.

s = ½ a t2 + vot + so

Write a program which first prompts the user
to enter values for acceleration, initial velocity, and initial displacement.
Then the program should prompt the user to enter a value for time and the
program should compute the displacement after the entered time.

You will need to define and use a function
called GenDisplaceFn() which takes three float64
arguments, acceleration a, initial velocity vo, and initial
displacement so. GenDisplaceFn()
should return a function which computes displacement as a function of time,
assuming the given values acceleration, initial velocity, and initial
displacement. The function returned by GenDisplaceFn() should take one float64 argument t, representing time, and return one
float64 argument which is the displacement travelled after time t.

For example, let’s say that I want to assume
the following values for acceleration, initial velocity, and initial
displacement: a = 10, vo = 2, so = 1. I can use the
following statement to call GenDisplaceFn() to
generate a function fn which will compute displacement as a function of time.

fn := GenDisplaceFn(10, 2, 1)

Then I can use the following statement to
print the displacement after 3 seconds.

fmt.Println(fn(3))

And I can use the following statement to print
the displacement after 5 seconds.

fmt.Println(fn(5))
*/
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Enter acceleration, initial velocity, and initial displacement (floats): in a form: <a> <v> <s>")
	a, v, s := getInputs()

	fn := GenDisplaceFn(a, v, s)

	fmt.Println("Enter time value (float): ")
	var t float64
	fmt.Scanln(&t)
	if t == 0 {
		t = 3.00
	}

	displacement := fn(t)
	fmt.Println(displacement)

}

// GenDisplaceFn
// s = (½ * a * t2) + (v * t) + s
// /**GenDisplaceFn
func GenDisplaceFn(a, v, s float64) func(float64) float64 {
	return func(t float64) float64 {
		t2 := t * t
		return (0.5 * a * t2) + (v * t) + s
	}
}

func getInputs() (float64, float64, float64) {
	var a, v, s float64
	fmt.Scanln(&a, &v, &s)

	// Get defaults for easier development and testing
	if a == 0 {
		a = 10.00
	}
	if v == 0 {
		v = 2.00
	}
	if s == 0 {
		s = 1.00
	}
	return a, v, s

}
