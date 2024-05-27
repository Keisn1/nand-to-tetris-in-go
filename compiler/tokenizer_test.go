package compiler_test

import (
	"os"
	"regexp"
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
				wantSymbol    string
				wantTokenType compiler.TokenType
			}
			testCases := []testCase{
				{input: "{", wantSymbol: compiler.LBRACE, wantTokenType: compiler.SYMBOL},
				{input: "}", wantSymbol: compiler.RBRACE, wantTokenType: compiler.SYMBOL},
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
	)-  +~ }
}`
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
				{wantSymbol: '}', wantTokenType: compiler.SYMBOL},
				{wantSymbol: '}', wantTokenType: compiler.SYMBOL},
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
				{wantInt: 123, wantTokenType: compiler.INT_CONST},
				{wantInt: 789, wantTokenType: compiler.INT_CONST},
				{wantInt: 99, wantTokenType: compiler.INT_CONST},
				{wantInt: 987, wantTokenType: compiler.INT_CONST},
				{wantInt: 111, wantTokenType: compiler.INT_CONST},
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

		t.Run("Multiple keywords", func(t *testing.T) {
			input := `  class
			method
		function 		static
		`
			tknzr := compiler.NewTokenizer(input)

			type testCase struct {
				wantKeyword   string
				wantTokenType compiler.TokenType
			}
			testCases := []testCase{
				{wantKeyword: compiler.CLASS, wantTokenType: compiler.KEYWORD},
				{wantKeyword: compiler.METHOD, wantTokenType: compiler.KEYWORD},
				{wantKeyword: compiler.FUNCTION, wantTokenType: compiler.KEYWORD},
				{wantKeyword: compiler.STATIC, wantTokenType: compiler.KEYWORD},
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
				wantTokenType  compiler.TokenType
			}
			testCases := []testCase{
				{name: "check Main is identifier", input: "Main", wantIdentifier: "Main", wantTokenType: compiler.IDENTIFIER},
				{name: "check foo is identifier", input: "foo", wantIdentifier: "foo", wantTokenType: compiler.IDENTIFIER},
				{name: "check bar is identifier", input: "bar", wantIdentifier: "bar", wantTokenType: compiler.IDENTIFIER},
				{name: "check c2 is identifier", input: "c2;", wantIdentifier: "c2", wantTokenType: compiler.IDENTIFIER},
			}
			for _, tc := range testCases {
				tknzr := compiler.NewTokenizer(tc.input)

				tt := tknzr.TokenType()
				assert.Equal(t, compiler.TokenType(""), tt)

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
			tknzr := compiler.NewTokenizer(input)

			type testCase struct {
				wantIdentifier string
				wantTokenType  compiler.TokenType
			}
			testCases := []testCase{
				{wantIdentifier: "first", wantTokenType: compiler.IDENTIFIER},
				{wantIdentifier: "string", wantTokenType: compiler.IDENTIFIER},
				{wantIdentifier: "constant", wantTokenType: compiler.IDENTIFIER},
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
		tknzr := compiler.NewTokenizer(input)
		type testCase struct {
			wantTokenType  compiler.TokenType
			wantKeyword    string
			wantSymbol     string
			wantIdentifier string
			wantStringVal  string
			wantIntVal     int
		}
		testCases := []testCase{
			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.CLASS},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "Main"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.LBRACE},

			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.FUNCTION},
			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.VOID},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "main"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.LPAREN},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.RPAREN},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.LBRACE},

			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.VAR},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "Array"},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "a"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.SEMICOLON},

			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.VAR},
			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.INT},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "length"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.SEMICOLON},

			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.VAR},
			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.INT},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "i"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.KOMMA},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "sum"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.SEMICOLON},

			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.LET},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "sum"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.EQUAL},
			{wantTokenType: compiler.INT_CONST, wantIntVal: 1},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.SEMICOLON},

			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.LET},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "length"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.EQUAL},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "Keyboard"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.POINT},
			{wantTokenType: compiler.IDENTIFIER, wantIdentifier: "readInt"},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.LPAREN},
			{wantTokenType: compiler.STRING_CONST, wantStringVal: "HOW MANY NUMBERS? "},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.RPAREN},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.SEMICOLON},

			{wantTokenType: compiler.KEYWORD, wantKeyword: compiler.RETURN},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.SEMICOLON},

			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.RBRACE},
			{wantTokenType: compiler.SYMBOL, wantSymbol: compiler.RBRACE},
		}

		for _, tc := range testCases {
			err := tknzr.Advance()
			assert.NoError(t, err)

			gotTokenType := tknzr.TokenType()
			assert.Equal(t, tc.wantTokenType, gotTokenType)

			switch gotTokenType {
			case compiler.KEYWORD:
				gotKeyword := tknzr.Keyword()
				assert.Equal(t, tc.wantKeyword, gotKeyword)
			case compiler.SYMBOL:
				gotSymbol := tknzr.Symbol()
				assert.Equal(t, tc.wantSymbol, gotSymbol)
			case compiler.IDENTIFIER:
				gotIdent := tknzr.Identifier()
				assert.Equal(t, tc.wantIdentifier, gotIdent)
			case compiler.INT_CONST:
				gotIntVal := tknzr.IntVal()
				assert.Equal(t, tc.wantIntVal, gotIntVal)
			case compiler.STRING_CONST:
				gotStrVal := tknzr.StringVal()
				assert.Equal(t, tc.wantStringVal, gotStrVal)
			default:
				t.Fatalf("default case, no tokenType hit")
			}
		}

		err := tknzr.Advance()
		assert.Error(t, err)

		gotTokenType := tknzr.TokenType()
		assert.Equal(t, compiler.EOF_CONST, gotTokenType)

		got := tknzr.EOF()
		assert.Equal(t, compiler.EOF, got)

	})
}
func Test_tokenizerFullPrograms1(t *testing.T) {
	t.Run("Test output to xml (file without comments)", func(t *testing.T) {
		fp := "test_programs/ArrayTest/Main_wo_comments.jack"
		input, err := os.ReadFile(fp)
		assert.NoError(t, err, "error reading file")

		tknzr := compiler.NewTokenizer(string(input))

		out := "MainT_test.xml"
		tknzr.OutputXML(out)

		gotByte, err := os.ReadFile(out)
		assert.NoError(t, err, "error reading file")

		wantFP := "test_programs/ArrayTest/MainT.xml"
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
				inputFile:   "test_programs/ArrayTest/Main.jack",
				compareFile: "test_programs/ArrayTest/MainT.xml",
			},
			{
				inputFile:   "test_programs/ExpressionLessSquare/Main.jack",
				compareFile: "test_programs/ExpressionLessSquare/MainT.xml",
			},
			{
				inputFile:   "test_programs/ExpressionLessSquare/Square.jack",
				compareFile: "test_programs/ExpressionLessSquare/SquareT.xml",
			},
			{
				inputFile:   "test_programs/ExpressionLessSquare/SquareGame.jack",
				compareFile: "test_programs/ExpressionLessSquare/SquareGameT.xml",
			},
			{
				inputFile:   "test_programs/Square/Main.jack",
				compareFile: "test_programs/Square/MainT.xml",
			},
			{
				inputFile:   "test_programs/Square/Square.jack",
				compareFile: "test_programs/Square/SquareT.xml",
			},
			{
				inputFile:   "test_programs/Square/SquareGame.jack",
				compareFile: "test_programs/Square/SquareGameT.xml",
			},
		}
		for _, tc := range testCases {

			input, err := os.ReadFile(tc.inputFile)
			assert.NoError(t, err, "error reading file")

			tknzr := compiler.NewTokenizer(string(input))

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
