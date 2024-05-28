package compiler_test

import (
	"errors"
	"fmt"
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
			{name: "main with mulitple empty subroutine", fp: "MainWithMultEmptySubroutine"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.fp+".jack")
				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

				tknzr := compiler.NewTokenizer(string(input))
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				got := engine.CompileClass()
				assert.Equal(t, want, removeWhiteSpaces(got))

				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
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
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.CLASS, compiler.EOF),
			},
			{
				inputs:  []string{"var"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.CLASS, compiler.VAR),
			},
			{
				inputs:  []string{"class var"},
				wantErr: compiler.NewErrSyntaxUnexpectedTokenType(compiler.IDENTIFIER, compiler.VAR),
			},
			{
				inputs:  []string{"class name name"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.LBRACE, "name"),
			},
			{
				inputs:  []string{"class name ."},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.LBRACE, "."),
			},
			{
				inputs:  []string{"class Main {"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.RBRACE, compiler.EOF),
			},
			{
				inputs:  []string{"class Main { ."},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.RBRACE, compiler.POINT),
			},
			{
				inputs:  []string{"class Main {static name"},
				wantErr: compiler.NewErrSyntaxUnexpectedTokenType(compiler.IDENTIFIER, compiler.EOF),
			},
			{
				inputs: []string{"class Main {static "},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(
					fmt.Sprintf("KEYWORD %s / %s / %s or className", compiler.INT, compiler.CHAR, compiler.BOOLEAN),
					compiler.EOF,
				),
			},
			{
				inputs: []string{"class Main {function var"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(
					fmt.Sprintf("KEYWORD %s / %s / %s or className", compiler.INT, compiler.CHAR, compiler.BOOLEAN),
					compiler.VAR,
				),
			},
			{
				inputs: []string{"class Main {function ;"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(
					fmt.Sprintf("KEYWORD %s / %s / %s or className", compiler.INT, compiler.CHAR, compiler.BOOLEAN),
					compiler.SEMICOLON,
				),
			},
		}

		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				engine.CompileClass()

				assertErrorFound(t, engine.Errors, tc.wantErr)
				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
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
				got := engine.CompileClassVarDec()
				assert.Equal(t, want, removeWhiteSpaces(got))

				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
			})
		}
	})

	t.Run("Testing falsy class variable declarations", func(t *testing.T) {
		type testCase struct {
			inputs  []string
			wantErr error
		}
		testCases := []testCase{
			{
				inputs: []string{"", "   "},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(
					fmt.Sprintf("KEYWORD %s or %s ", compiler.STATIC, compiler.FIELD),
					compiler.EOF,
				),
			},
			{
				inputs: []string{"var"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(
					fmt.Sprintf("KEYWORD %s or %s ", compiler.STATIC, compiler.FIELD),
					compiler.VAR,
				),
			},
			{
				inputs:  []string{"static boolean boo{"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.SEMICOLON, compiler.LBRACE),
			},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				engine.CompileClassVarDec()

				assertErrorFound(t, engine.Errors, tc.wantErr)

				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
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
				got := engine.CompileSubroutineDec()

				assert.Equal(t, want, removeWhiteSpaces(got))
				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
			})
		}
	})

	t.Run("Testing falsy subroutine declarations", func(t *testing.T) {
		type testCase struct {
			inputs   []string
			wantErrs []error
		}
		testCases := []testCase{
			{
				inputs: []string{"", "   "},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(
					fmt.Sprintf("KEYWORD %s / %s / %s ", compiler.CONSTRUCTOR, compiler.FUNCTION, compiler.METHOD),
					compiler.EOF,
				)},
			},
			{
				inputs: []string{"var"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(
					fmt.Sprintf("KEYWORD %s / %s / %s ", compiler.CONSTRUCTOR, compiler.FUNCTION, compiler.METHOD),
					compiler.VAR,
				)},
			},
			{
				inputs: []string{"function var"},
				wantErrs: []error{
					compiler.NewErrSyntaxUnexpectedToken(
						fmt.Sprintf("KEYWORD %s / %s / %s or className", compiler.INT, compiler.CHAR, compiler.BOOLEAN),
						compiler.VAR,
					),
					compiler.NewErrSyntaxUnexpectedToken(fmt.Sprintf("expected KEYWORD %s or type", compiler.VOID), compiler.VAR),
				},
			},
			{
				inputs:   []string{"function int var"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedTokenType(compiler.IDENTIFIER, compiler.VAR)},
			},
			{
				inputs:   []string{"function int name;"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(compiler.LPAREN, compiler.SEMICOLON)},
			},
			{
				inputs: []string{"function int name(var", "function int name(int x, var"},
				wantErrs: []error{
					compiler.NewErrSyntaxUnexpectedToken(
						fmt.Sprintf("KEYWORD %s / %s / %s or className", compiler.INT, compiler.CHAR, compiler.BOOLEAN),
						compiler.VAR,
					),
				},
			},

			{
				inputs:   []string{"function int name(int var", "function int name(int name1, int var"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedTokenType(compiler.IDENTIFIER, compiler.VAR)},
			},
			{
				inputs:   []string{"function int name(int name1}"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(compiler.RPAREN, compiler.RBRACE)},
			},
			{
				inputs:   []string{"function int name(int name1)}"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(compiler.LBRACE, compiler.RBRACE)},
			},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				engine.CompileSubroutineDec()

				for _, wantErr := range tc.wantErrs {
					assertErrorFound(t, engine.Errors, wantErr)
				}

				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
			}
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

				got := engine.CompileSubroutineBody()
				assert.Equal(t, want, removeWhiteSpaces(got))

				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
			})
		}
	})

	t.Run("Testing falsy subroutine bodies", func(t *testing.T) {
		type testCase struct {
			inputs   []string
			wantErrs []error
		}
		testCases := []testCase{
			{
				inputs:   []string{"", "   "},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(compiler.LBRACE, compiler.EOF)},
			},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				engine.CompileSubroutineBody()

				for _, wantErr := range tc.wantErrs {
					assertErrorFound(t, engine.Errors, wantErr)
				}

				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
			}
		}
	})
}

