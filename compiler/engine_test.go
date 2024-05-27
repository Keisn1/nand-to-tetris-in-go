package compiler_test

import (
	"errors"
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
			{name: "main arbitrary class variable declarations", fp: "MainWith3ClassVarDec"},
			{name: "main with empty subroutine", fp: "MainWithEmptySubroutine"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.fp+".jack")
				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

				tknzr := compiler.NewTokenizer(string(input))
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				got, err := engine.CompileClass()
				assert.NoError(t, err)
				assert.Equal(t, want, removeWhiteSpaces(got))
			})
		}
	})

	t.Run("Testing falsy class statements", func(t *testing.T) {
		type testCase struct {
			inputs  []string
			wantErr error
		}
		testCases := []testCase{
			{
				inputs:  []string{"", "   "},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.NewToken(string(compiler.CLASS), compiler.KEYWORD), compiler.Token{}),
			},
			{
				inputs:  []string{"var"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.NewToken(compiler.CLASS, compiler.KEYWORD), compiler.NewToken(compiler.VAR, compiler.KEYWORD)),
			},
			// {inputs: []string{"class var"}, errors: "compileClass: expected tokenType IDENTIFIER"},
			// {inputs: []string{"class name name", "class name ."}, errors: "compileClass: expected SYMBOL {"},
			// {inputs: []string{"class Main {", "class name { ."}, errors: "compileClass: expected SYMBOL }"},
			// {inputs: []string{"class Main {static name"}, errors: "compileClass: compileClassVarDec"},
			// {inputs: []string{"class Main {function var"}, errors: "compileClass: compileSubroutineDec"},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				engine.CompileClass()

				foundErr := false
				for _, gotErr := range engine.Errors {
					if errors.Is(gotErr, tc.wantErr) {
						foundErr = true
					}
				}

				if !foundErr {
					t.Log("Could not find error:", tc.wantErr, "\nPresent errors: ", engine.Errors)
					t.Fail()
				}
			}
		}
	})

	t.Run("Testing falsy class statements", func(t *testing.T) {
		type testCase struct {
			inputs []string
			error  string
		}
		testCases := []testCase{
			{inputs: []string{"class var"}, error: "compileClass: expected tokenType IDENTIFIER"},
			{inputs: []string{"class name name", "class name ."}, error: "compileClass: expected SYMBOL {"},
			{inputs: []string{"class Main {", "class name { ."}, error: "compileClass: expected SYMBOL }"},
			{inputs: []string{"class Main {static name"}, error: "compileClass: compileClassVarDec"},
			{inputs: []string{"class Main {function var"}, error: "compileClass: compileSubroutineDec"},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				_, err := engine.CompileClass()

				assert.ErrorContains(t, err, tc.error)
			}
		}
	})
}

func Test_classVarDec(t *testing.T) {
	t.Run("Testing happy ClassVarDeclarations", func(t *testing.T) {
		type testCase struct {
			name string
			fp   string
		}

		dir := "test_programs/own/classVarDec/"
		testCases := []testCase{
			{name: "mulitple static char in one line", fp: "multipleStatic"},
			{name: "static and field var declaration", fp: "staticAndField"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.fp+".jack")
				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				got, err := engine.CompileClassVarDec()
				assert.NoError(t, err)
				assert.Equal(t, want, removeWhiteSpaces(got))
			})
		}
	})

	t.Run("Testing falsy class variable declarations", func(t *testing.T) {
		type testCase struct {
			inputs []string
			error  string
		}
		testCases := []testCase{
			{inputs: []string{"", "var", "   "}, error: "expected KEYWORD static or field"},
			{inputs: []string{"static boolean boo{"}, error: "compileClassVarDec: expected SYMBOL ;"},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				_, err := engine.CompileClassVarDec()
				assert.ErrorContains(t, err, tc.error)
			}
		}
	})
}

