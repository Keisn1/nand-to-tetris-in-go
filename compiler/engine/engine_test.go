package engine_test

import (
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"hack/compiler/engine"
	"hack/compiler/token"
)

func Test_compileClass(t *testing.T) {
	t.Run("Testing happy classes", func(t *testing.T) {
		type testCase struct {
			name string
			fp   string
		}

		dir := "../test_programs/project_11/engine/own/classes/"
		testCases := []testCase{
			{name: "empty main", fp: "emptyMain"},
			{name: "one static class variable", fp: "MainWith1StaticClassVarDec"},
			// {name: "two static class variable", fp: "MainWith2StaticClassVarDec"},
			// {name: "two static class variable different type", fp: "MainWith2StaticClassVarDec2Types"},
			// {name: "main arbitrary class variable declarations", fp: "MainWith3ClassVarDec"},
			// {name: "main with empty subroutine", fp: "MainWithEmptySubroutine"},
			// {name: "main with mulitple subroutines", fp: "main2Subroutine"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.fp+".jack")
				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

				tknzr := token.NewTokenizer(string(input))
				e := engine.NewEngine(&tknzr)

				e.Tknzr.Advance()
				got := e.CompileClass()

				assert.Empty(t, e.Errors)
				assert.Equal(t, want, removeWhiteSpaces(got))
				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
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
				wantErr: engine.NewErrSyntaxUnexpectedToken(token.CLASS, token.EOF),
			},
			{
				inputs:  []string{"var"},
				wantErr: engine.NewErrSyntaxUnexpectedToken(token.CLASS, token.VAR),
			},
			{
				inputs:  []string{"class var"},
				wantErr: engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.VAR),
			},
			{
				inputs:  []string{"class name name"},
				wantErr: engine.NewErrSyntaxUnexpectedToken(token.LBRACE, "name"),
			},
			{
				inputs:  []string{"class name ."},
				wantErr: engine.NewErrSyntaxUnexpectedToken(token.LBRACE, "."),
			},
			{
				inputs:  []string{"class Main {"},
				wantErr: engine.NewErrSyntaxUnexpectedToken(token.RBRACE, token.EOF),
			},
			{
				inputs:  []string{"class Main { ."},
				wantErr: engine.NewErrSyntaxUnexpectedToken(token.RBRACE, token.DOT),
			},
			{
				inputs:  []string{"class Main {static name"},
				wantErr: engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.EOF),
			},
			{
				inputs:  []string{"class Main {static "},
				wantErr: engine.NewErrSyntaxNotAType(token.EOF),
			},
			{
				inputs:  []string{"class Main {function var"},
				wantErr: engine.NewErrSyntaxNotAType(token.VAR),
			},
			{
				inputs:  []string{"class Main {function ;"},
				wantErr: engine.NewErrSyntaxNotAType(token.SEMICOLON),
			},
		}

		for _, tc := range testCases {
			for _, input := range tc.inputs {
				tknzr := token.NewTokenizer(input)
				e := engine.NewEngine(&tknzr)

				e.Tknzr.Advance()
				e.CompileClass()

				assertErrorFound(t, e.Errors, tc.wantErr)
			}
		}
	})
}

// func Test_classVarDec(t *testing.T) {
// 	t.Run("Testing happy ClassVarDeclarations", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}

// 		dir := "../test_programs/project_10/engine/own/classVarDec/"
// 		testCases := []testCase{
// 			{name: "mulitple static char in one line", fp: "multipleStatic"},
// 			{name: "static and field var declaration", fp: "staticAndField"},
// 		}

// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				input := readFile(t, dir+tc.fp+".jack")
// 				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)

// 				e.Tknzr.Advance()
// 				got := e.CompileClassVarDec()

// 				assert.Empty(t, e.Errors)
// 				assert.Equal(t, want, removeWhiteSpaces(got))
// 				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
// 			})
// 		}
// 	})

