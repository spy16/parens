package stdlib

import (
	"fmt"
)

var io = []mapEntry{
	entry("println", fmt.Println,
		"Concatenates arguments and prints with a newline at the end",
	),
	entry("print", fmt.Print,
		"Concatenates arguments and prints without a newline at the end",
	),
	entry("printf", fmt.Printf,
		"Formats the first string using remaining arguments and prints",
	),
}
