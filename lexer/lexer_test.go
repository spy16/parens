package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spy16/parens/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexer(suite *testing.T) {
	src := `
(begin
	(print "hello world"))`

	checkValidTokens(suite, strings.TrimSpace(src),
		result{lexer.LPAREN, "("},
		result{lexer.SYMBOL, "begin"},
		result{lexer.NEWLINE, "\n"},
		result{lexer.WHITESPACE, "\t"},
		result{lexer.LPAREN, "("},
		result{lexer.SYMBOL, "print"},
		result{lexer.WHITESPACE, " "},
		result{lexer.STRING, `"hello world"`},
		result{lexer.RPAREN, ")"},
		result{lexer.RPAREN, ")"},
	)
}

func TestLexer_Parens(suite *testing.T) {
	suite.Parallel()

	suite.Run("EmptyParens", func(t *testing.T) {
		src := `()`
		checkValidTokens(t, src,
			result{lexer.LPAREN, "("},
			result{lexer.RPAREN, ")"},
		)
	})
	suite.Run("ValidList", func(t *testing.T) {
		checkValidTokens(t, "(hello)",
			result{lexer.LPAREN, "("},
			result{lexer.SYMBOL, "hello"},
			result{lexer.RPAREN, ")"},
		)
	})
}

func TestLexer_Braces(suite *testing.T) {
	suite.Parallel()

	suite.Run("EmptyVector", func(t *testing.T) {
		checkValidTokens(t, "[]",
			result{lexer.LVECT, "["},
			result{lexer.RVECT, "]"},
		)
	})

	suite.Run("ValidVector", func(t *testing.T) {
		checkValidTokens(t, "[1 -2.3]",
			result{lexer.LVECT, "["},
			result{lexer.NUMBER, "1"},
			result{lexer.WHITESPACE, " "},
			result{lexer.NUMBER, "-2.3"},
			result{lexer.RVECT, "]"},
		)
	})
}

func TestLexer_CurlyBraces(suite *testing.T) {
	suite.Parallel()

	suite.Run("EmptyMap", func(t *testing.T) {
		checkValidTokens(t, "{}",
			result{lexer.LDICT, "{"},
			result{lexer.RDICT, "}"},
		)
	})

	suite.Run("ValidVector", func(t *testing.T) {
		checkValidTokens(t, "{:a 1}",
			result{lexer.LDICT, "{"},
			result{lexer.KEYWORD, ":a"},
			result{lexer.WHITESPACE, " "},
			result{lexer.NUMBER, "1"},
			result{lexer.RDICT, "}"},
		)
	})
}

func TestLexer_WhiteSpaces(suite *testing.T) {
	suite.Parallel()

	suite.Run("SingleWhitespace", func(t *testing.T) {
		checkValidTokens(t, " ",
			result{lexer.WHITESPACE, " "},
		)
		checkValidTokens(t, "\t",
			result{lexer.WHITESPACE, "\t"},
		)
	})

	suite.Run("SingleNewline", func(t *testing.T) {
		checkValidTokens(t, "\n",
			result{lexer.NEWLINE, "\n"},
		)

		checkValidTokens(t, "\r",
			result{lexer.NEWLINE, "\r"},
		)
	})

	suite.Run("MixedWhitespaces", func(t *testing.T) {
		checkValidTokens(t, " \t(\nhello)",
			result{lexer.WHITESPACE, " "},
			result{lexer.WHITESPACE, "\t"},
			result{lexer.LPAREN, "("},
			result{lexer.NEWLINE, "\n"},
			result{lexer.SYMBOL, "hello"},
			result{lexer.RPAREN, ")"},
		)
	})
}

