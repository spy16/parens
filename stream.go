package parens

import (
	"io"
	"strings"
	"unicode"
)

// Stream provides functions to read from a rune reader and also maintains
// the stream information such as filename and position in stream.
type Stream struct {
	File string

	rs        io.RuneReader
	buf       []rune
	line, col int
	lastCol   int
}

// NextRune returns next rune from the stream and advances the stream.
func (stream *Stream) NextRune() (rune, error) {
	var r rune
	if len(stream.buf) > 0 {
		r = stream.buf[0]
		stream.buf = stream.buf[1:]
	} else {
		temp, _, err := stream.rs.ReadRune()
		if err != nil {
			return -1, err
		}

		r = temp
	}

	if r == '\n' {
		stream.line++
		stream.lastCol = stream.col
		stream.col = 0
	} else {
		stream.col++
	}

	return r, nil
}

// Unread can be used to return runes consumed from the stream back to the
// stream. Un-reading more runes than read is guaranteed to work but might
// cause inconsistency in stream positional information.
func (stream *Stream) Unread(runes ...rune) {
	newLine := false
	for _, r := range runes {
		if r == '\n' {
			newLine = true
			break
		}
	}

	if newLine {
		stream.line--
		stream.col = stream.lastCol
	} else {
		stream.col--
	}

	stream.buf = append(runes, stream.buf...)
}

// Info returns information about the stream including file name and the
// position of the reader.
func (stream Stream) Info() (file string, line, col int) {
	file = strings.TrimSpace(stream.File)
	return file, stream.line + 1, stream.col
}

// SkipSpaces consumes and discards runes from stream repeatedly until a
// character that is not a whitespace is identified. Along with standard
// unicode  white-space characters "," is also considered  a white-space
// and discarded.
func (stream *Stream) SkipSpaces() error {
	for {
		r, err := stream.NextRune()
		if err != nil {
			return err
		}

		if !isSpace(r) {
			stream.Unread(r)
			break
		}
	}

	return nil
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r) || r == ','
}
