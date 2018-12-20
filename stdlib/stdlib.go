package stdlib

import "github.com/spy16/parens"

// RegisterAll registers different built-in functions into the
// given scope.
func RegisterAll(scope parens.Scope) error {
	return doUntilErr(scope,
		RegisterCore,
		RegisterMath,
		RegisterIO,
		RegisterSystem,
	)
}

// RegisterSystem binds system functions into the scope.
func RegisterSystem(scope parens.Scope) error {
	return registerList(scope, system)
}

// RegisterIO binds input/output functions into the scope.
func RegisterIO(scope parens.Scope) error {
	return registerList(scope, io)
}

// RegisterCore binds all the core macros and functions into
// the scope.
func RegisterCore(scope parens.Scope) error {
	return registerList(scope, core)
}

// RegisterMath binds basic math operators into the scope.
func RegisterMath(scope parens.Scope) error {
	return registerList(scope, math)
}

func doUntilErr(scope parens.Scope, fns ...func(scope parens.Scope) error) error {
	for _, fn := range fns {
		if err := fn(scope); err != nil {
			return err
		}
	}

	return nil
}

func registerList(scope parens.Scope, entries []mapEntry) error {
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
