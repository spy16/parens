package utfstrings

import (
	"unicode/utf8"
)

// EOS is returned when cursor reaches the end of the string.
const EOS = int32(-1)

// Cursor provides functions to navigate a Unicode string.
type Cursor struct {
	Selection

	String string
}

// Selection represents the current selection in the cursor.
type Selection struct {
	Start int
	Pos   int
	width int
}

// Next returns the next rune in the string.
func (cur *Cursor) Next() rune {
	if int(cur.Pos) >= len(cur.String) {
		cur.width = 0
		return EOS
	}
	ru, width := utf8.DecodeRuneInString(cur.String[cur.Pos:])
	cur.width = width
	cur.Pos += cur.width
	return ru
}

// Peek returns but does not consume the next rune in the string.
func (cur *Cursor) Peek() rune {
	r := cur.Next()
	cur.Backup()
	return r
}

// Backup steps back one rune. Can only be called once per call of next.
func (cur *Cursor) Backup() {
	cur.Pos -= cur.width
}

// Build will build a string from the cursor. move will be called before
// every rune and move can advance the cursor one or more times to ignore
// runes. Build will always start from the beginning of the string. This
// will not modify the original cursor since it has value receiver.
func (cur Cursor) Build(move MoveFunc) string {
	if len(cur.String) == 0 {
		return ""
	}

	cur.Selection = Selection{}
	rus := []rune{}

	for {
		move(&cur)
		ru := cur.Next()
		if ru == EOS {
			break
		}

		rus = append(rus, ru)
	}

	return string(rus)
}

// MoveFunc is called by Build. MoveFunc can advance the cursor one or
// more times. Once this function returns Build will start from the current
// position and consume one rune and call MoveFunc again.
type MoveFunc func(cur *Cursor)
