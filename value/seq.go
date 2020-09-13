package value

var (
	_ Any = (*LinkedList)(nil)
	_ Seq = (*LinkedList)(nil)
)

// Seq represents a sequence of values.
type Seq interface {
	Count() (int, error)
	First() (Any, error)
	Next() (Seq, error)
	Conj(items ...Any) (Seq, error)
}

type Seqable interface {
	Seq() (Seq, error)
}
