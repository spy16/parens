package reflection

import "errors"

var (
	// ErrNameNotFound is returned when a lookup is performed with a
	// non-bound name.
	ErrNameNotFound = errors.New("name not bound to a value")

	// ErrNotCallable is returned when a Call is attempted on a non-
	// callable value.
	ErrNotCallable = errors.New("value is not callable")

	// ErrConversionImpossible is returned when the Value type cannot be
	// converted to the expected type.
	ErrConversionImpossible = errors.New("cannot be converted")

	// ErrInvalidNumberOfArgs is returned when a function call is attempted
	// with invalid number of arguments.
	ErrInvalidNumberOfArgs = errors.New("invalid number of arguments")
)
