package value

// Any represents any value.
type Any interface {
	// SExpr MUST return a parsable s-expression that can be consumed by
	// a reader.Reader.
	//
	// For a human-readable implementation, implement `repl.Renderable`.
	SExpr() (string, error)
}
