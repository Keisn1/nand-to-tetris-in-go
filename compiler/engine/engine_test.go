package engine_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"hack/compiler/engine"
	"hack/compiler/token"
	"hack/compiler/vmWriter"

	"github.com/stretchr/testify/assert"
)

func Test_subroutineBody(t *testing.T) {
	t.Run("Testing happy subRoutineBodies", func(t *testing.T) {
		type testCase struct {
			name     string
			filename string
		}

		dir := "own/bodies/"
		testCases := []testCase{
			{name: "var int x;let x = 2+3;return;", filename: "letStatement2plus3"},
			{name: "var int x;let x = 3+4;return;", filename: "letStatement3plus4"},
			{name: "var int x;let x = -1;return;", filename: "unaryOperatorMinus"},
			{name: "var int x;let x = ~true;return;", filename: "unaryOperatorTilde"},
			{name: "var int x;let x = 3-4;return;", filename: "letStatement3minus4"},
			{name: "var int x, y;let x = x * y;return;", filename: "mathMulitply"},
			{name: "the body of the seven program", filename: "seven"},
			{name: "Screen.drawRectangle(x, y, x + width, y + height);", filename: "drawRectangle"},
			{name: "expression from the course", filename: "expressionFromCourse"},
			{name: "expression from the course with do", filename: "expressionFromCourseWithDo"},
			{name: "let statement function call with List", filename: "letStatementFunctionCall"},
			{name: "if without else", filename: "ifWithoutElse"},
			{name: "if with else", filename: "ifWithElse"},
			{name: "while statement", filename: "whileStatement"},
			{name: "multiple flow control statements", filename: "flowControl"},
			{name: "test that false constant is correct", filename: "testFalseToken"},
			{name: "construction and manipulation of objects", filename: "objectManipulation"},
			// {name: "string handling", filename: "outputString"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := readFile(t, dir+tc.filename+".jack")
				want := readFile(t, dir+tc.filename+".vm")

				tempDir := t.TempDir()
				out := filepath.Join(tempDir, tc.filename+".vm")
				vw := vmWriter.NewVmWriter(out)
				tknzr := token.NewTokenizer(input)
				e := engine.NewEngine(&tknzr, &vw)

				e.Tknzr.Advance()
				e.CompileSubroutineDec()
				vw.Close()

				assert.Empty(t, e.Errors)
				assert.FileExists(t, out)

				got, err := os.ReadFile(out)
				fmt.Println(string(got))
				assert.NoError(t, err)
				assert.Equal(t, want, string(got))
				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
			})
		}
	})
}