func Test_VarDec(t *testing.T) {
	t.Run("Testing falsy varDeclarations", func(t *testing.T) {
		type testCase struct {
			inputs   []string
			wantErrs []error
		}
		testCases := []testCase{
			{
				inputs:   []string{"", "   "},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(compiler.VAR, compiler.EOF)},
			},
			{
				inputs: []string{"var;"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(
					fmt.Sprintf("KEYWORD %s / %s / %s or className", compiler.INT, compiler.CHAR, compiler.BOOLEAN),
					compiler.SEMICOLON,
				)},
			},
			{
				inputs:   []string{"var int ;"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedTokenType(compiler.IDENTIFIER, compiler.SEMICOLON)},
			},
			{
				inputs:   []string{"var int name1, int"},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedTokenType(compiler.IDENTIFIER, compiler.INT)},
			},
			{
				inputs:   []string{"var int name1{ "},
				wantErrs: []error{compiler.NewErrSyntaxUnexpectedToken(compiler.SEMICOLON, compiler.LBRACE)},
			},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				engine.CompileVarDec()

				for _, wantErr := range tc.wantErrs {
					assertErrorFound(t, engine.Errors, wantErr)
				}

				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
			}
		}
	})
}

func Test_Return(t *testing.T) {
	t.Run("Testing falsy return statements", func(t *testing.T) {
		type testCase struct {
			inputs  []string
			wantErr error
		}
		testCases := []testCase{
			{
				inputs:  []string{"", "   "},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.RETURN, compiler.EOF),
			},
			{
				inputs:  []string{"return{"},
				wantErr: compiler.NewErrSyntaxUnexpectedToken(compiler.SEMICOLON, compiler.LBRACE),
			},
		}
		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := compiler.NewTokenizer(input)
				engine := compiler.NewEngine(&tknzr)

				engine.Tknzr.Advance()
				engine.CompileReturn()

				assertErrorFound(t, engine.Errors, tc.wantErr)

				assert.Equal(t, engine.Tknzr.Keyword(), compiler.EOF)
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

func assertErrorFound(t *testing.T, gotErrs []error, wantErr error) {
	foundErr := false
	for _, gotErr := range gotErrs {
		if errors.Is(gotErr, wantErr) {
			foundErr = true
		}
	}

	if !foundErr {
		t.Log("Could not find error:", wantErr, "\nPresent errors: ", gotErrs)
		t.Fail()
	}
}
