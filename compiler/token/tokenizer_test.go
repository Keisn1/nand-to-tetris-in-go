package token_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"hack/compiler/token"
)

func TestSpecial(t *testing.T) {
	t.Run("if statement with comment at the end", func(t *testing.T) {
		input := "if (key = 81)  { let exit = true; }     // q key"
		tknzr := token.NewTokenizer(input)

		tknzr.Advance()
		for tknzr.Symbol() != token.RBRACE {
			tknzr.Advance()
		}
		assert.Equal(t, token.RBRACE, tknzr.Symbol())

		assert.False(t, tknzr.HasMoreTokens())
	})
}

func Test_tokenizer(t *testing.T) {
	t.Run("only comment", func(t *testing.T) {
		input := "// a comment"
		tknzr := token.NewTokenizer(input)

		tknzr.Advance()
		assert.False(t, tknzr.HasMoreTokens())

		err := tknzr.Advance()
		assert.Error(t, err)
	})
	t.Run("x // comment", func(t *testing.T) {
		input := "x // a comment"
		tknzr := token.NewTokenizer(input)

		err := tknzr.Advance()
		assert.NoError(t, err)
		assert.Equal(t, "x", tknzr.Identifier())

		assert.False(t, tknzr.HasMoreTokens())

		err = tknzr.Advance()
		assert.False(t, tknzr.HasMoreTokens())
		assert.Error(t, err)

	})

	t.Run("Non empty file has more tokens", func(t *testing.T) {
		input := "some text"
		tknzr := token.NewTokenizer(input)

		assert.True(t, tknzr.HasMoreTokens())
	})

	t.Run("Empty file has no more tokens", func(t *testing.T) {
		input := ""
		tknzr := token.NewTokenizer(input)

		err := tknzr.Advance()
		assert.False(t, tknzr.HasMoreTokens())
		assert.Error(t, err)
	})

	t.Run("Identifying symbols", func(t *testing.T) {
		t.Run("single symbols", func(t *testing.T) {
			type testCase struct {
				input         string
				wantSymbol    string
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{input: "{", wantSymbol: token.LBRACE, wantTokenType: token.SYMBOL},
				{input: "}", wantSymbol: token.RBRACE, wantTokenType: token.SYMBOL},
			}
			for _, tc := range testCases {
				tknzr := token.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, token.TokenType(""), tt)

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
		)-  +~ }
	}`
			tknzr := token.NewTokenizer(input)

			type testCase struct {
				wantSymbol    rune
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{wantSymbol: '{', wantTokenType: token.SYMBOL},
				{wantSymbol: '}', wantTokenType: token.SYMBOL},
				{wantSymbol: '[', wantTokenType: token.SYMBOL},
				{wantSymbol: ']', wantTokenType: token.SYMBOL},
				{wantSymbol: '(', wantTokenType: token.SYMBOL},
				{wantSymbol: ')', wantTokenType: token.SYMBOL},
				{wantSymbol: '-', wantTokenType: token.SYMBOL},
				{wantSymbol: '+', wantTokenType: token.SYMBOL},
				{wantSymbol: '~', wantTokenType: token.SYMBOL},
				{wantSymbol: '}', wantTokenType: token.SYMBOL},
				{wantSymbol: '}', wantTokenType: token.SYMBOL},
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
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{name: "check int_constant 1", input: "1", wantInt: 1, wantTokenType: token.INT_CONST},
				{name: "check int_constant 12", input: "12", wantInt: 12, wantTokenType: token.INT_CONST},
				{name: "check int_constant 012", input: "012", wantInt: 12, wantTokenType: token.INT_CONST},
			}
			for _, tc := range testCases {
				tknzr := token.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, token.TokenType(""), tt)

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
			tknzr := token.NewTokenizer(input)

			type testCase struct {
				wantInt       int
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{wantInt: 123, wantTokenType: token.INT_CONST},
				{wantInt: 789, wantTokenType: token.INT_CONST},
				{wantInt: 99, wantTokenType: token.INT_CONST},
				{wantInt: 987, wantTokenType: token.INT_CONST},
				{wantInt: 111, wantTokenType: token.INT_CONST},
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
		t.Run("Single keywords", func(t *testing.T) {

			type testCase struct {
				name          string
				input         string
				wantKeyword   string
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{name: "check class is keyword", input: "class", wantKeyword: token.CLASS, wantTokenType: token.KEYWORD},
				{name: "check method is keyword", input: "method", wantKeyword: token.METHOD, wantTokenType: token.KEYWORD},
				{name: "check function is keyword", input: "function", wantKeyword: token.FUNCTION, wantTokenType: token.KEYWORD},
				{name: "check static is keyword", input: "static", wantKeyword: token.STATIC, wantTokenType: token.KEYWORD},
			}
			for _, tc := range testCases {
				tknzr := token.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, token.TokenType(""), tt)

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

		t.Run("Multiple keywords", func(t *testing.T) {
			input := `  class
				method
			function 		static
			`
			tknzr := token.NewTokenizer(input)

			type testCase struct {
				wantKeyword   string
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{wantKeyword: token.CLASS, wantTokenType: token.KEYWORD},
				{wantKeyword: token.METHOD, wantTokenType: token.KEYWORD},
				{wantKeyword: token.FUNCTION, wantTokenType: token.KEYWORD},
				{wantKeyword: token.STATIC, wantTokenType: token.KEYWORD},
			}

			for _, tc := range testCases {
				err := tknzr.Advance()
				assert.NoError(t, err)

				gotTokenType := tknzr.TokenType()
				assert.Equal(t, tc.wantTokenType, gotTokenType)

				gotToken := tknzr.Keyword()
				assert.Equal(t, tc.wantKeyword, gotToken)
			}

		})
	})

	t.Run("Identifying Identifiers", func(t *testing.T) {
		t.Run("sinlge identifiers", func(t *testing.T) {
			type testCase struct {
				name           string
				input          string
				wantIdentifier string
				wantTokenType  token.TokenType
			}
			testCases := []testCase{
				{name: "check Main is identifier", input: "Main", wantIdentifier: "Main", wantTokenType: token.IDENTIFIER},
				{name: "check foo is identifier", input: "foo", wantIdentifier: "foo", wantTokenType: token.IDENTIFIER},
				{name: "check bar is identifier", input: "bar", wantIdentifier: "bar", wantTokenType: token.IDENTIFIER},
				{name: "check c2 is identifier", input: "c2;", wantIdentifier: "c2", wantTokenType: token.IDENTIFIER},
			}
			for _, tc := range testCases {
				tknzr := token.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, token.TokenType(""), tt)

				err := tknzr.Advance()
				assert.NoError(t, err)

				gotTokenType := tknzr.TokenType()
				assert.Equal(t, tc.wantTokenType, gotTokenType)

				gotToken := tknzr.Identifier()
				assert.Equal(t, tc.wantIdentifier, gotToken)
			}
		})

		t.Run("Multiple identifiers", func(t *testing.T) {
			input := `  first
				string
			constant
			`
			tknzr := token.NewTokenizer(input)

			type testCase struct {
				wantIdentifier string
				wantTokenType  token.TokenType
			}
			testCases := []testCase{
				{wantIdentifier: "first", wantTokenType: token.IDENTIFIER},
				{wantIdentifier: "string", wantTokenType: token.IDENTIFIER},
				{wantIdentifier: "constant", wantTokenType: token.IDENTIFIER},
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
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{name: "check recognition of string constant", input: `"A string constant"`, wantStringVal: "A string constant", wantTokenType: token.STRING_CONST},
			}
			for _, tc := range testCases {
				tknzr := token.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, token.TokenType(""), tt)

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
			tknzr := token.NewTokenizer(input)
			type testCase struct {
				wantStringVal string
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{wantStringVal: "first string constant", wantTokenType: token.STRING_CONST},
				{wantStringVal: "second string constant", wantTokenType: token.STRING_CONST},
				{wantStringVal: "third string constant", wantTokenType: token.STRING_CONST},
				{wantStringVal: `fourth string
	constant`, wantTokenType: token.STRING_CONST},
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

		t.Run("Multiple string constants", func(t *testing.T) {
			input := `  "first string constant"
		"second string constant"
			"third string constant"
			"fourth string
	constant"
	  `
			tknzr := token.NewTokenizer(input)

			type testCase struct {
				wantStringVal string
				wantTokenType token.TokenType
			}
			testCases := []testCase{
				{wantStringVal: "first string constant", wantTokenType: token.STRING_CONST},
				{wantStringVal: "second string constant", wantTokenType: token.STRING_CONST},
				{wantStringVal: "third string constant", wantTokenType: token.STRING_CONST},
				{wantStringVal: `fourth string
	constant`, wantTokenType: token.STRING_CONST},
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

func Test_tokenizerFullPrograms(t *testing.T) {
	t.Run("Test 1", func(t *testing.T) {
		input := `
class Main {
    function void main() {
        var Array a;
        var int length;
        var int i, sum;


		let sum = 1;

		let length = Keyboard.readInt("HOW MANY NUMBERS? ");

		return;
    }
}
`
		tknzr := token.NewTokenizer(input)
		type testCase struct {
			wantTokenType  token.TokenType
			wantKeyword    string
			wantSymbol     string
			wantIdentifier string
			wantStringVal  string
			wantIntVal     int
		}
		testCases := []testCase{
			{wantTokenType: token.KEYWORD, wantKeyword: token.CLASS},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "Main"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.LBRACE},

			{wantTokenType: token.KEYWORD, wantKeyword: token.FUNCTION},
			{wantTokenType: token.KEYWORD, wantKeyword: token.VOID},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "main"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.LPAREN},
			{wantTokenType: token.SYMBOL, wantSymbol: token.RPAREN},
			{wantTokenType: token.SYMBOL, wantSymbol: token.LBRACE},

			{wantTokenType: token.KEYWORD, wantKeyword: token.VAR},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "Array"},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "a"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.SEMICOLON},

			{wantTokenType: token.KEYWORD, wantKeyword: token.VAR},
			{wantTokenType: token.KEYWORD, wantKeyword: token.INT},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "length"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.SEMICOLON},

			{wantTokenType: token.KEYWORD, wantKeyword: token.VAR},
			{wantTokenType: token.KEYWORD, wantKeyword: token.INT},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "i"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.KOMMA},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "sum"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.SEMICOLON},

			{wantTokenType: token.KEYWORD, wantKeyword: token.LET},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "sum"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.EQUAL},
			{wantTokenType: token.INT_CONST, wantIntVal: 1},
			{wantTokenType: token.SYMBOL, wantSymbol: token.SEMICOLON},

			{wantTokenType: token.KEYWORD, wantKeyword: token.LET},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "length"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.EQUAL},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "Keyboard"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.DOT},
			{wantTokenType: token.IDENTIFIER, wantIdentifier: "readInt"},
			{wantTokenType: token.SYMBOL, wantSymbol: token.LPAREN},
			{wantTokenType: token.STRING_CONST, wantStringVal: "HOW MANY NUMBERS? "},
			{wantTokenType: token.SYMBOL, wantSymbol: token.RPAREN},
			{wantTokenType: token.SYMBOL, wantSymbol: token.SEMICOLON},

			{wantTokenType: token.KEYWORD, wantKeyword: token.RETURN},
			{wantTokenType: token.SYMBOL, wantSymbol: token.SEMICOLON},

			{wantTokenType: token.SYMBOL, wantSymbol: token.RBRACE},
			{wantTokenType: token.SYMBOL, wantSymbol: token.RBRACE},
		}

		for _, tc := range testCases {
			err := tknzr.Advance()
			assert.NoError(t, err)

			gotTokenType := tknzr.TokenType()
			assert.Equal(t, tc.wantTokenType, gotTokenType)

			switch gotTokenType {
			case token.KEYWORD:
				gotKeyword := tknzr.Keyword()
				assert.Equal(t, tc.wantKeyword, gotKeyword)
			case token.SYMBOL:
				gotSymbol := tknzr.Symbol()
				assert.Equal(t, tc.wantSymbol, gotSymbol)
			case token.IDENTIFIER:
				gotIdent := tknzr.Identifier()
				assert.Equal(t, tc.wantIdentifier, gotIdent)
			case token.INT_CONST:
				gotIntVal := tknzr.IntVal()
				assert.Equal(t, tc.wantIntVal, gotIntVal)
			case token.STRING_CONST:
				gotStrVal := tknzr.StringVal()
				assert.Equal(t, tc.wantStringVal, gotStrVal)
			default:
				t.Fatalf("default case, no tokenType hit")
			}
		}

		err := tknzr.Advance()
		assert.Error(t, err)

		gotTokenType := tknzr.TokenType()
		assert.Equal(t, token.EOF_CONST, gotTokenType)
	})
}

func Test_tokenizerFullPrograms1(t *testing.T) {
	t.Run("Test output to xml (file without comments)", func(t *testing.T) {
		fp := "../test_programs/ArrayTest/Main_wo_comments.jack"
		input, err := os.ReadFile(fp)
		assert.NoError(t, err, "error reading file")

		tknzr := token.NewTokenizer(string(input))

		out := "MainT_test.xml"
		tknzr.OutputXML(out)

		gotByte, err := os.ReadFile(out)
		assert.NoError(t, err, "error reading file")

		wantFP := "../test_programs/ArrayTest/MainT.xml"
		wantByte, err := os.ReadFile(wantFP)
		assert.NoError(t, err, "error reading file")

		regex := regexp.MustCompile(`\s+`)
		got := regex.ReplaceAllString(string(gotByte), "")
		want := regex.ReplaceAllString(string(wantByte), "")
		assert.Equal(t, want, got)
		os.Remove(out)
	})

	t.Run("Test output to xml (file with comments)", func(t *testing.T) {
		testCases := []struct {
			inputFile   string
			compareFile string
		}{
			{
				inputFile:   "../test_programs/ArrayTest/Main.jack",
				compareFile: "../test_programs/ArrayTest/MainT.xml",
			},
			{
				inputFile:   "../test_programs/ExpressionLessSquare/Main.jack",
				compareFile: "../test_programs/ExpressionLessSquare/MainT.xml",
			},
			{
				inputFile:   "../test_programs/ExpressionLessSquare/Square.jack",
				compareFile: "../test_programs/ExpressionLessSquare/SquareT.xml",
			},
			{
				inputFile:   "../test_programs/ExpressionLessSquare/SquareGame.jack",
				compareFile: "../test_programs/ExpressionLessSquare/SquareGameT.xml",
			},
			{
				inputFile:   "../test_programs/Square/Main.jack",
				compareFile: "../test_programs/Square/MainT.xml",
			},
			{
				inputFile:   "../test_programs/Square/Square.jack",
				compareFile: "../test_programs/Square/SquareT.xml",
			},
			{
				inputFile:   "../test_programs/Square/SquareGame.jack",
				compareFile: "../test_programs/Square/SquareGameT.xml",
			},
		}
		for _, tc := range testCases {

			input, err := os.ReadFile(tc.inputFile)
			assert.NoError(t, err, "error reading file")

			token.NewTokenizer(string(input))
			tknzr := token.NewTokenizer(string(input))

			out := "MainT_test.xml"
			tknzr.OutputXML(out)

			gotByte, err := os.ReadFile(out)
			assert.NoError(t, err, "error reading file")

			wantByte, err := os.ReadFile(tc.compareFile)
			assert.NoError(t, err, "error reading file")

			regex := regexp.MustCompile(`\s+`)
			got := regex.ReplaceAllString(string(gotByte), "")
			want := regex.ReplaceAllString(string(wantByte), "")
			assert.Equal(t, want, got)
			os.Remove(out)
		}
	})
}
