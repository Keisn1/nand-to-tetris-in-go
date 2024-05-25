package compiler_test

import (
	"hack/compiler"
	"os"
	"testing"

	"regexp"

	"github.com/stretchr/testify/assert"
)

func Test_compileClass(t *testing.T) {
	t.Run("Testing happy classes", func(t *testing.T) {
		type testCase struct {
			name string
			fp   string
		}

		dir := "test_programs/own/classes/"
		testCases := []testCase{
			{name: "empty main", fp: "emptyMain"},
			{name: "one static class variable", fp: "MainWith1StaticClassVarDec"},
			{name: "two static class variable", fp: "MainWith2StaticClassVarDec"},
			{name: "two static class variable different type", fp: "MainWith2StaticClassVarDec2Types"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.fp+".jack")
				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

				tknzr := compiler.NewTokenizer(string(input))
				engine := compiler.NewEngine(&tknzr)

				for tknzr.HasMoreTokens() {
					tknzr.Advance()
					switch tknzr.TokenType() {
					case compiler.KEYWORD:
						switch tknzr.Keyword() {
						case compiler.CLASS:
							got, _ := engine.CompileClass()
							assert.Equal(t, want, removeWhiteSpaces(got))
						}
					}
				}
			})
		}
	})

	t.Run("Testing falsy class statements", func(t *testing.T) {
		type testCase struct {
			inputs []string
			error  string
		}
		testCases := []testCase{
			{inputs: []string{"", "var", "   "}, error: "expected keyword CLASS"},
			{inputs: []string{"class var"}, error: "expected tokenType IDENTIFIER"},
			{inputs: []string{"class name name", "class name ."}, error: "expected symbol {"},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(string(input))
				engine := compiler.NewEngine(&tknzr)
				tknzr.Advance()
				_, err := engine.CompileClass()
				assert.ErrorContains(t, err, tc.error)
			}
		}
	})
}

func Test_classVarDev(t *testing.T) {

	t.Run("Testing happy ClassVarDeclarations", func(t *testing.T) {
		type testCase struct {
			name string
			fp   string
		}

		dir := "test_programs/own/classVarDec/"
		testCases := []testCase{
			// {name: "mulitple static char in one line", fp: "multipleStatic"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.fp+".jack")
				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

				tknzr := compiler.NewTokenizer(string(input))
				engine := compiler.NewEngine(&tknzr)
				tknzr.Advance()

				got, err := engine.CompileClassVarDec()
				assert.NoError(t, err)
				assert.Equal(t, want, removeWhiteSpaces(got))
			})
		}
	})

	t.Run("Testing  falsy class variable declarations", func(t *testing.T) {
		type testCase struct {
			inputs []string
			error  string
		}
		testCases := []testCase{
			{inputs: []string{"", "var", "   "}, error: "expected keyword STATIC or FIELD"},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(string(input))
				engine := compiler.NewEngine(&tknzr)
				tknzr.Advance()
				_, err := engine.CompileClassVarDec()
				assert.ErrorContains(t, err, tc.error)
			}
		}
	})
}

func readFile(t *testing.T, fp string) string {
	t.Helper()
	input, err := os.ReadFile(fp)
	assert.NoError(t, err, "error reading file")
	return string(input)
}

func removeWhiteSpaces(input string) string {
	regex := regexp.MustCompile(`\s+`)
	return regex.ReplaceAllString(input, "")
}