func TestLexer_Strings(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleString", func(t *testing.T) {
		checkValidTokens(t, `"hello"`,
			result{lexer.STRING, `"hello"`},
		)
		checkValidTokens(t, `"hello world"`,
			result{lexer.STRING, `"hello world"`},
		)
		checkValidTokens(t, "\"hello\tworld\"",
			result{lexer.STRING, "\"hello\tworld\""},
		)
		checkValidTokens(t, `"hello world in \"english\""`,
			result{lexer.STRING, `"hello world in \"english\""`},
		)

		checkValidTokens(t, `(tokenize "\"hello world\"")`,
			result{lexer.LPAREN, "("},
			result{lexer.SYMBOL, "tokenize"},
			result{lexer.WHITESPACE, " "},
			result{lexer.STRING, `"\"hello world\""`},
			result{lexer.RPAREN, ")"},
		)
		checkValidTokens(t, `(tokenize "'hello'")`,
			result{lexer.LPAREN, "("},
			result{lexer.SYMBOL, "tokenize"},
			result{lexer.WHITESPACE, " "},
			result{lexer.STRING, `"'hello'"`},
			result{lexer.RPAREN, ")"},
		)

		checkInvalidTokens(t, `"hello`, nil)
	})

	suite.Run("MultiLineString", func(t *testing.T) {
		src := `"hello
world"`
		checkValidTokens(t, src, result{lexer.STRING, "\"hello\nworld\""})
	})

	suite.Run("StringInList", func(t *testing.T) {
		src := `("hello")`
		checkValidTokens(t, src,
			result{lexer.LPAREN, "("},
			result{lexer.STRING, "\"hello\""},
			result{lexer.RPAREN, ")"},
		)
	})

	suite.Run("StringInVector", func(t *testing.T) {
		src := `["hello"]`
		checkValidTokens(t, src,
			result{lexer.LVECT, "["},
			result{lexer.STRING, "\"hello\""},
			result{lexer.RVECT, "]"},
		)
	})
}

func TestLexer_Numbers(suite *testing.T) {
	suite.Parallel()

	suite.Run("Integer", func(t *testing.T) {
		checkValidTokens(t, "1",
			result{lexer.NUMBER, "1"},
		)

		checkValidTokens(t, "-1",
			result{lexer.NUMBER, "-1"},
		)

		checkValidTokens(t, "-100",
			result{lexer.NUMBER, "-100"},
		)

	})

	suite.Run("Float", func(t *testing.T) {
		checkValidTokens(t, "1.3",
			result{lexer.NUMBER, "1.3"},
		)

		checkValidTokens(t, "-1.35",
			result{lexer.NUMBER, "-1.35"},
		)

		checkValidTokens(t, "100.34",
			result{lexer.NUMBER, "100.34"},
		)

	})

	suite.Run("InvalidNumbers", func(t *testing.T) {
		checkInvalidTokens(t, "1.09.9", nil)
	})
}

func TestLexer_Symbols(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleASCIISymbols", func(t *testing.T) {
		checkValidTokens(t, "hello", result{lexer.SYMBOL, "hello"})
		checkValidTokens(t, "hello-world", result{lexer.SYMBOL, "hello-world"})
		checkValidTokens(t, "is-true?", result{lexer.SYMBOL, "is-true?"})
		checkValidTokens(t, "+", result{lexer.SYMBOL, "+"})
		checkValidTokens(t, "+-=", result{lexer.SYMBOL, "+-="})
	})

	suite.Run("UnicodeSymbols", func(t *testing.T) {
		checkValidTokens(t, "≠", result{lexer.SYMBOL, "≠"})
		checkValidTokens(t, "∂", result{lexer.SYMBOL, "∂"})
	})
}

func TestLexer_Comment(suite *testing.T) {
	suite.Parallel()

	suite.Run("SimpleComments", func(t *testing.T) {
		checkValidTokens(t, `; this is a sample comment`,
			result{lexer.COMMENT, "; this is a sample comment"},
		)
	})

	suite.Run("CommentWithExpressions", func(t *testing.T) {
		src := `(add 10) ; call the add func`
		checkValidTokens(t, src,
			result{lexer.LPAREN, "("},
			result{lexer.SYMBOL, "add"},
			result{lexer.WHITESPACE, " "},
			result{lexer.NUMBER, "10"},
			result{lexer.RPAREN, ")"},
			result{lexer.WHITESPACE, " "},
			result{lexer.COMMENT, "; call the add func"},
		)
	})
}

type result struct {
	typ lexer.TokenType
	val string
}

func checkValidTokens(t *testing.T, src string, results ...result) {
	tokens, err := lexer.New(src).Tokens()
	require.NoError(t, err)

	require.Equal(t, len(results), len(tokens), fmt.Sprintf("got: %s", tokens))
	for i, res := range results {
		assert.Equal(t, res.typ, tokens[i].Type)
		assert.Equal(t, res.val, tokens[i].Value)
	}
}

func checkInvalidTokens(t *testing.T, src string, expectErr error) {
	tokens, err := lexer.New(src).Tokens()
	require.Error(t, err)
	if expectErr != nil {
		assert.Equal(t, expectErr, err)
	}
	assert.Nil(t, tokens)
}
