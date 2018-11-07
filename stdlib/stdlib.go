package stdlib

import (
	"github.com/spy16/parens/parser"
)

// RegisterBuiltins registers different built-in functions into the
// given scope.
func RegisterBuiltins(scope parser.Scope) error {
	return doUntilErr(scope,
		RegisterMacros,
		RegisterMath,
		RegisterIO,
	)
}

// RegisterIO binds input/output functions into the scope.
func RegisterIO(scope parser.Scope) error {
	return registerList(scope, io)
}

// RegisterMacros binds all the macros into the scope.
func RegisterMacros(scope parser.Scope) error {
	return registerList(scope, macros)
}

// RegisterMath binds basic math operators into the scope.
func RegisterMath(scope parser.Scope) error {
	return registerList(scope, math)
}

func doUntilErr(scope parser.Scope, fns ...func(scope parser.Scope) error) error {
	for _, fn := range fns {
		if err := fn(scope); err != nil {
			return err
		}
	}

	return nil
}

func registerList(scope parser.Scope, entries []mapEntry) error {
	for _, entry := range entries {
		if err := scope.Bind(entry.name, entry.val, entry.doc...); err != nil {
			return err
		}
	}

	return nil
}

func entry(name string, val interface{}, doc ...string) mapEntry {
	return mapEntry{
		name: name,
		val:  val,
		doc:  doc,
	}
}

type mapEntry struct {
	name string
	val  interface{}
	doc  []string
}
