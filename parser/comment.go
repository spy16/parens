package parser

// CommentExpr is returned to represent a lisp-style comment.
type CommentExpr struct {
	comment string
}

// Eval returns the comment string.
func (ce CommentExpr) Eval(_ Scope) (interface{}, error) {
	return ce.comment, nil
}