func Test_subroutineDec(t *testing.T) {
	t.Run("Testing happy subroutineDec", func(t *testing.T) {
		type testCase struct {
			name string
			fp   string
		}

		dir := "test_programs/own/subRoutineDec/"
		testCases := []testCase{
			{name: "constructor void no parameter", fp: "constVoidNoParam"},
			{name: "function int no parameter", fp: "funcIntNoParam"},
			{name: "method char methodName", fp: "methCharMethodNameNoParam"},
			{name: "one parameter", fp: "oneParameter"},
			{name: "two parameter", fp: "twoParameter"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.fp+".jack")
				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				got, err := engine.CompileSubroutineDec()
				assert.NoError(t, err)
				assert.Equal(t, want, removeWhiteSpaces(got))
			})
		}
	})

	t.Run("Testing falsy subroutine declarations", func(t *testing.T) {
		type testCase struct {
			name   string
			inputs []string
			error  []string
		}
		testCases := []testCase{
			{name: "1", inputs: []string{"", "var", "   "}, error: []string{"compileSubroutineDec", "expected KEYWORD constructor / function / method"}},
			{name: "2", inputs: []string{"function var"}, error: []string{
				"compileSubroutineDec",
				"expected type or KEYWORD void",
				"expected className or KEYWORD int / char / boolean"},
			},
			{name: "3", inputs: []string{"function int var"}, error: []string{"compileSubroutineDec", "expected tokenType IDENTIFIER"}},
			{name: "4", inputs: []string{"function int name;"}, error: []string{"compileSubroutineDec", "expected SYMBOL ("}},
			{name: "5", inputs: []string{"function int name()}"}, error: []string{"compileSubroutineDec", "expected SYMBOL {"}},
			{name: "6", inputs: []string{"function int name(var)", "function int name(int x, var)"}, error: []string{"compileSubroutineDec", "compileParameterList", "expected className or KEYWORD int / char / boolean"}},
			{name: "7", inputs: []string{"function int name(int var)"}, error: []string{"compileSubroutineDec", "compileParameterList", "expected tokenType IDENTIFIER"}},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				for _, input := range tc.inputs {
					tknzr := compiler.NewTokenizer(input)
					engine := compiler.NewEngine(&tknzr)

					engine.Tknzr.Advance()
					_, err := engine.CompileSubroutineDec()
					for _, wantErr := range tc.error {
						assert.ErrorContains(t, err, wantErr)
					}
				}
			})
		}
	})
}

func Test_subroutineBody(t *testing.T) {
	t.Run("Testing happy subRoutineBodies", func(t *testing.T) {
		type testCase struct {
			name string
			fp   string
		}

		dir := "test_programs/own/subRoutineBody/"
		testCases := []testCase{
			{name: "one variable declaration", fp: "oneVarDec"},
			{name: "multiple variable declaration", fp: "multVarDec"},
			{name: "multiple variable declaration one line", fp: "multVarDec1line"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.fp+".jack")
				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)
				engine.Tknzr.Advance()

				got, err := engine.CompileSubroutineBody()
				assert.NoError(t, err)
				assert.Equal(t, want, removeWhiteSpaces(got))
			})
		}
	})

	t.Run("Testing falsy subroutine bodies", func(t *testing.T) {
		type testCase struct {
			name   string
			inputs []string
			errors []string
		}
		testCases := []testCase{
			{
				name:   "false varDec",
				inputs: []string{"{var ;}"},
				errors: []string{
					"compileSubroutineBody",
					"compileVarDec",
					"expected className or KEYWORD int / char / boolean",
				},
			},
			{
				name:   "false varDec",
				inputs: []string{"{var ;}"},
				errors: []string{
					"compileSubroutineBody",
					"compileVarDec",
					"expected className or KEYWORD int / char / boolean",
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				for _, input := range tc.inputs {
					tknzr := compiler.NewTokenizer(input)
					engine := compiler.NewEngine(&tknzr)

					engine.Tknzr.Advance()
					_, err := engine.CompileSubroutineBody()
					for _, wantErr := range tc.errors {
						assert.ErrorContains(t, err, wantErr)
					}
				}
			})
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
