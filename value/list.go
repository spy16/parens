package value

var (
	_ Any = (*LinkedList)(nil)
	_ Seq = (*LinkedList)(nil)
)

// Cons returns a new seq with `v` added as the first and `seq` as the rest.
// seq can be nil as well.
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

// NewList returns a new linked-list containing given values.
func NewList(items []Any) (lst Seq, err error) {
	if len(items) == 0 {
		lst = Seq((*LinkedList)(nil))
		return
	}

	lst = Seq(&LinkedList{})
	for i := len(items) - 1; i >= 0; i-- {
		if lst, err = Cons(items[i], lst); err != nil {
			break
		}
	}

	return
}

// LinkedList implements an immutable Seq using linked-list data structure.
type LinkedList struct {
	count int
	first Any
	rest  Seq
}

// SExpr returns a valid s-expression for LinkedList.
func (ll *LinkedList) SExpr() (string, error) {
	if ll == nil {
		return "()", nil
	}

	return SeqString(ll, "(", ")", " ")
}

// Conj returns a new list with all the items added at the head of the list.
func (ll *LinkedList) Conj(items ...Any) (res Seq, err error) {
	if ll == nil {
		res = &LinkedList{}
	} else {
		res = ll
	}

	for _, item := range items {
		if res, err = Cons(item, res); err != nil {
			break
		}
	}

	return
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

	return ll.count, nil
}
