package stdlib

import (
	"os"
)

var system = []mapEntry{
	entry("env", os.Getenv,
		"Returns the value of environment variable",
	),
	entry("set-env", setenv,
		"Sets value of environment variable",
		"Example: (set-env \"HELLO\" \"world\")",
	),
}

func setenv(name, val string) string {
	if err := os.Setenv(name, val); err != nil {
		panic(err)
	}

	return val
}
