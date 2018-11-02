package lexer_test

import (
	"log"
	"testing"

	"github.com/spy16/parens/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexer_NextToken_EmptyList(suite *testing.T) {
	src := `()`
	lxr := getLexer(src)

	token, err := lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, token.Type, lexer.LPAREN)
	assert.Equal(suite, 0, token.Start)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, token.Type, lexer.RPAREN)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.NotNil(suite, err)
	assert.Nil(suite, token)
	assert.Equal(suite, err, lexer.ErrEOF)
}

func TestLexer_NextToken_NestedList(suite *testing.T) {
	src := `((get-func add) 1 2 3 4)`
	lxr := getLexer(src)

	// consume '('
	token, err := lxr.NextToken()
	assert.Nil(suite, err)
	require.NotNil(suite, token)

	// consume '('
	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	require.NotNil(suite, token)

	// consume 'get-func'
	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	require.NotNil(suite, token)
	assert.Equal(suite, lexer.SYMBOL, token.Type)
	assert.Equal(suite, "get-func", token.Value)
}

func TestLexer_NextToken_Symbols(suite *testing.T) {
	src := `(-> 10)`

	lxr := getLexer(src)

	token, err := lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	require.NotNil(suite, token)
	assert.Equal(suite, lexer.SYMBOL, token.Type)
	assert.Equal(suite, "->", token.Value)
}

func TestLexer_NextToken_Numner(suite *testing.T) {
	suite.Parallel()

	suite.Run("ValidInt", func(t *testing.T) {
		token, err := getLexer(`100`).NextToken()
		assert.Nil(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, token.Type, lexer.NUMBER)
	})

	suite.Run("ValidNegativeInt", func(t *testing.T) {
		token, err := getLexer(`-100`).NextToken()
		assert.Nil(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, token.Type, lexer.NUMBER)
	})

	suite.Run("ValidFloat", func(t *testing.T) {
		token, err := getLexer(`100.1234`).NextToken()
		assert.Nil(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, token.Type, lexer.NUMBER)
	})

	suite.Run("ValidFloat", func(t *testing.T) {
		token, err := getLexer(`100.`).NextToken()
		assert.Nil(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, token.Type, lexer.NUMBER)
	})

	suite.Run("ValidNegativeFloat", func(t *testing.T) {
		token, err := getLexer(`-100.123`).NextToken()
		assert.Nil(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, token.Type, lexer.NUMBER)
	})

	suite.Run("ValidFloat", func(t *testing.T) {
		token, err := getLexer(`.100.`).NextToken()
		assert.Nil(t, err)
		assert.NotNil(t, token)
		assert.Equal(t, token.Type, lexer.NUMBER)
	})

	suite.Run("InvalidNumber", func(t *testing.T) {
		token, err := getLexer(`100x1234`).NextToken()
		assert.NotNil(t, err)
		assert.Nil(t, token)
	})

}

func TestLexer_NextToken_DoubleQuote_String(suite *testing.T) {
	suite.Parallel()

	suite.Run("Valid", func(t *testing.T) {
		token, err := getLexer(`"hello world"`).NextToken()
		assert.Nil(t, err)
		assert.NotNil(t, token)
	})

	suite.Run("MissingEndQuote", func(t *testing.T) {
		token, err := getLexer(`"hello world`).NextToken()
		assert.NotNil(t, err)
		assert.Nil(t, token)
	})

	suite.Run("EscapeDoubleQuote", func(t *testing.T) {
		token, err := getLexer(`"hello\" world"`).NextToken()
		assert.Nil(t, err)
		require.NotNil(t, token)
		assert.Contains(t, token.Value, "world")
	})
}

func TestLexer_NextToken_All(suite *testing.T) {
	src := `(println "hello world" 10 1.3)`

	lxr := lexer.New([]byte(src))
	token, err := lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.LPAREN, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.SYMBOL, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.WHITESPACE, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.DSTRING, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.WHITESPACE, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.NUMBER, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.WHITESPACE, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.NUMBER, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))

	token, err = lxr.NextToken()
	assert.Nil(suite, err)
	assert.NotNil(suite, token)
	assert.Equal(suite, lexer.RPAREN, token.Type)
	log.Printf("got token: '%s', value='%s'", token, string(token.Value))
}

func TestLexer_Tokens(t *testing.T) {
	src := `(println "hello world" 10 1.3)`
	lxr := getLexer(src)
	tokens, err := lxr.Tokens()
	assert.Nil(t, err)
	require.NotNil(t, tokens)
	assert.Equal(t, 9, len(tokens))
}

func getLexer(src string) *lexer.Lexer {
	return lexer.New([]byte(src))
}