func Test_compileClass(t *testing.T) {
	t.Run("Testing happy classes", func(t *testing.T) {
		type testCase struct {
			name string
			fp   string
			dir  string
		}

		dir := "own/classes/"
		testCases := []testCase{
			{name: "2d point from course", fp: "Point2d"},
			{name: "3d point from course", fp: "Point3d"},
			{name: "method call in Point class", fp: "PointMethod"},
			{name: "seven program", fp: "Seven"},
			{name: "Main of Square game", fp: "SquareGameMain"},
			{name: "field var is object, calling method on that object", fp: "SquareGameFieldVarWithDo"},
			{name: "field var is object, calling method on that object with do", fp: "SquareGameFieldVar"},
			{name: "print the score left bottom", dir: "PrintScore/", fp: "Main"},
			// {name: "method call different params", fp: "PointMethod2"},
			// 			{name: "one static class variable", fp: "MainWith1StaticClassVarDec"},
			// 			{name: "two static class variable", fp: "MainWith2StaticClassVarDec"},
			// 			{name: "two static class variable different type", fp: "MainWith2StaticClassVarDec2Types"},
			// 			{name: "main arbitrary class variable declarations", fp: "MainWith3ClassVarDec"},
			// 			{name: "main with empty subroutine", fp: "MainWithEmptySubroutine"},
			// 			{name: "varDec in one line", fp: "classVarDecInARow"},
			// 			{name: "main with mulitple subroutines", fp: "main2Subroutine"},
			// 			{name: "with do statement method call", fp: "mainWithDoStatements"},
			// 			{name: "with method call other class", fp: "withMethodCallOtherClass"},
			// 			{name: "with ext classes", fp: "arrayTest1"},
			// 			{name: "with more ext classes", fp: "arrayTest2"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.dir != "" {
					dir = dir + tc.dir
				}
				input := readFile(t, dir+tc.fp+".jack")
				want := readFile(t, dir+tc.fp+".vm")

				tempDir := t.TempDir()
				out := filepath.Join(tempDir, tc.fp+".vm")
				vw := vmWriter.NewVmWriter(out)
				tknzr := token.NewTokenizer(input)
				e := engine.NewEngine(&tknzr, &vw)

				e.Tknzr.Advance()
				e.CompileClass()
				vw.Close()

				assert.Empty(t, e.Errors)
				assert.FileExists(t, out)

				got, err := os.ReadFile(out)
				fmt.Println(string(got))
				assert.NoError(t, err)
				assert.Equal(t, want, string(got))
				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
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

// 	t.Run("Testing falsy class statements", func(t *testing.T) {
// 		type testCase struct {
// 			inputs  []string
// 			wantErr error
// 		}
// 		testCases := []testCase{
// 			{
// 				inputs:  []string{"", "   "},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.CLASS, token.EOF),
// 			},
// 			{
// 				inputs:  []string{"var"},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.CLASS, token.VAR),
// 			},
// 			{
// 				inputs:  []string{"class var"},
// 				wantErr: engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.VAR),
// 			},
// 			{
// 				inputs:  []string{"class name name"},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.LBRACE, "name"),
// 			},
// 			{
// 				inputs:  []string{"class name ."},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.LBRACE, "."),
// 			},
// 			{
// 				inputs:  []string{"class Main {"},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.RBRACE, token.EOF),
// 			},
// 			{
// 				inputs:  []string{"class Main { ."},
// 				wantErr: engine.NewErrSyntaxUnexpectedToken(token.RBRACE, token.DOT),
// 			},
// 			{
// 				inputs:  []string{"class Main {static name"},
// 				wantErr: engine.NewErrSyntaxUnexpectedTokenType(token.IDENTIFIER, token.EOF),
// 			},
// 			{
// 				inputs:  []string{"class Main {static "},
// 				wantErr: engine.NewErrSyntaxNotAType(token.EOF),
// 			},
// 			{
// 				inputs:  []string{"class Main {function var"},
// 				wantErr: engine.NewErrSyntaxNotAType(token.VAR),
// 			},
// 			{
// 				inputs:  []string{"class Main {function ;"},
// 				wantErr: engine.NewErrSyntaxNotAType(token.SEMICOLON),
// 			},
// 		}

// 		for _, tc := range testCases {
// 			for _, input := range tc.inputs {
// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr)

// 				e.Tknzr.Advance()
// 				e.CompileClass()

// 				assertErrorFound(t, e.Errors, tc.wantErr)
// 			}
// 		}
// 	})
// }

// func Test_classVarDec(t *testing.T) {
// 	t.Run("Testing happy ClassVarDeclarations", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}

// 		dir := "../test_programs/project_11/engine/own/classVarDec/"
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

// 		dir := "../test_programs/project_11/engine/own/subRoutineDec/"
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

// func Test_classes(t *testing.T) {
// 	t.Run("Testing happy subRoutineBodies", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}

// 		dir := "own/classes"
// 		testCases := []testCase{
// 			{name: "emptyMain", fp: "emptyMain"},
// 		}

// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				input := readFile(t, dir+tc.fp+".jack")
// 				want := readFile(t, dir+tc.fp+".vm")

// 				tempDir := t.TempDir()
// 				out := filepath.Join(tempDir, tc.fp+".vm")
// 				vw := vmWriter.NewVmWriter(out)
// 				tknzr := token.NewTokenizer(input)
// 				e := engine.NewEngine(&tknzr, &vw)

// 				e.Tknzr.Advance()
// 				e.CompileSubroutineBody()
// 				vw.Close()

// 				assert.Empty(t, e.Errors)
// 				assert.FileExists(t, out)

// 				got, err := os.ReadFile(out)
// 				fmt.Println(string(got))
// 				assert.NoError(t, err)
// 				assert.Equal(t, want, string(got))
// 				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
// 			})
// 		}
// 	})
// }

// func Test_compileProgram(t *testing.T) {
// 	t.Run("Testing happy classes", func(t *testing.T) {
// 		type testCase struct {
// 			name string
// 			fp   string
// 		}
// 		dir := "own/test_programs/project_11/engine/ArrayTest/"
// 		testCases := []testCase{
// 			{name: "arrayTest", fp: "Main"},
// 		}
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				input := readFile(t, dir+tc.fp+".jack")
// 				want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 				tknzr := token.NewTokenizer(string(input))
// 				e := engine.NewEngine(&tknzr)

// 				e.Tknzr.Advance()
// 				got := e.CompileClass()

// 				assert.Empty(t, e.Errors)
// 				assert.Equal(t, want, removeWhiteSpaces(got))
// 				assert.Equal(t, e.Tknzr.Keyword(), token.EOF)
// 			})
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

// 		dir := "../test_programs/project_11/engine/own/subRoutineBody/"
// 		testCases := []testCase{
// 			{name: "two let statements", fp: "twoLetStatements"},
// 			{name: "let, do and return statement", fp: "letDoReturn"},
// 			// {name: "do statement with dot", fp: "doWithDot"},
// 			// {name: "while statement", fp: "while"},
// 			// {name: "while with < symbol", fp: "while"},
// 			{name: "if without else statement", fp: "ifWithoutElse"},
// 			{name: "if with else statement", fp: "ifWithElse"},
// 			{name: "return with value", fp: "returnWithValue"},
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
// 				assert.Equal(t, token.EOF, e.Tknzr.Keyword())
// 			})
// 		}
// 	})
// }

// func Test_letStatement(t *testing.T) {
// t.Run("Testing happy letStatements", func(t *testing.T) {
// 	type testCase struct {
// 		name string
// 		fp   string
// 	}

// 	dir := "../test_programs/project_10/engine/own/letStatements/"
// 	testCases := []testCase{
// 		{name: "one array assignment", fp: "arrayAssignment"},
// 		{name: "array assignment 1+2 expression", fp: "arrayAssignmentExpressionLHS"},
// 		{name: "let statement string RHS ", fp: "stringRHS"},
// 		{name: "let statement keywordConst RHS ", fp: "kwConstRHS"},
// 		{name: "let statement variableName RHS ", fp: "varNameRHS"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			input := readFile(t, dir+tc.fp+".jack")
// 			want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 			tknzr := token.NewTokenizer(input)
// 			engine := engine.NewEngine(&tknzr)
// 			engine.Tknzr.Advance()

// 			got := engine.CompileLetStatement()

// 			assert.Empty(t, engine.Errors)
// 			assert.Equal(t, want, removeWhiteSpaces(got))
// 			assert.Equal(t, engine.Tknzr.Keyword(), token.EOF)
// 		})
// 	}
// })

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
// t.Run("Testing happy terms", func(t *testing.T) {
// 	type testCase struct {
// 		name string
// 		fp   string
// 	}

// 	dir := "../test_programs/project_10/engine/own/terms/"
// 	testCases := []testCase{
// 		{name: "arr[5+6]", fp: "arr5Plus6"},
// 		{name: "(x+y)", fp: "xPlusYinParenthesis"},
// 		{name: "-x", fp: "minusX"},
// 		{name: "~x", fp: "tildeX"},
// 		{name: "foo()", fp: "foo"},
// 		{name: "foo.bar(true)", fp: "foobar"},
// 		{name: "foo(x, (y + z))", fp: "withExpList"},
// 		{name: "foo.bar(x, (y + z))", fp: "fooBarExpList"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			input := readFile(t, dir+tc.fp+".jack")
// 			want := removeWhiteSpaces(readFile(t, dir+tc.fp+".xml"))

// 			tknzr := token.NewTokenizer(input)
// 			engine := engine.NewEngine(&tknzr)
// 			engine.Tknzr.Advance()

// 			got := engine.CompileTerm()

// 			assert.Empty(t, engine.Errors)
// 			assert.Equal(t, want, removeWhiteSpaces(got))
// 			assert.Equal(t, engine.Tknzr.Keyword(), token.EOF)
// 		})
// 	}
// })

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
