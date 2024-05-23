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

	t.Run("Identifying symbols", func(t *testing.T) {
		t.Run("single symbols", func(t *testing.T) {
			type testCase struct {
				input         string
				wantSymbol    rune
				wantTokenType compiler.TokenType
			}
			testCases := []testCase{
				{input: "{", wantSymbol: '{', wantTokenType: compiler.SYMBOL},
				{input: "}", wantSymbol: '}', wantTokenType: compiler.SYMBOL},
			}
			for _, tc := range testCases {
				tknzr := compiler.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, compiler.TokenType(""), tt)

				for tknzr.HasMoreTokens() {
					err := tknzr.Advance()
					assert.NoError(t, err)

					gotTokenType := tknzr.TokenType()
					assert.Equal(t, tc.wantTokenType, gotTokenType)

					gotToken := tknzr.Symbol()
					assert.Equal(t, tc.wantSymbol, gotToken)
				}
			}
		})

		t.Run("multiple symbols", func(t *testing.T) {
			input := `   {}
  [] (
	)-  +~`
			tknzr := compiler.NewTokenizer(input)

			type testCase struct {
				wantSymbol    rune
				wantTokenType compiler.TokenType
			}
			testCases := []testCase{
				{wantSymbol: '{', wantTokenType: compiler.SYMBOL},
				{wantSymbol: '}', wantTokenType: compiler.SYMBOL},
				{wantSymbol: '[', wantTokenType: compiler.SYMBOL},
				{wantSymbol: ']', wantTokenType: compiler.SYMBOL},
				{wantSymbol: '(', wantTokenType: compiler.SYMBOL},
				{wantSymbol: ')', wantTokenType: compiler.SYMBOL},
				{wantSymbol: '-', wantTokenType: compiler.SYMBOL},
				{wantSymbol: '+', wantTokenType: compiler.SYMBOL},
				{wantSymbol: '~', wantTokenType: compiler.SYMBOL},
			}
			for _, tc := range testCases {
				err := tknzr.Advance()
				assert.NoError(t, err)

				gotTokenType := tknzr.TokenType()
				assert.Equal(t, tc.wantTokenType, gotTokenType)

				gotToken := tknzr.Symbol()
				assert.Equal(t, string(tc.wantSymbol), string(gotToken))
			}
		})
	})

	t.Run("Identifying Integers", func(t *testing.T) {
		t.Run("Identifying single integers", func(t *testing.T) {
			type testCase struct {
				name          string
				input         string
				wantInt       int
				wantTokenType compiler.TokenType
			}
			testCases := []testCase{
				{name: "check int_constant 1", input: "1", wantInt: 1, wantTokenType: compiler.INT_CONST},
				{name: "check int_constant 12", input: "12", wantInt: 12, wantTokenType: compiler.INT_CONST},
				{name: "check int_constant 012", input: "012", wantInt: 12, wantTokenType: compiler.INT_CONST},
			}
			for _, tc := range testCases {
				tknzr := compiler.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, compiler.TokenType(""), tt)

				for tknzr.HasMoreTokens() {
					err := tknzr.Advance()
					assert.NoError(t, err)

					gotTokenType := tknzr.TokenType()
					assert.Equal(t, tc.wantTokenType, gotTokenType)

					gotToken := tknzr.IntVal()
					assert.Equal(t, tc.wantInt, gotToken)
				}
			}
		})

		t.Run("Multiple integers", func(t *testing.T) {
			input := `  123 789
			0099 987
		111  `
			tknzr := compiler.NewTokenizer(input)

			type testCase struct {
				wantInt       int
				wantTokenType compiler.TokenType
			}
			testCases := []testCase{
				{wantInt: 1237890099987111, wantTokenType: compiler.INT_CONST},
			}

			for _, tc := range testCases {
				err := tknzr.Advance()
				assert.NoError(t, err)

				gotTokenType := tknzr.TokenType()
				assert.Equal(t, tc.wantTokenType, gotTokenType)

				gotToken := tknzr.IntVal()
				assert.Equal(t, tc.wantInt, gotToken)
			}
		})
	})

	t.Run("Identifying Keywords", func(t *testing.T) {
		type testCase struct {
			name          string
			input         string
			wantKeyword   compiler.Keyword
			wantTokenType compiler.TokenType
		}
		testCases := []testCase{
			{name: "check class is keyword", input: "class", wantKeyword: compiler.CLASS, wantTokenType: compiler.KEYWORD},
			{name: "check method is keyword", input: "method", wantKeyword: compiler.METHOD, wantTokenType: compiler.KEYWORD},
			{name: "check function is keyword", input: "function", wantKeyword: compiler.FUNCTION, wantTokenType: compiler.KEYWORD},
			{name: "check static is keyword", input: "static", wantKeyword: compiler.STATIC, wantTokenType: compiler.KEYWORD},
		}
		for _, tc := range testCases {
			tknzr := compiler.NewTokenizer(tc.input)

			tt := tknzr.TokenType()
			assert.Equal(t, compiler.TokenType(""), tt)

			for tknzr.HasMoreTokens() {
				err := tknzr.Advance()
				assert.NoError(t, err)

				gotTokenType := tknzr.TokenType()
				assert.Equal(t, tc.wantTokenType, gotTokenType)

				gotToken := tknzr.Keyword()
				assert.Equal(t, tc.wantKeyword, gotToken)
			}
		}
	})

	t.Run("Identifying Identifiers", func(t *testing.T) {
		t.Run("sinlge identifiers", func(t *testing.T) {
			type testCase struct {
				name           string
				input          string
				wantIdentifier string
				wantTokenType  compiler.TokenType
			}
			testCases := []testCase{
				{name: "check foo is identifier", input: "foo", wantIdentifier: "foo", wantTokenType: compiler.IDENTIFIER},
				{name: "check bar is identifier", input: "bar", wantIdentifier: "bar", wantTokenType: compiler.IDENTIFIER},
			}
			for _, tc := range testCases {
				tknzr := compiler.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, compiler.TokenType(""), tt)

				for tknzr.HasMoreTokens() {
					err := tknzr.Advance()
					assert.NoError(t, err)

					gotTokenType := tknzr.TokenType()
					assert.Equal(t, tc.wantTokenType, gotTokenType)

					gotToken := tknzr.Identifier()
					assert.Equal(t, tc.wantIdentifier, gotToken)
				}
			}
		})

		t.Run("Multiple identifiers", func(t *testing.T) {
			input := `  first
	string
constant
`
			tknzr := compiler.NewTokenizer(input)

			type testCase struct {
				wantIdentifier string
				wantTokenType  compiler.TokenType
			}
			testCases := []testCase{
				{wantIdentifier: "first string constant", wantTokenType: compiler.IDENTIFIER},
				{wantIdentifier: "second string constant", wantTokenType: compiler.IDENTIFIER},
				{wantIdentifier: "third string constant", wantTokenType: compiler.IDENTIFIER},
			}

			for _, tc := range testCases {
				err := tknzr.Advance()
				assert.NoError(t, err)

				gotTokenType := tknzr.TokenType()
				assert.Equal(t, tc.wantTokenType, gotTokenType)

				gotToken := tknzr.Identifier()
				assert.Equal(t, tc.wantIdentifier, gotToken)
			}
		})
	})

	t.Run("Identifying String constants", func(t *testing.T) {
		t.Run("single string constants", func(t *testing.T) {
			type testCase struct {
				name          string
				input         string
				wantStringVal string
				wantTokenType compiler.TokenType
			}
			testCases := []testCase{
				{name: "check recognition of string constant", input: `"A string constant"`, wantStringVal: "A string constant", wantTokenType: compiler.STRING_CONST},
			}
			for _, tc := range testCases {
				tknzr := compiler.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, compiler.TokenType(""), tt)

				for tknzr.HasMoreTokens() {
					err := tknzr.Advance()
					assert.NoError(t, err)

					gotTokenType := tknzr.TokenType()
					assert.Equal(t, tc.wantTokenType, gotTokenType)

					gotToken := tknzr.StringVal()
					assert.Equal(t, tc.wantStringVal, gotToken)
				}
			}
		})

		t.Run("Multiple string constants", func(t *testing.T) {
			input := `  "first string constant"
	"second string constant"
		"third string constant"
		"fourth string
constant"
  `
			tknzr := compiler.NewTokenizer(input)

			type testCase struct {
				wantStringVal string
				wantTokenType compiler.TokenType
			}
			testCases := []testCase{
				{wantStringVal: "first string constant", wantTokenType: compiler.STRING_CONST},
				{wantStringVal: "second string constant", wantTokenType: compiler.STRING_CONST},
				{wantStringVal: "third string constant", wantTokenType: compiler.STRING_CONST},
				{wantStringVal: `fourth string
constant`, wantTokenType: compiler.STRING_CONST},
			}

			for _, tc := range testCases {
				err := tknzr.Advance()
				assert.NoError(t, err)

				gotTokenType := tknzr.TokenType()
				assert.Equal(t, tc.wantTokenType, gotTokenType)

				gotToken := tknzr.StringVal()
				assert.Equal(t, tc.wantStringVal, gotToken)
			}
		})
	})
}
