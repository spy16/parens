package parser

// Result represents result of an Eval call.
type Result struct {
	Value interface{}
	Err   error
}

// IsErr returns true if there is an error in the result.
func (res *Result) IsErr() bool {
	return res.Err != nil
}
