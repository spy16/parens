package stdlib

import (
	"fmt"
	"reflect"
)

var math = []mapEntry{
	// logical constants
	entry("true", true,
		"Represents logical true",
	),
	entry("false", false,
		"Represents logical false",
	),
	entry("nil", false,
		"Represents logical false. Same as false",
	),

	// functions and operators
	entry("+", Add,
		"Returns sum of all number arguments",
		"Usage: (+ num1 num2 ...)",
	),
	entry("-", Sub,
		"Returns result of subtracting all the right-most args from first arg.",
	),
	entry("*", Mul,
		"Returns result of multiplying all number arguments",
	),
	entry("/", Div,
		"Returns result of continuously dividing first arg by remaining args",
	),
	entry(">", Gt,
		"Returns true if 1st arg is greater than the 2nd",
	),
	entry("<", Lt,
		"Returns true if 1st arg is less than second arg",
	),
	entry("==", Eq,
		"Returns true if all arguments are equal to each other",
	),
	entry("not", Not,
		"Returns true if argument is falsy, false otherwise",
	),
}

// Add returns sum of all the arguments.
func Add(vals ...float64) float64 {
	sum := 0.0
	for _, val := range vals {
		sum += val
	}

	return sum
}

// Sub returns result of subtracting from left-to-right.
func Sub(vals ...float64) float64 {
	if len(vals) == 1 {
		if vals[0] == 0 {
			return 0
		}
		return -1 * vals[0]
	}

	for i := 1; i < len(vals); i++ {
		vals[i] = -1 * vals[i]
	}

	return Add(vals...)
}

// Mul multiplies all numbers.
func Mul(vals ...float64) float64 {
	result := 1.0
	for _, val := range vals {
		result = result * val
	}

	return result
}

// Div divides from left to right.
func Div(vals ...float64) float64 {
	if len(vals) < 2 {
		panic(fmt.Errorf("division requires at least 2 arguments, got %d", len(vals)))
	}

	result := vals[0]
	for i := 1; i < len(vals); i++ {
		result = result / vals[i]
	}

	return result
}

// Gt checks if lval is greater than rval
func Gt(lval, rval float64) bool {
	return lval > rval
}

// Lt checks if lval is lesser than rval
func Lt(lval, rval float64) bool {
	return lval < rval
}

// Eq checks if lval is same as rval
func Eq(vals ...interface{}) bool {
	if len(vals) <= 1 {
		return true
	}

	lval := vals[0]

	for i := 1; i < len(vals); i++ {
		if !reflect.DeepEqual(lval, vals[i]) {
			return false
		}
	}

	return true
}

// Not returns true if val is nil or false value and false
// otherwise.
func Not(val interface{}) bool {
	if b, ok := val.(bool); ok {
		return !b
	}

	if val == nil {
		return true
	}

	return false
}
