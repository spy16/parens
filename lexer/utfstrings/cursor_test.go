package utfstrings_test

import (
	"testing"

	"github.com/spy16/parens/lexer/utfstrings"
	"github.com/stretchr/testify/assert"
)

func TestCursor_Next(suite *testing.T) {
	suite.Parallel()

	suite.Run("WithEmptyString", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "",
		}

		assert.Equal(t, utfstrings.EOS, cur.Next())
	})

	suite.Run("WithOneRune", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "≠",
		}

		assert.Equal(t, '≠', cur.Next())
		assert.Equal(t, utfstrings.EOS, cur.Next())
	})

	suite.Run("WithMultiRunes", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "≠∂",
		}

		assert.Equal(t, '≠', cur.Next())
		assert.Equal(t, '∂', cur.Next())
		assert.Equal(t, utfstrings.EOS, cur.Next())
	})
}

func TestCursor_Peek(suite *testing.T) {
	suite.Parallel()

	suite.Run("WithEmptyString", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "",
		}

		assert.Equal(t, utfstrings.EOS, cur.Peek())
	})

	suite.Run("WithOneRune", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "∑",
		}

		assert.Equal(t, '∑', cur.Peek())
		assert.Equal(t, '∑', cur.Peek())
	})

	suite.Run("WithMultipleRune", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "ƒø",
		}

		assert.Equal(t, 'ƒ', cur.Peek())
		assert.Equal(t, 'ƒ', cur.Peek())
	})
}

func TestCursor_Backup(suite *testing.T) {
	suite.Parallel()

	suite.Run("WithEmptyString", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "",
		}

		assert.Equal(t, utfstrings.EOS, cur.Peek())
		cur.Backup()
		assert.Equal(t, utfstrings.EOS, cur.Peek())
	})

	suite.Run("WithString", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "≈ß",
		}

		assert.Equal(t, '≈', cur.Next())
		cur.Backup()
		assert.Equal(t, '≈', cur.Next())
	})
}

func TestCursor_Build(suite *testing.T) {
	suite.Parallel()

	suite.Run("WithEmptyString", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "",
		}

		cur.Build(func(arg *utfstrings.Cursor) {
			t.Error("move func should not have been called")
		})
	})

	suite.Run("WithSomeString", func(t *testing.T) {
		cur := utfstrings.Cursor{
			String: "abcd∑fg",
		}

		out := cur.Build(func(arg *utfstrings.Cursor) {
			// skips ∑ character
			if arg.Peek() == '∑' {
				arg.Next()
			}
		})

		assert.Equal(t, "abcdfg", out)
	})
}