// 	t.Run("Testing falsy class variable declarations", func(t *testing.T) {
// 		type testCase struct {
// 			inputs  []string
// 			wantErr error
// 		}
// 		testCases := []testCase{
// 			{
// 				inputs:  []string{"", "   "},
// 				wantErr: engine.NewErrSyntaxNotAClassVarDec(token.EOF),
// 			},
// 			{
// 				inputs:  []string{"var"},
// 				wantErr: engine.NewErrSyntaxNotAClassVarDec(token.VAR),
// 			},
// 			{
// 				inputs:  []string{"static boolean boo{"},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.SEMICOLON, token.LBRACE),
// 			},
// 		}
// 		for _, tc := range testCases {
// 			for _, input := range tc.inputs {
// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)

// 				e.Tknzr.Advance()
// 				e.CompileClassVarDec()

// 				assertErrorFound(t, e.Errors, tc.wantErr)
// 			}
// 		}
// 	})
// }

// func Test_subroutineDec(t *testing.T) {
// 	t.Run("Testing happy subroutineDec", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}

// 		dir := "../test_programs/project_10/engine/own/subRoutineDec/"
// 		testCases := []testCase{
// 			{name: "constructor void no parameter", fp: "constVoidNoParam"},
// 			{name: "function int no parameter", fp: "funcIntNoParam"},
// 			{name: "method char methodName", fp: "methCharMethodNameNoParam"},
// 			{name: "one parameter", fp: "oneParameter"},
// 			{name: "two parameter", fp: "twoParameter"},
// 		}

// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				input := readFile(t, dir+tc.fp+".jack")
// 				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)

// 				e.Tknzr.Advance()
// 				got := e.CompileSubroutineDec()

// 				assert.Empty(t, e.Errors)
// 				assert.Equal(t, want, removeWhiteSpaces(got))
// 				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
// 			})
// 		}
// 	})

// 	t.Run("Testing falsy subroutine declarations", func(t *testing.T) {
// 		type testCase struct {
// 			inputs   []string
// 			wantErrs []error
// 		}
// 		testCases := []testCase{
// 			{
// 				inputs:   []string{"", "   "},
// 				wantErrs: []error{engine.NewErrSyntaxNotASubroutineDec(token.EOF)},
// 			},
// 			{
// 				inputs:   []string{"var"},
// 				wantErrs: []error{engine.NewErrSyntaxNotASubroutineDec(token.VAR)},
// 			},
// 			{
// 				inputs:   []string{"function var"},
// 				wantErrs: []error{engine.NewErrSyntaxNotAType(token.VAR)},
// 			},
// 			{
// 				inputs:   []string{"function int var"},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.VAR)},
// 			},
// 			{
// 				inputs:   []string{"function int name;"},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedToken(token.LPAREN, token.SEMICOLON)},
// 			},
// 			{
// 				inputs:   []string{"function int name(var", "function int name(int x, var"},
// 				wantErrs: []error{engine.NewErrSyntaxNotAType(token.VAR)},
// 			},
// 			{
// 				inputs:   []string{"function int name(int var", "function int name(int name1, int var"},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.VAR)},
// 			},
// 			{
// 				inputs:   []string{"function int name(int name1}"},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedToken(token.RPAREN, token.RBRACE)},
// 			},
// 			{
// 				inputs:   []string{"function int name(int name1)}"},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedToken(token.LBRACE, token.RBRACE)},
// 			},
// 		}
// 		for _, tc := range testCases {
// 			for _, input := range tc.inputs {
// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)

// 				e.Tknzr.Advance()
// 				e.CompileSubroutineDec()

// 				for _, wantErr := range tc.wantErrs {
// 					assertErrorFound(t, e.Errors, wantErr)
// 				}
// 			}
// 		}
// 	})
// }

// func Test_subroutineBody(t *testing.T) {
// 	t.Run("Testing happy subRoutineBodies", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}

// 		dir := "../test_programs/project_10/engine/own/subRoutineBody/"
// 		testCases := []testCase{
// 			{name: "one variable declaration", fp: "oneVarDec"},
// 			{name: "multiple variable declaration", fp: "multVarDec"},
// 			{name: "multiple variable declaration one line", fp: "multVarDec1line"},
// 			{name: "one var + a let statement", fp: "oneLetStatement"},
// 		}

// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				input := readFile(t, dir+tc.fp+".jack")
// 				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)
// 				e.Tknzr.Advance()

// 				got := e.CompileSubroutineBody()

