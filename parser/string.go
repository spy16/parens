package parser

import (
	"strings"

	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/reflection"
)

// StringExp represents single and double quoted strings.
type StringExp struct {
	Token lexer.Token
}

func (se StringExp) Eval(_ *reflection.Env) (interface{}, error) {
	return strings.Trim(se.Token.Value, "\""), nil
}
