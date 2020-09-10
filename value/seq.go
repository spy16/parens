package value

var (
	_ Any = (*LinkedList)(nil)
	_ Seq = (*LinkedList)(nil)
)

// NewSeq returns a new sequence containing given values.
func NewSeq(items ...Any) Seq {
	if len(items) == 0 {
		return Seq((*LinkedList)(nil))
	}
	lst := Seq(&LinkedList{})
	for i := len(items) - 1; i >= 0; i-- {
		lst = Cons(items[i], lst)
	}
	return lst
}

// Cons returns a new seq with `v` added as the first and `seq` as the rest. Seq
// can be nil as well.
func Cons(v Any, seq Seq) Seq {
	return &LinkedList{
		first: v,
		rest:  seq,
	}
}

// Seq represents a sequence of values.
type Seq interface {
	Count() (int, error)
	First() (Any, error)
	Next() (Seq, error)
	Conj(items ...Any) (Seq, error)
}

// Seqable can provide a sequence of its values.
type Seqable interface {
	Seq() (Seq, error)
}

// LinkedList implements an immutable Seq using linked-list data structure.
type LinkedList struct {
	first Any
	rest  Seq
}

// Conj returns a new list with all the items added at the head of the list.
func (ll *LinkedList) Conj(items ...Any) (Seq, error) {
	var res Seq
	if ll == nil {
		res = &LinkedList{}
	} else {
		res = ll
	}

	for _, item := range items {
		res = Cons(item, res)
	}
	return res, nil
}

// First returns the head or first item of the list.
func (ll *LinkedList) First() (Any, error) {
	if ll == nil {
		return nil, nil
	}
	return ll.first, nil
}

// Next returns the tail of the list.
func (ll *LinkedList) Next() (Seq, error) {
	if ll == nil {
		return nil, nil
	}
	return ll.rest, nil
}

// Count returns the number of the list.
func (ll *LinkedList) Count() (int, error) {
	if ll == nil {
		return 0, nil
	}

	count := 0
	if ll.first != nil {
		count++
	}

	if ll.rest != nil {
		seqLen, err := ll.rest.Count()
		if err != nil {
			return 0, err
		}
		count += seqLen
	}

	return count, nil
}