// 				assert.Empty(t, e.Errors)
// 				assert.Equal(t, want, removeWhiteSpaces(got))
// 				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
// 			})
// 		}
// 	})

// 	t.Run("Testing falsy subroutine bodies", func(t *testing.T) {
// 		type testCase struct {
// 			inputs   []string
// 			wantErrs []error
// 		}
// 		testCases := []testCase{
// 			{
// 				inputs:   []string{"", "   "},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedToken(token.LBRACE, token.EOF)},
// 			},
// 		}
// 		for _, tc := range testCases {
// 			for _, input := range tc.inputs {
// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)

// 				e.Tknzr.Advance()
// 				e.CompileSubroutineBody()

// 				for _, wantErr := range tc.wantErrs {
// 					assertErrorFound(t, e.Errors, wantErr)
// 				}
// 			}
// 		}
// 	})
// }

// func Test_VarDec(t *testing.T) {
// 	t.Run("Testing falsy varDeclarations", func(t *testing.T) {
// 		type testCase struct {
// 			inputs   []string
// 			wantErrs []error
// 		}
// 		testCases := []testCase{
// 			{
// 				inputs:   []string{"", "   "},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedToken(token.VAR, token.EOF)},
// 			},
// 			{
// 				inputs:   []string{"var;"},
// 				wantErrs: []error{engine.NewErrSyntaxNotAType(token.SEMICOLON)},
// 			},
// 			{
// 				inputs:   []string{"var int ;"},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.SEMICOLON)},
// 			},
// 			{
// 				inputs:   []string{"var int name1, int"},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.INT)},
// 			},
// 			{
// 				inputs:   []string{"var int name1{ "},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedToken(token.SEMICOLON, token.LBRACE)},
// 			},
// 		}
// 		for _, tc := range testCases {
// 			for _, input := range tc.inputs {
// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)

// 				e.Tknzr.Advance()
// 				e.CompileVarDec()

// 				for _, wantErr := range tc.wantErrs {
// 					assertErrorFound(t, e.Errors, wantErr)
// 				}
// 			}
// 		}
// 	})
// }

// func Test_statements(t *testing.T) {
// 	t.Run("Testing happy statements", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}

// 		dir := "../test_programs/project_10/engine/own/statements/"
// 		testCases := []testCase{
// 			{name: "two let statements", fp: "twoLetStatements"},
// 			{name: "let, do and return statement", fp: "letDoReturn"},
// 			{name: "do statement with dot", fp: "doWithDot"},
// 			{name: "while statement", fp: "while"},
// 			{name: "while with < symbol", fp: "while"},
// 			{name: "if without else statement", fp: "ifWithoutElse"},
// 			{name: "if with else statement", fp: "whileWithLowerThan"},
// 		}
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				input := readFile(t, dir+tc.fp+".jack")
// 				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)
// 				e.Tknzr.Advance()

// 				got := e.CompileStatements()

// 				assert.Empty(t, e.Errors)
// 				assert.Equal(t, want, removeWhiteSpaces(got))
// 				assert.Equal(t, token.EOF, e.Tknzr.Keyword())
// 			})
// 		}
// 	})
// }

// func Test_letStatement(t *testing.T) {
// 	t.Run("Testing happy letStatements", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}

// 		dir := "../test_programs/project_10/engine/own/letStatements/"
// 		testCases := []testCase{
// 			{name: "one array assignment", fp: "arrayAssignment"},
// 			{name: "array assignment 1+2 expression", fp: "arrayAssignmentExpressionLHS"},
// 			{name: "let statement string RHS ", fp: "stringRHS"},
// 			{name: "let statement keywordConst RHS ", fp: "kwConstRHS"},
// 			{name: "let statement variableName RHS ", fp: "varNameRHS"},
// 		}

// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				input := readFile(t, dir+tc.fp+".jack")
// 				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 				tknzr := token.NewTokenizer(input)
// 				engine := engine.NewEngine(&tknzr)
// 				engine.Tknzr.Advance()

// 				got := engine.CompileLetStatement()

// 				assert.Empty(t, engine.Errors)
// 				assert.Equal(t, want, removeWhiteSpaces(got))
// 				assert.Equal(t, engine.Tknzr.Keyword(), token.EOF)
// 			})
// 		}
// 	})

