package value

// Any represents any value.
type Any interface {
	// SExpr MUST return a parsable s-expression that can be consumed by
	// a reader.Reader.
	//
	// For a human-readable implementation, implement `repl.Renderable`.
	SExpr() (string, error)
}

// Seq represents a sequence of values.
type Seq interface {
	Any
	Count() (int, error)
	First() (Any, error)
	Next() (Seq, error)
	Conj(items ...Any) (Seq, error)
}
