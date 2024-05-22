package compiler_test

import (
	"testing"

	"hack/compiler"

	"github.com/stretchr/testify/assert"
)

func Test_tokenizer(t *testing.T) {
	t.Run("Non empty file has more tokens", func(t *testing.T) {
		input := "some text"
		tknzr := compiler.NewTokenizer(input)
		assert.True(t, tknzr.HasMoreTokens())
	})

	t.Run("Empty file has no more tokens", func(t *testing.T) {
		input := ""
		tknzr := compiler.NewTokenizer(input)
		assert.False(t, tknzr.HasMoreTokens())

		err := tknzr.Advance()
		assert.Error(t, err)
	})

	t.Run("Identifying tokenTypes", func(t *testing.T) {
		type testCase struct {
			name      string
			input     string
			tokenType compiler.TokenType
			token     string
		}
		testCases := []testCase{
			{name: "check symbol", input: "{", tokenType: compiler.SYMBOL, token: "{"},
			{name: "check symbol", input: "}", tokenType: compiler.SYMBOL, token: "}"},
			{name: "check int_const", input: "1", tokenType: compiler.INT_CONST},
			{name: "check int_const", input: "12", tokenType: compiler.INT_CONST},
			{name: "check class is keyword", input: "class", tokenType: compiler.KEYWORD},
			{name: "check method is keyword", input: "method", tokenType: compiler.KEYWORD},
			{name: "check function is keyword", input: "function", tokenType: compiler.KEYWORD},
			{name: "check static is keyword", input: "static", tokenType: compiler.KEYWORD},
			{name: "check main is identifier", input: "bar", tokenType: compiler.IDENTIFIER},
			{name: "check main is identifier", input: "foo", tokenType: compiler.IDENTIFIER},
		}
		for _, tc := range testCases {
			tknzr := compiler.NewTokenizer(tc.input)

			tt := tknzr.TokenType()
			assert.Equal(t, compiler.TokenType(""), tt)

			for tknzr.HasMoreTokens() {
				err := tknzr.Advance()
				assert.NoError(t, err)

				gotTokenType := tknzr.TokenType()
				assert.Equal(t, tc.tokenType, gotTokenType)

				if gotTokenType == compiler.SYMBOL {
					gotSymbol := tknzr.Symbol()
					assert.Equal(t, tc.token, gotSymbol)
				}
			}
		}
	})

}
