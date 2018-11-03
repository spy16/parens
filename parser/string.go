package parser

import (
	"github.com/spy16/parens/lexer"
	"github.com/spy16/parens/lexer/utfstrings"
	"github.com/spy16/parens/reflection"
)

// StringExp represents single and double quoted strings.
type StringExp struct {
	Token lexer.Token
}

// Eval returns unquoted version of the STRING token.
func (se StringExp) Eval(_ *reflection.Scope) (interface{}, error) {
	return unquoteStr(se.Token.Value), nil
}

func unquoteStr(str string) string {
	sc := &utfstrings.Cursor{
		String: str[1 : len(str)-1],
	}

	final := sc.Build(func(cur *utfstrings.Cursor) {
		if ru := cur.Next(); ru == '\\' && cur.Peek() == '"' {
			return
		}
		cur.Backup()
	})
	return final
}
