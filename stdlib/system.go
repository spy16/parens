package stdlib

import (
	"fmt"
	"os"

	"github.com/spy16/parens/parser"
)

var system = []mapEntry{
	entry("env", os.Getenv,
		"Returns the value of environment variable",
	),
	entry("set-env", setenv,
		"Sets value of environment variable",
		"Example: (set-env \"HELLO\" \"world\")",
	),
	entry("dump-scope", parser.MacroFunc(dumpScope),
		"Formats and displays the entire scope",
	),
}

func dumpScope(scope parser.Scope, _ string, exprs []parser.Expr) (interface{}, error) {
	return fmt.Sprintln(scope), nil
}

func setenv(name, val string) string {
	if err := os.Setenv(name, val); err != nil {
		panic(err)
	}

	return val
}
