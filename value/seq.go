package value

var (
	_ Any = (*LinkedList)(nil)
	_ Seq = (*LinkedList)(nil)
)

// Seq represents a sequence of values.
type Seq interface {
	SExpr() (string, error)
	Count() (int, error)
	First() (Any, error)
	Next() (Seq, error)
	Conj(items ...Any) (Seq, error)
}

// Seqable can provide a sequence of its values.
type Seqable interface {
	Seq() (Seq, error)
}

// NewSeq returns a new sequence containing given values.
// It is an alias of NewList.
func NewSeq(items ...Any) (Seq, error) {
	return NewList(items)
}

// Cons returns a new seq with `v` added as the first and `seq` as the rest. Seq
// can be nil as well.
func Cons(v Any, seq Seq) (Seq, error) {
	newSeq := &LinkedList{
		first: v,
		rest:  seq,
		count: 1,
	}

	if seq != nil {
		cnt, err := seq.Count()
		if err != nil {
			return nil, err
		}
		newSeq.count = cnt + 1
	}

	return newSeq, nil
}