// 	t.Run("Testing falsy letStatements", func(t *testing.T) {
// 		type testCase struct {
// 			inputs   []string
// 			wantErrs []error
// 		}
// 		testCases := []testCase{
// 			{
// 				inputs:   []string{"", "   "},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedToken(token.LET, token.EOF)},
// 			},
// 			{
// 				inputs:   []string{"let var"},
// 				wantErrs: []error{engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.VAR)},
// 			},
// 			{
// 				inputs:   []string{"let arr[;"},
// 				wantErrs: []error{engine.NewErrSyntaxNotATerm(token.SEMICOLON)},
// 			},
// 			{
// 				inputs:   []string{"let x = var;"},
// 				wantErrs: []error{engine.NewErrSyntaxNotAKeywordConst(token.VAR)},
// 			},
// 		}
// 		for _, tc := range testCases {
// 			for _, input := range tc.inputs {
// 				tknzr := token.NewTokenizer(input)
// 				engine := engine.NewEngine(&tknzr)

// 				engine.Tknzr.Advance()
// 				engine.CompileLetStatement()

// 				for _, wantErr := range tc.wantErrs {
// 					assertErrorFound(t, engine.Errors, wantErr)
// 				}
// 			}
// 		}
// 	})
// }

// func Test_term(t *testing.T) {
// 	t.Run("Testing happy terms", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}

// 		dir := "../test_programs/project_10/engine/own/terms/"
// 		testCases := []testCase{
// 			{name: "arr[5+6]", fp: "arr5Plus6"},
// 			{name: "(x+y)", fp: "xPlusYinParenthesis"},
// 			{name: "-x", fp: "minusX"},
// 			{name: "~x", fp: "tildeX"},
// 			{name: "foo()", fp: "foo"},
// 			{name: "foo.bar(true)", fp: "foobar"},
// 			{name: "foo(x, (y + z))", fp: "withExpList"},
// 			{name: "foo.bar(x, (y + z))", fp: "fooBarExpList"},
// 		}

// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				input := readFile(t, dir+tc.fp+".jack")
// 				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 				tknzr := token.NewTokenizer(input)
// 				engine := engine.NewEngine(&tknzr)
// 				engine.Tknzr.Advance()

// 				got := engine.CompileTerm()

// 				assert.Empty(t, engine.Errors)
// 				assert.Equal(t, want, removeWhiteSpaces(got))
// 				assert.Equal(t, engine.Tknzr.Keyword(), token.EOF)
// 			})
// 		}
// 	})

// 	t.Run("Testing falsy terms", func(t *testing.T) {
// 		type testCase struct {
// 			inputs   []string
// 			wantErrs []error
// 		}
// 		testCases := []testCase{
// 			{
// 				inputs:   []string{"", "   "},
// 				wantErrs: []error{engine.NewErrSyntaxNotATerm(token.EOF)},
// 			},
// 		}
// 		for _, tc := range testCases {
// 			for _, input := range tc.inputs {
// 				tknzr := token.NewTokenizer(input)
// 				engine := engine.NewEngine(&tknzr)

// 				engine.Tknzr.Advance()
// 				engine.CompileTerm()

// 				for _, wantErr := range tc.wantErrs {
// 					assertErrorFound(t, engine.Errors, wantErr)
// 				}
// 			}
// 		}
// 	})
// }

// func Test_Return(t *testing.T) {
// 	t.Run("Testing falsy return statements", func(t *testing.T) {
// 		type testCase struct {
// 			inputs  []string
// 			wantErr error
// 		}
// 		testCases := []testCase{
// 			{
// 				inputs:  []string{"", "   "},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.RETURN, token.EOF),
// 			},
// 			{
// 				inputs:  []string{"return{"},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.SEMICOLON, token.LBRACE),
// 			},
// 		}
// 		for _, tc := range testCases {
// 			for _, input := range tc.inputs {
// 				tknzr := token.NewTokenizer(input)
// 				engine := engine.NewEngine(&tknzr)

// 				engine.Tknzr.Advance()
// 				engine.CompileReturn()

// 				assertErrorFound(t, engine.Errors, tc.wantErr)
// 			}
// 		}
// 	})
// }

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
	t.Helper()
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
